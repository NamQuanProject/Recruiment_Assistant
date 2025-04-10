package evaluation

import (
	"bytes"
	"encoding/json"
	"io"

	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func evaluateJobHandler(c *gin.Context) {
	inputFolder := c.Query("path")
	if inputFolder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'path' query parameter"})
		return
	}

	// Đọc job_type
	jobTypePath := filepath.Join(inputFolder, "job_type.txt")
	jobTypeBytes, err := os.ReadFile(jobTypePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read job_type.txt"})
		return
	}
	jobType := string(jobTypeBytes)

	// Đọc criteria (JD JSON)
	criteriaPath := filepath.Join(inputFolder, "jd.json")
	criteriaBytes, err := os.ReadFile(criteriaPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read jd.json"})
		return
	}
	criteria := string(criteriaBytes)

	// Đọc các file CV trong thư mục con "cvs"
	cvFolder := filepath.Join(inputFolder, "cvs")
	files, err := os.ReadDir(cvFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read cvs folder"})
		return
	}

	results := make(map[string]interface{})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		cvPath := filepath.Join(cvFolder, file.Name())
		cvBytes, err := os.ReadFile(cvPath)
		if err != nil {
			results[file.Name()] = "Failed to read file"
			continue
		}
		cv := string(cvBytes)

		// Tạo request JSON gửi sang server 8081
		payload := map[string]string{
			"job_type": jobType,
			"criteria": criteria,
			"cv":       cv,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			results[file.Name()] = "JSON marshal error: " + err.Error()
			continue
		}

		// Gửi POST request tới server AI tại :8081
		resp, err := http.Post("http://localhost:8081/evaluate", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			results[file.Name()] = "Request error: " + err.Error()
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			results[file.Name()] = "Failed to read response: " + err.Error()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			results[file.Name()] = "Error from AI server: " + string(respBody)
			continue
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			results[file.Name()] = "Failed to parse JSON: " + err.Error()
			continue
		}

		results[file.Name()] = parsed
	}

	// Trả kết quả
	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}
