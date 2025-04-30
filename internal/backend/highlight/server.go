package highlight

import (
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

	"github.com/gin-gonic/gin"
)

// HighlightRequest represents the request to highlight a CV
type HighlightRequest struct {
	PDFPath string `json:"pdf_path" binding:"required"`
	Areas   []Area `json:"areas" binding:"required"`
}

// Area represents an area in the CV that needs to be highlighted
type Area struct {
	Text        string  `json:"text"`
	Page        int     `json:"page"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "weak" or "strong"
}

// HighlightResponse represents the response from the highlight server
type HighlightResponse struct {
	HighlightedPDFPath string `json:"highlighted_pdf_path"`
	Message            string `json:"message"`
}

// RunServer starts the highlight server
func RunServer() {
	log.Println("[Highlight] Starting server...")

	r := gin.Default()

	r.POST("/highlight", func(c *gin.Context) {
		var req HighlightRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		if _, err := os.Stat(req.PDFPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", req.PDFPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist"})
			return
		}

		// Create a timestamp for unique file naming
		timestamp := time.Now().Format("20060102_150405")
		baseName := strings.TrimSuffix(filepath.Base(req.PDFPath), ".pdf")

		// Create output directory
		outputDir := filepath.Join("storage", "highlighted_pdfs", "highlight_"+timestamp)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			log.Printf("Failed to create output directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output directory"})
			return
		}

		// Copy the original PDF to the output directory
		copiedPDFPath := filepath.Join(outputDir, baseName+".pdf")
		if err := copyFile(req.PDFPath, copiedPDFPath); err != nil {
			log.Printf("Failed to copy PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy PDF"})
			return
		}

		// Save areas to JSON file
		areasJSON, err := json.MarshalIndent(req.Areas, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process areas"})
			return
		}

		// Use consistent file name for areas JSON
		areasPath := filepath.Join(outputDir, "areas.json")
		if err := os.WriteFile(areasPath, areasJSON, 0644); err != nil {
			log.Printf("Failed to write areas file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save areas"})
			return
		}

		// Run the Python script to highlight the PDF
		pythonScriptPath := filepath.Join("internal", "backend", "highlight", "highlight_pdf.py")
		highlightedPDFPath := filepath.Join(outputDir, baseName+"_highlighted.pdf")

		cmd := exec.Command("python3", pythonScriptPath, copiedPDFPath, areasPath, highlightedPDFPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing Python script: %v\nOutput: %s", err, string(output))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to highlight PDF"})
			return
		}

		// Return the path to the highlighted PDF
		c.JSON(http.StatusOK, HighlightResponse{
			HighlightedPDFPath: highlightedPDFPath,
			Message:            "PDF highlighted successfully",
		})
	})

	fmt.Println("Highlight server running at http://localhost:8083")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // fallback khi chạy local
	}

	r.Run(":" + port)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
