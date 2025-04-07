package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// Process the PDF
	if err := ProcessPDF(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to process PDF: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded and processed successfully", file.Filename),
		"path":    filePath,
	})
}

func ProcessPDF(filePath string) error {
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

	// Here you can add your PDF processing logic
	// For example, using unipdf or pdfcpu libraries

	

	return nil
}
