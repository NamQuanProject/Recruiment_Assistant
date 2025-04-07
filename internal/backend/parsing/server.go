// internal/backend/parsing/server.go
package parsing

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct to receive the request
type ParseRequest struct {
	PDFPath  string `json:"pdf_path" binding:"required"`
	TextPath string `json:"txt_path" binding:"required"`
}

// Struct for the response
type ParseResponse struct {
	Text string `json:"text"`
}

// StartParsingServer starts the parsing service on port 8082
func RunServer() {
	r := gin.Default()

	r.POST("/parse", func(c *gin.Context) {
		var req ParseRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		text, err := ExtractTextFromPDF(req.PDFPath, req.TextPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ParseResponse{Text: text})
	})

	fmt.Println("Parsing server running at http://localhost:8082")
	r.Run(":8082")
}
