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

func SubmitCVsHandler(c *gin.Context) {
	file, err := c.FormFile("pdfFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded under field 'pdfFile'"})
		return
	}

	// ext := strings.ToLower(filepath.Ext(file.Filename))
	// if ext != ".pdf" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"file":  file.Filename,
	// 		"error": "Only PDF files are allowed",
	// 	})
	// 	return
	// }

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	dst := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"file":  file.Filename,
			"error": "Failed to save file",
		})
		return
	}

	if err := ProcessCV(dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"file":   file.Filename,
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded and processed successfully", file.Filename),
		"path":    dst,
	})
}

func ProcessCV(filePath string) error {
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

	log.Printf("Processing CV: %s (Size: %d bytes)", absPath, fileInfo.Size())

	parseRequest := struct {
		InputPath  string `json:"input_path" binding:"required"`
		OutputPath string `json:"output_path" binding:"required"`
	}{
		InputPath:  absPath,
		OutputPath: "storage",
	}

	reqBody, err := json.Marshal(parseRequest)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}

	resp, err := http.Post("http://localhost:8085/parse/cv", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call CV parsing server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CV parsing server returned error: %v", resp.Status)
	}

	log.Printf("CV parsed successfully: %s", absPath)
	return nil
}
