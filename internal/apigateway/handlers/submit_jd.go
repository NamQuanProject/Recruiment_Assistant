package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadJDHandler(c *gin.Context) {
	// Single file upload
	file, err := c.FormFile("pdfFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get uploaded file",
		})
		return
	}

	// Check file extension
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only PDF files are allowed",
		})
		return
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	// Save the file
	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	// // Process the PDF
	if err := ProcessPDF(filePath, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to process PDF: %v", err),
		})
		return
	}
	// parsing.ExtractTextFromPDF(filePath)

	// Return success response without processed text
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded and processed successfully", file.Filename),
		"path":    filePath,
	})
}

func ProcessPDF(filePath string, c *gin.Context) error {
	// Open the PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open PDF file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %v", err)
	}

	log.Printf("Processing PDF file: %s (Size: %d bytes)", filePath, fileInfo.Size())

	// Prepare the request to the parsing server
	absPath, _ := filepath.Abs(filePath)
	txtFilePath := strings.TrimSuffix(absPath, ".pdf") + ".txt"

	parseRequest := struct {
		PDFPath  string `json:"pdf_path"`
		TextPath string `json:"txt_path"`
	}{
		PDFPath:  absPath,
		TextPath: txtFilePath,
	}

	// Send POST request to the parsing server
	reqBody, err := json.Marshal(parseRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to prepare request: %v", err),
		})
		return err
	}

	resp, err := http.Post("http://localhost:8082/parse", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to call parsing server: %v", err),
		})
		return err
	}
	defer resp.Body.Close()

	// Read the response
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Parsing server returned error: %v", resp.Status),
		})
		return nil
	}

	// No need to return the processed text now, so just skip decoding the response.
	// We only need to send a status.

	return nil
}
