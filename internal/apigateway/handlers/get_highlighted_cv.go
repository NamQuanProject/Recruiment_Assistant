package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type HlCVRequest struct {
	Index int `json:"index"`
}

// Cấu trúc đúng với output.json
type OutputItem struct {
	FullName        string  `json:"full_name"`
	WorkedFor       string  `json:"worked_for"`
	ExperienceLevel string  `json:"experience_level"`
	Authenticity    float64 `json:"authenticity"`
	FinalScore      float64 `json:"final_score"`
	PathToCV        string  `json:"path_to_cv"`
	PathToEval      string  `json:"path_to_evaluation"`
}

type OutputFile struct {
	List []OutputItem `json:"list"`
}

func GetHlCVHandler(c *gin.Context) {
	var req HlCVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	index := req.Index

	currentPathBytes, err := os.ReadFile("storage/current.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read storage/current.txt"})
		return
	}
	basePath := strings.TrimSpace(string(currentPathBytes))

	jobNamePath := filepath.Join(basePath, "parse", "jobname.txt")
	if _, err := os.Stat(jobNamePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "jobname.txt not found"})
		return
	}
	jobNameBytes, err := os.ReadFile(jobNamePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read jobname.txt"})
		return
	}
	jobName := strings.TrimSpace(string(jobNameBytes))

	// Read output.json
	finalOutputPath := filepath.Join("internal", "backend", "output", "output.json")
	if _, err := os.Stat(finalOutputPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "finaloutput.json not found"})
		return
	}
	data, err := os.ReadFile(finalOutputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read finaloutput.json"})
		return
	}

	var output OutputFile
	if err := json.Unmarshal(data, &output); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot parse finaloutput.json"})
		return
	}
	if index < 0 || index >= len(output.List) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index out of bounds"})
		return
	}

	item := output.List[index]
	pathtocv := item.PathToCV
	pathtoeval := item.PathToEval
	fmt.Println("Path to CV:", pathtocv)
	fmt.Println("Path to Evaluation:", pathtoeval)

	request := struct {
		JobTitle       string `json:"job_title"`
		JobDetailsPath string `json:"job_details_path"`
		PdfPath        string `json:"pdf_path"`
		EvalRefPath    string `json:"evaluation_path"`
	}{
		JobTitle:       jobName,
		JobDetailsPath: filepath.Join(basePath, "parse", "jd.txt"),
		PdfPath:        pathtocv,
		EvalRefPath:    pathtoeval,
	}
	fmt.Print("Request to server:", request)

	requestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot marshal request body"})
		return
	}
	URL := os.Getenv("BACKEND_URL")
	if URL == "" {
		URL = "http://localhost:8080" // Default URL for local testing
	}
	// URL := "https://backend-service-mjv8.onrender.com"
	// resp, err := http.Post("http://localhost:4000/analyze-cv", "application/json", strings.NewReader(string(requestBody)))
	resp, err := http.Post(fmt.Sprintf("%s/analyze-cv", URL), "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call server"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server returned an error", "status": resp.StatusCode})
		return
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode server response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"highlighted_pdf_path": responseBody["highlighted_pdf_path"],
	})
}
