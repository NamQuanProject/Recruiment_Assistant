package output

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Handler for processing and returning the final output
func outputHandler(c *gin.Context) {
    var job struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Company     string `json:"company"`
    }

    // Bind JSON input to the job struct
    if err := c.ShouldBindJSON(&job); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Simulate final processing
    job.Description = "Final Output: " + job.Description

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Job processed successfully",
        "job":     job,
    })
}

// Initialize and run the Gin server
func RunServer() {
    r := gin.Default()

    // Define API routes
    r.POST("/output", outputHandler)

    // Start the server on port 8082
    r.Run(":8082")
}