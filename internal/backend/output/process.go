package output

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"

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

func outputHandler(c *gin.Context) {
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
				FullName string `json:"FullName"`
			} `json:"PersonalInfo"`
			FinalScore float64 `json:"FinalScore"`
		}

		if err := json.Unmarshal(data, &evaluation); err != nil {
			continue // Skip files that can't be parsed
		}

		// Create candidate entry
		candidate := Candidate{
			FullName:         evaluation.PersonalInfo.FullName,
			WorkedFor:        "N/A", // This would need to be extracted from the evaluation
			ExperienceLevel:  "N/A", // This would need to be extracted from the evaluation
			Authenticity:     0.0,   // This would need to be extracted from the evaluation
			FinalScore:       evaluation.FinalScore,
			PathToCV:         "N/A", // This would need to be extracted from the evaluation
			PathToEvaluation: file,
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
	currentDir := "internal/backend/output"
	os.WriteFile(filepath.Join(currentDir, "output.json"), jsonData, 0644)

	c.JSON(http.StatusOK, response)
}
