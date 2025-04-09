package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SubmitJDHandler(c *gin.Context) {
	// Single file upload
	file, err := c.FormFile("pdfFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get uploaded file",
		})
		return
	}

	// // Check file extension
	// if filepath.Ext(file.Filename) != ".pdf" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Only PDF files are allowed",
	// 	})
	// 	return
	// }

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
	if err := ProcessJD(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to process File: %v", err),
		})
		return
	}

	// Return success response with the extracted text
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded and processed successfully", file.Filename),
		"path":    filePath,
	})
}

func ProcessJD(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %v", err)
	}

	absPath, _ := filepath.Abs(filePath)
	absPath = filepath.ToSlash(absPath)

	log.Printf("Processing file: %s (Size: %d bytes)", absPath, fileInfo.Size())

	parseRequest := struct {
		InputPath string `json:"input_path"`
	}{
		InputPath: absPath,
	}

	reqBody, err := json.Marshal(parseRequest)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}

	resp, err := http.Post("http://localhost:8085/parse", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call parsing server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("parsing server returned error: %v", resp.Status)
	}

	return nil
}
