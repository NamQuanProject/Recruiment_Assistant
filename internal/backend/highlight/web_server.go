package highlight

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// WebServer handles the web interface for CV highlighting
type WebServer struct {
	uploadDir string
}

// NewWebServer creates a new web server instance
func NewWebServer() *WebServer {
	// Create upload directory if it doesn't exist
	uploadDir := filepath.Join("storage", "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	return &WebServer{
		uploadDir: uploadDir,
	}
}

// Run starts the web server
func (s *WebServer) Run() {
	r := gin.Default()

	// Serve static files from templates directory
	r.LoadHTMLGlob("templates/*")
	r.Static("/storage", "./storage")

	// Serve the CV highlight page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cv_highlight.html", nil)
	})

	// Handle CV upload and analysis
	r.POST("/analyze-cv", func(c *gin.Context) {
		// Get the uploaded file
		file, err := c.FormFile("cvFile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Get job title and details
		jobTitle := c.PostForm("jobTitle")
		jobDetails := c.PostForm("jobDetails")

		if jobTitle == "" || jobDetails == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Job title and details are required"})
			return
		}

		// Create a unique filename
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(file.Filename))
		filepath := filepath.Join(s.uploadDir, filename)

		// Save the uploaded file
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			log.Printf("Failed to save file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Find weak areas using AI
		weakAreas, err := FindWeakAreas(filepath, jobTitle, jobDetails, "http://localhost:8081")
		if err != nil {
			log.Printf("Failed to find weak areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze CV"})
			return
		}

		// Create highlight client and highlight the PDF
		highlightClient := NewClient("http://localhost:8083")
		highlightResp, err := highlightClient.HighlightPDF(filepath, weakAreas)
		if err != nil {
			log.Printf("Failed to highlight PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to highlight PDF"})
			return
		}

		// Return success response with download link
		c.JSON(http.StatusOK, gin.H{
			"message":            "CV analyzed and highlighted successfully",
			"highlighted_pdf_path": highlightResp.HighlightedPDFPath,
		})
	})

	// Start the server
	fmt.Println("Web server running at http://localhost:3001")
	r.Run(":3001")
} 