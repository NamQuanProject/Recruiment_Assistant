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

type JDRequest struct {
	JobName                string `json:"job_name"`
	CompanyDescriptionPath string `json:"company_jd"`
}

func RunServer() {
	log.Println("[Parsing] Starting server...")

	r := gin.Default()

	r.POST("/parse/cv", func(c *gin.Context) {
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
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Print("ExtractJsonFromText:", textPath)
			_, er := ExtractJsonFromText(textPath, "")

			if er != nil {
				log.Print(er)
				c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "PDF processed successfully"})

		case ".zip":
			log.Printf("Detected ZIP: %s", req.InputPath)
			extractedPath, err := ExtractTextFromZip(req.InputPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			er := ExtractJsonFromTextBatch(extractedPath)
			if er != nil {
				log.Print(er)
				c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "ZIP processed successfully"})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported file type. Only .pdf and .zip are allowed.",
			})
		}
	})

	r.POST("/parse/jd", func(c *gin.Context) {

		var req JDRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		// Check if the file exists
		if _, err := os.Stat(req.CompanyDescriptionPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", req.CompanyDescriptionPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist"})
			return
		}

		err := ExtractCategoriesFromJDText(req.JobName, req.CompanyDescriptionPath)
		if err != nil {
			log.Printf("Failed to extract and send JD: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "JD processed and sent for AI extraction successfully",
		})
	})

	fmt.Println("Parsing server running at http://localhost:8085")
	r.Run(":8085")
}
