package aiservices

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	agent, err := NewAIAgent(Config{}, true)
	r.GET("/ai", func(c *gin.Context) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI agent"})
			return
		}

		prompt := "List 3 popular cookie recipes."

		final_prompt := prompt

		result := agent.CallChatGemini(final_prompt)

		c.JSON(http.StatusOK, gin.H{
			"Question": prompt,
			"Response": result["Response"],
		})
	})

	r.POST("/ai", func(c *gin.Context) {})

	agent.Close()
	r.Run(":8081")
}
