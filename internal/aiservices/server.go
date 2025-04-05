package aiservices

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	r.GET("/ai", func(c *gin.Context) {
		// Initialize the AI agent
		agent, err := NewAIAgent(Config{}, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI agent"})
			return
		}
		prompt := "List 3 popular cookie recipes."

		agent.CallChatGemini(prompt)

		c.JSON(http.StatusOK, gin.H{
			"prompt": prompt,
		})
		agent.Close()

	})

	// Start the server on port 8081
	r.Run(":8081")
}
