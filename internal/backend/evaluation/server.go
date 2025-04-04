package evaluation

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
)

// Handler for evaluating job data
func evaluateJobHandler(c *gin.Context) {
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

    // Simulate evaluation logic
    job.Description = "Evaluated " + job.Description

    // Forward the evaluated job to the output service
    jsonData, _ := json.Marshal(job)
    resp, err := http.Post("http://localhost:8082/output", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to output service"})
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "application/json", body)
}

// Initialize and run the Gin server
func RunServer() {
    r := gin.Default()

    // Define API routes
    r.POST("/evaluate", evaluateJobHandler)

    // Start the server on port 8081
    r.Run(":8081")
}