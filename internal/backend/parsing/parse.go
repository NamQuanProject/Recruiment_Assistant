// internal/backend/parsing/parse.go
package parsing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return "", err
	}

	// Copy original zip file into storage
	zipCopyPath := filepath.Join(baseDir, "original.zip")
	zipCopy, err := os.Create(zipCopyPath)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, 0) // reset to start of zip
	if err != nil {
		return "", err
	}
	_, err = io.Copy(zipCopy, file)
	if err != nil {
		return "", err
	}

	// Extract files
	extractedPath := filepath.Join(baseDir, "extracted")
	err = unzip(zipCopyPath, extractedPath)
	if err != nil {
		return "", err
	}

	// Output path for texts
	outputPath := filepath.Join(baseDir, "texts")
	os.Mkdir(outputPath, 0755)

	// Call Python batch processing
	pythonScriptPath := filepath.Join("internal", "backend", "parsing", "extract_pdf.py")
	cmd := exec.Command("python", pythonScriptPath, "-batch", "true", extractedPath, outputPath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing script: %v\nOutput: %s", err, string(out))
	}

	return baseDir, nil
}
