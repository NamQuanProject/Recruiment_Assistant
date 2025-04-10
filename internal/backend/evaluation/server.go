package evaluation

import (
	"github.com/gin-gonic/gin"
)

// Initialize and run the Gin server
func RunServer() {
	r := gin.Default()

	// Define API routes
	r.POST("/evaluate", evaluateJobHandler)

	// Start the server on port 8082
	r.Run(":8082")
}
