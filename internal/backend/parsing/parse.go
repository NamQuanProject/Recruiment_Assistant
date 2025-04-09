// internal/backend/parsing/parse.go
package parsing

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func ExtractTextFromPDF(pdfPath string) (string, error) {
	// Get timestamp and base name of the PDF
	timeStamp := time.Now().Format("20060102_150405")
	baseName := strings.TrimSuffix(filepath.Base(pdfPath), ".pdf")

	// Define storage paths
	uploadDir := filepath.Join("storage", "cv_pdfs", "upload_"+timeStamp)
	textDir := uploadDir

	// Create necessary directories
	if err := os.MkdirAll(textDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories: %v", err)
	}

	// Copy the original PDF to the upload directory
	copiedPDFPath := filepath.Join(uploadDir, baseName+".pdf")
	if err := copyFile(pdfPath, copiedPDFPath); err != nil {
		return "", fmt.Errorf("failed to copy PDF: %v", err)
	}

	// Prepare Python script and output path
	pythonScriptPath := filepath.Join("internal", "backend", "parsing", "extract_pdf.py")
	outputPath := filepath.Join(textDir, baseName+".txt")

	// Run the Python script
	cmd := exec.Command("python", pythonScriptPath, copiedPDFPath, outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing Python script: %v\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(outputPath)), nil
}

func ExtractTextFromZip(zipPath string) (string, error) {
	// Open and hash zip for uniqueness
	file, err := os.Open(zipPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashStr := hex.EncodeToString(hash.Sum(nil))[:10]

	// Create unique output directory
	timestamp := time.Now().Format("20060102_150405")
	baseDir := filepath.Join("storage", "cv_zips", fmt.Sprintf("upload_%s_%s", timestamp, hashStr))
	log.Printf("Base dir: %s", baseDir)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	// Save original zip file with its actual name
	originalName := filepath.Base(zipPath)
	zipCopyPath := filepath.Join(baseDir, originalName)
	zipCopy, err := os.Create(zipCopyPath)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(zipCopy, file); err != nil {
		return "", err
	}

	// Unzip into "extracted" directory
	extractedPath := filepath.Join(baseDir)
	if err := os.MkdirAll(extractedPath, 0755); err != nil {
		return "", err
	}
	if err := unzip(zipCopyPath, extractedPath); err != nil {
		return "", err
	}

	// Create output folder for texts
	outputPath := filepath.Join(baseDir, "texts")
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return "", err
	}

	extractedPath = filepath.Join(extractedPath, strings.TrimSuffix(originalName, filepath.Ext(originalName)))
	log.Printf("Extracted_path: %s", extractedPath)

	// Run python script
	pythonScriptPath := filepath.Join("internal", "backend", "parsing", "extract_pdf.py")
	cmd := exec.Command("python", pythonScriptPath, "-batch", "true", extractedPath, outputPath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing script: %v\nOutput: %s", err, string(out))
	}

	return outputPath, nil
}

func ExtractJsonFromText(textPath string, outputPath string) (map[string]interface{}, error) {
	log.Println("Function entered")

	// 1. Read the file content
	fileBytes, err := os.ReadFile(textPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	cvText := string(fileBytes)

	// 2. Prepare JSON body
	requestBody := map[string]string{
		"job_raw_text": cvText,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 3. Send POST request
	url := "http://localhost:8081/ai/parsing"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// 4. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 5. Parse top-level JSON
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// 6. Access the "Response" field and parse the wrapped JSON string
	outerMap, ok := jsonResponse["Response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'Response' to be a map")
	}

	responseStr, ok := outerMap["Response"].(string)
	if !ok {
		return nil, fmt.Errorf("expected inner 'Response' to be a string")
	}

	// 7. Clean up the string (remove Markdown formatting)
	responseStr = strings.TrimSpace(responseStr)
	responseStr = strings.TrimPrefix(responseStr, "```json")
	responseStr = strings.TrimSuffix(responseStr, "```")
	responseStr = strings.TrimSpace(responseStr)

	// 8. Parse the cleaned JSON string into a map
	var parsedData map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &parsedData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cleaned JSON: %w", err)
	}

	// 9. Determine output path
	var jsonFilename string
	if outputPath == "" {
		jsonFilename = strings.TrimSuffix(textPath, filepath.Ext(textPath)) + ".json"
	} else {
		jsonFilename = outputPath
	}

	// 10. Save the final parsed JSON to a file
	jsonBytes, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format final JSON: %w", err)
	}

	err = os.WriteFile(jsonFilename, jsonBytes, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write final JSON file: %w", err)
	}

	log.Printf("✅ JSON saved to: %s", jsonFilename)
	return parsedData, nil
}

func ExtractJsonFromTextBatch(folderPath string) error {
	// Step 1: Create "jsons" directory next to folderPath
	jsonDir := filepath.Join(filepath.Dir(folderPath), "jsons")
	err := os.MkdirAll(jsonDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create json output directory: %w", err)
	}

	// Step 2: Read directory entries
	log.Print(folderPath)
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read folder: %w", err)
	}

	// Step 3: Process each .txt file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".txt" {
			continue
		}

		txtFilePath := filepath.Join(folderPath, file.Name())
		log.Printf("📄 Processing file: %s", txtFilePath)
		jsonFilename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + ".json"
		outputPath := filepath.Join(jsonDir, jsonFilename)
		_, er := ExtractJsonFromText(txtFilePath, outputPath)
		if er != nil {
			return fmt.Errorf("failed to extract file %s to json: %w", txtFilePath, er)
		}
	}

	return nil
}
