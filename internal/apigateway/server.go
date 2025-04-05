package apigateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Example handler for a health check endpoint
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "API Gateway is running",
	})
}

// Forward job creation request to the backend evaluation service
func createJobHandler(c *gin.Context) {
	var job struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Company     string `json:"company" binding:"required"`
	}

	// Bind JSON input to the job struct
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Serialize the job struct back to JSON
	jsonData, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize request body"})
		return
	}

	// Forward the request to the backend evaluation service
	resp, err := http.Post("http://localhost:8081/evaluate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to evaluation service"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from evaluation service"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// Initialize and run the Gin server
func RunServer() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Define API routes
	r.GET("/health", healthCheckHandler)
	r.POST("/jobs", createJobHandler)
	// Start the server on port 8080
	r.Run(":8080")
}
