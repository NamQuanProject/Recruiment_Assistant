package parsing

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ParseRequest struct {
	InputPath string `json:"input_path" binding:"required"`
}

func RunServer() {
	log.Println("[Parsing] Starting server...")

	r := gin.Default()

	r.POST("/parse", func(c *gin.Context) {
		var req ParseRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", req.InputPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist"})
			return
		}

		ext := strings.ToLower(filepath.Ext(req.InputPath))

		switch ext {
		case ".pdf":
			log.Printf("Detected PDF: %s", req.InputPath)
			textPath, err := ExtractTextFromPDF(req.InputPath)
			_, er := ExtractJsonFromText(textPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if er != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "PDF processed successfully"})

		case ".zip":
			log.Printf("Detected ZIP: %s", req.InputPath)
			_, err := ExtractTextFromZip(req.InputPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "ZIP processed successfully"})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported file type. Only .pdf and .zip are allowed.",
			})
		}
	})

	fmt.Println("Parsing server running at http://localhost:8085")
	r.Run(":8085")
}
