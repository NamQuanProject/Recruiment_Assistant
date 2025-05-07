package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	EvaluationID   string `json:"evaluation_id"`
	CVID           string `json:"cv_id"`
}

func EvaluateJobHandler(c *gin.Context) {
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

	// Parse JDMainQuiteria to get total scoring scale
	var jdData struct {
		MainCategory []struct {
			ScoringScale float64 `json:"ScoringScale"`
		} `json:"MainCategory"`
		SubCategory []struct {
			ScoringScale float64 `json:"ScoringScale"`
		} `json:"SubCategory"`
	}
	if err := json.Unmarshal(jdMainBytes, &jdData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid jd.json format"})
		return
	}
	var totalScale float64
	for _, m := range jdData.MainCategory {
		totalScale += m.ScoringScale
	}
	for _, s := range jdData.SubCategory {
		totalScale += s.ScoringScale
	}

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
			EvaluationID:   strings.TrimPrefix(basePath, "storage/evaluation_"),
			CVID:           strings.TrimSuffix(f.Name(), ".txt"),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			continue
		}

		URL := os.Getenv("AI_URL")
		if URL == "" {
			URL = "http://localhost:8081"
		}
		// URL := "https://aiservice-service.onrender.com"
		// resp, err := http.Post("http://localhost:8081/ai/evaluate", "application/json", bytes.NewBuffer(payloadBytes))
		resp, err := http.Post(fmt.Sprintf("%s/ai/evaluate", URL), "application/json", bytes.NewBuffer(payloadBytes))

		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(body, &jsonResponse); err != nil {
			return
		}

		evalMap, ok := jsonResponse["evaluation"].(map[string]interface{})
		if !ok {
			return
		}

		responseStr, ok := evalMap["Response"].(string)
		if !ok {
			return
		}

		responseStr = strings.TrimSpace(responseStr)
		responseStr = strings.TrimPrefix(responseStr, "```json")
		responseStr = strings.TrimSuffix(responseStr, "```")
		responseStr = strings.TrimSpace(responseStr)

		var parsedData map[string]interface{}
		if err := json.Unmarshal([]byte(responseStr), &parsedData); err != nil {
			return
		}

		evalList, ok := parsedData["Evaluation"].([]interface{})
		if !ok {
			return
		}

		var totalScore float64
		for _, item := range evalList {
			evalItem, ok := item.(map[string]interface{})
			if !ok {
				fmt.Println("Error: invalid evaluation item format")
				continue
			}
			var score1 float64
			score, ok := evalItem["score"].(string)
			if !ok {
				fmt.Println(score, "Error: score not found or not a number")
				return
			}
			if ok {
				score1, err = strconv.ParseFloat(score, 64)
			}
			totalScore += score1
		}

		finalScore := totalScore / totalScale * 100
		parsedData["FinalScore"] = finalScore // Adding the new key "FinalScore"
		personalInfo, ok := parsedData["PersonalInfo"].(map[string]interface{})
		if ok {
			personalInfo["PathToCV"] = filepath.Join(basePath, "original", "cvs", strings.TrimSuffix(f.Name(), ".txt")+".pdf")
		}
		fmt.Print("pathtocv:", cvPath)

		jsonOutput := filepath.Join(outputFolder, strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))+".json")
		jsonBytes, err := json.MarshalIndent(parsedData, "", "  ")
		if err != nil {
			return
		}
		if err := os.WriteFile(jsonOutput, jsonBytes, 0644); err != nil {
			continue
		}
		// fmt.Println(string(jsonBytes))
		results = append(results, jsonOutput)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Evaluation completed",
		"evaluations": results,
	})
}
