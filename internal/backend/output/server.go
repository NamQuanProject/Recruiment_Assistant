package output

import (
	"os"

	"github.com/gin-gonic/gin"
)

// CriteriaScore represents the score and explanation for a single evaluation criteria
type CriteriaScore struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
}

// CVOwner represents the information about the CV owner
type CVOwner struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`
}

// EvaluationResult represents the complete evaluation output
type EvaluationResult struct {
	JobTitle     string          `json:"jobTitle"`
	Company      string          `json:"company"`
	TotalScore   float64         `json:"totalScore"`
	CriteriaList []CriteriaScore `json:"criteriaList"`
	Summary      string          `json:"summary"`
	CVOwner      CVOwner         `json:"cvOwner"`
}

// Handler for processing and returning the final output

// Initialize and run the Gin server
func RunServer() {
	r := gin.Default()

	// Define API routes
	// r.POST("/output", outputHandler)

	// Start the server on port 8084
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // fallback khi chạy local
	}

	r.Run(":" + port)
}
