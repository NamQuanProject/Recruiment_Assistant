package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Candidate struct {
	FullName         string  `json:"full_name"`
	WorkedFor        string  `json:"worked_for"`
	ExperienceLevel  string  `json:"experience_level"`
	Authenticity     float64 `json:"authenticity"`
	FinalScore       float64 `json:"final_score"`
	PathToCV         string  `json:"path_to_cv"`
	PathToEvaluation string  `json:"path_to_evaluation"`
}

type OutputResponse struct {
	List []Candidate `json:"list"`
}

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func OutputHandler(c *gin.Context) {
	var evaluationFolder string

	// Handle both GET and POST requests
	if c.Request.Method == "POST" {
		var request struct {
			EvaluationFolder string `json:"evaluation_folder"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		evaluationFolder = request.EvaluationFolder
	} else {
		// For GET requests, use the default evaluation folder
		evaluationFolder = "storage/evaluation"
	}

	// Read all JSON files in the evaluation folder
	files, err := filepath.Glob(filepath.Join(evaluationFolder, "*.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read evaluation folder"})
		return
	}

	var candidates []Candidate

	// Process each evaluation file
	for _, file := range files {
		// Read the JSON file
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}

		var evaluation struct {
			PersonalInfo struct {
				FullName        string `json:"FullName"`
				WorkFor         string `json:"WorkFor"`
				ExperienceYears string `json:"Experience_Years"`
				PathToCV        string `json:"PathToCV"`
				PathToCVAlt     string `json:"path_to_cv"`
			} `json:"PersonalInfo"`
			Authenticity interface{} `json:"Authenticity"`
			FinalScore   float64     `json:"FinalScore"`
		}

		if err := json.Unmarshal(data, &evaluation); err != nil {
			continue // Skip files that can't be parsed
		}

		// Convert Authenticity to float64
		var authenticity float64
		switch v := evaluation.Authenticity.(type) {
		case string:
			authenticity, _ = strconv.ParseFloat(v, 64)
		case float64:
			authenticity = v
		case int:
			authenticity = float64(v)
		}

		// Get the CV path, checking both field names
		cvPath := evaluation.PersonalInfo.PathToCV
		if cvPath == "" {
			cvPath = evaluation.PersonalInfo.PathToCVAlt
		}

		// Convert file path to use forward slashes
		evalPath := strings.ReplaceAll(file, "\\", "/")

		// Create candidate entry
		candidate := Candidate{
			FullName:         evaluation.PersonalInfo.FullName,
			WorkedFor:        evaluation.PersonalInfo.WorkFor,
			ExperienceLevel:  evaluation.PersonalInfo.ExperienceYears,
			Authenticity:     authenticity,
			FinalScore:       evaluation.FinalScore,
			PathToCV:         cvPath,
			PathToEvaluation: evalPath,
		}

		candidates = append(candidates, candidate)
	}

	// Sort candidates by FinalScore in descending order
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].FinalScore > candidates[j].FinalScore
	})

	// Return the sorted list
	response := OutputResponse{
		List: candidates,
	}

	// Save response to a file
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}
	// Write to the current directory
	os.WriteFile("internal/backend/output/output.json", jsonData, 0644)

	c.JSON(http.StatusOK, response)
}
