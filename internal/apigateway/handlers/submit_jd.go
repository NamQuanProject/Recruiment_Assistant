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
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitJDHandler(c *gin.Context) {
	// Get job_name from form-data
	jobName := c.PostForm("job_name")
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing job_name in form-data",
		})
		return
	}

	// Get uploaded PDF file
	file, err := c.FormFile("pdf_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing or invalid pdf_file",
		})
		return
	}
	temptime := time.Now().Format("20060102_150405")

	path := filepath.Join("storage", fmt.Sprintf("evaluation_%s", temptime))
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create uploads directory",
		})
		return
	}

	// Save path to current.txt in storage
	currentFilePath := filepath.Join("storage", "current.txt")
	currentFile, err := os.OpenFile(currentFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create or open current.txt",
		})
		return
	}
	defer currentFile.Close()

	if _, err := currentFile.WriteString(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write to current.txt",
		})
		return
	}

	// Create uploads folder if not exists

	// Save uploaded file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	originalPath := filepath.Join(path, "original")

	if err := os.MkdirAll(originalPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create create directory",
		})
		return
	}
	// Write job_name to jobnam.txt in the specified path
	jobNameFilePath := filepath.Join(originalPath, "jobname.txt")
	jobNameFile, err := os.OpenFile(jobNameFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create or open jobname.txt",
		})
		return
	}
	defer jobNameFile.Close()

	if _, err := jobNameFile.WriteString(jobName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write to jobname.txt",
		})
		return
	}

	filePath := filepath.Join(originalPath, fmt.Sprint("jd", ext))

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save uploaded file",
		})
		return
	}
	parsePath := filepath.Join(path, "parse")
	if err := os.MkdirAll(parsePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create create directory",
		})
		return
	}

	// Copy jobname.txt to parsePath/jobname.txt
	parseJobNameFilePath := filepath.Join(parsePath, "jobname.txt")
	input, err := os.ReadFile(jobNameFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read jobname.txt",
		})
		return
	}

	if err := os.WriteFile(parseJobNameFilePath, input, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to copy jobname.txt to parsePath",
		})
		return
	}

	TxtPath := filepath.Join(parsePath, "jd.txt")
	JsonPath := filepath.Join(parsePath, "jd.json")

	// Process JD file
	if err := ProcessJD(jobName, filePath, TxtPath, JsonPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to process JD: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("JD '%s' uploaded and processed successfully", file.Filename),
		"path":    JsonPath,
	})
}

func ProcessJD(jobName, filePath, TxtPath, JsonPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %v", err)
	}

	// absPath, _ := filepath.Abs(filePath)
	// absPath = filepath.ToSlash(absPath)

	log.Printf("Processing JD file: %s (Size: %d bytes)", TxtPath, fileInfo.Size())

	// Build JSON request
	parseRequest := struct {
		JobName                string `json:"job_name"`
		CompanyDescriptionPath string `json:"company_jd"`
		TxtPath                string `json:"txt_path"`
		JsonPath               string `json:"json_path"`
	}{
		JobName:                jobName,
		CompanyDescriptionPath: filePath,
		TxtPath:                TxtPath,
		JsonPath:               JsonPath,
	}

	reqBody, err := json.Marshal(parseRequest)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}

	// Call parsing server
	URL := os.Getenv("PARSE_URL")
	if URL == "" {
		URL = "http://localhost:8080" // Default URL for local testing
	}
	// URL := "https://parsing-service.onrender.com"
	// resp, err := http.Post("http://localhost:8085/parse/jd", "application/json", bytes.NewBuffer(reqBody))
	resp, err := http.Post(fmt.Sprintf("%s/parse/jd", URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call parsing server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("parsing server returned error: %v", resp.Status)
	}

	log.Printf("JD parsed successfully")
	return nil
}
