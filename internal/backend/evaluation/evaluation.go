package evaluation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type evalRequest struct {
	InputPath string `json:"input_path" binding:"required"`
}

type JDRequest struct {
	JobName        string `json:"job_name"`
	JDMainQuiteria string `json:"jd_main_quiteria"`
	CVRawText      string `json:"cv_raw_text"`
}

func evaluateJobHandler(c *gin.Context) {
	var req evalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	basePath := strings.TrimSpace(req.InputPath)
	parsePath := filepath.Join(basePath, "parse")

	jobnameBytes, err := os.ReadFile(filepath.Join(parsePath, "jobname.txt"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read jobname.txt"})
		return
	}
	jobName := strings.TrimSpace(string(jobnameBytes))

	jdMainBytes, err := os.ReadFile(filepath.Join(parsePath, "jd.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read jd.json"})
		return
	}
	// var jdMain []string
	// if err := json.Unmarshal(jdMainBytes, &jdMain); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid jd.json format"})
	// 	return
	// }

	cvsFolder := filepath.Join(parsePath, "cvs")
	outputFolder := filepath.Join(basePath, "evaluation")
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output folder"})
		return
	}

	files, err := os.ReadDir(cvsFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read cvs folder"})
		return
	}

	var results []string
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".txt") {
			continue
		}

		cvPath := filepath.Join(cvsFolder, f.Name())
		cvTextBytes, err := os.ReadFile(cvPath)
		if err != nil {
			continue
		}

		payload := JDRequest{
			JobName:        jobName,
			JDMainQuiteria: string(jdMainBytes),
			CVRawText:      string(cvTextBytes),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			continue
		}

		resp, err := http.Post("http://localhost:8081/ai/evaluate", "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		fmt.Printf("Response status: %s\n", resp.Status)
		fmt.Printf("Response body: %s\n", string(body))

		// 5. Parse top-level JSON
		var jsonResponse map[string]interface{}
		err = json.Unmarshal(body, &jsonResponse)
		if err != nil {
			return
		}

		// 6. Access the "Response" field and parse the wrapped JSON string
		outerMap, ok := jsonResponse["evaluation"].(map[string]interface{})
		if !ok {
			return
		}

		responseStr, ok := outerMap["Response"].(string)
		if !ok {
			return
		}

		// 7. Clean up the string (remove Markdown formatting)
		responseStr = strings.TrimSpace(responseStr)
		responseStr = strings.TrimPrefix(responseStr, "```json")
		responseStr = strings.TrimSuffix(responseStr, "```")
		responseStr = strings.TrimSpace(responseStr)

		// 8. Parse the cleaned JSON string into a map
		var parsedData map[string]interface{}
		err = json.Unmarshal([]byte(responseStr), &parsedData)
		if err != nil {
			return
		}

		// 9. Determine output path
		jsonOutput := filepath.Join(outputFolder, strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))+".json")

		// 10. Save the final parsed JSON to a file
		jsonBytes, err := json.MarshalIndent(parsedData, "", "  ")
		if err != nil {
			return
		}

		err = os.WriteFile(jsonOutput, jsonBytes, 0644)
		if err != nil {
			continue
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Evaluation completed",
		"evaluations": results,
	})
}
