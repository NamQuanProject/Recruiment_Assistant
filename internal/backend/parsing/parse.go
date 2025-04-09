// internal/backend/parsing/parse.go
package parsing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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

	return strings.TrimSpace(string(output)), nil
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

	return baseDir, nil
}
