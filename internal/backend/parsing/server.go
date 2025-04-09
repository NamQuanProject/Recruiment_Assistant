package parsing

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Struct to receive the request
type ParseRequest struct {
	PDFPath  string `json:"pdf_path" binding:"required"`
	TextPath string `json:"txt_path" binding:"required"`
}

// StartParsingServer starts the parsing service on port 8082
func RunServer() {
	r := gin.Default()

	// POST route for parsing
	r.POST("/parse", func(c *gin.Context) {
		var req ParseRequest

		// Bind the incoming JSON request to the ParseRequest struct
		log.Printf("Bind JSON...\n")
		if err := c.ShouldBindJSON(&req); err != nil {
			// If the JSON is invalid, return a bad request status
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		log.Printf("Extracting Text...\n")
		// Call the processing function (assuming it returns an error if any)
		text, err := ExtractTextFromPDF(req.PDFPath)
		log.Printf("Extracted Text...\n")

		if err != nil {
			// If there is an error during processing, return an internal server error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send only a status message indicating success
		c.JSON(http.StatusOK, gin.H{
			"message": "File processed successfully",
			"text":    text,
		})
	})

	// Start the server on port 8082
	fmt.Println("Parsing server running at http://localhost:8082")
	r.Run(":8082")
}
