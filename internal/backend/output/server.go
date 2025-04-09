package output

import (
	"net/http"

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
	JobTitle     string         `json:"jobTitle"`
	Company      string         `json:"company"`
	TotalScore   float64        `json:"totalScore"`
	CriteriaList []CriteriaScore `json:"criteriaList"`
	Summary      string         `json:"summary"`
	CVOwner      CVOwner        `json:"cvOwner"`
}

// Handler for processing and returning the final output
func outputHandler(c *gin.Context) {
	var evaluationResult EvaluationResult

	// Bind JSON input to the evaluation result struct
	if err := c.ShouldBindJSON(&evaluationResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input format",
			"error":   err.Error(),
		})
		return
	}

	// Format the response
	response := gin.H{
		"status":  "success",
		"message": "Evaluation completed successfully",
		"data": gin.H{
			"jobTitle":   evaluationResult.JobTitle,
			"company":    evaluationResult.Company,
			"totalScore": evaluationResult.TotalScore,
			"criteria":   evaluationResult.CriteriaList,
			"summary":    evaluationResult.Summary,
			"cvOwner": gin.H{
				"name":        evaluationResult.CVOwner.Name,
				"email":       evaluationResult.CVOwner.Email,
				"phoneNumber": evaluationResult.CVOwner.PhoneNumber,
				"location":    evaluationResult.CVOwner.Location,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// Initialize and run the Gin server
func RunServer() {
	r := gin.Default()

	// Define API routes
	r.POST("/output", outputHandler)

	// Start the server on port 8082
	r.Run(":8083")
}
