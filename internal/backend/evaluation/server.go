package evaluation

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Initialize and run the Gin server
func RunServer() {
	r := gin.Default()

	// Define API routes
	r.POST("/evaluate", evaluateJobHandler)

	// Start the server on port 8082
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // fallback khi chạy local
	}

	r.Run(":" + port)
}
