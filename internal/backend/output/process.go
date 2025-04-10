package output

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// type CriteriaScore struct {
// 	Name        string  `json:"name"`
// 	Score       float64 `json:"score"`
// 	Explanation string  `json:"explanation"`
// }

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
