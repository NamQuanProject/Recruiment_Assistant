package highlight

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PDFTextBlock represents a block of text in a PDF with its position
type PDFTextBlock struct {
	Text   string  `json:"text"`
	Page   int     `json:"page"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ExtractTextFromPDF extracts text and positions from a PDF file
func ExtractTextFromPDF(pdfPath string) ([]PDFTextBlock, error) {
	// Create a temporary directory for the output
	tempDir, err := os.MkdirTemp("", "pdf_extract_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create output file path
	outputPath := filepath.Join(tempDir, "text_blocks.json")

	// Run the Python script to extract text
	pythonScriptPath := filepath.Join("internal", "backend", "highlight", "extract_pdf_text.py")
	cmd := exec.Command("python3", pythonScriptPath, pdfPath, outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %v\nOutput: %s", err, string(output))
	}

	// Read the output file
	jsonData, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %v", err)
	}

	// Parse the JSON
	var textBlocks []PDFTextBlock
	if err := json.Unmarshal(jsonData, &textBlocks); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return textBlocks, nil
}
