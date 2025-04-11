package highlight

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// WebServer handles the web interface for CV highlighting
type WebServer struct {
	uploadDir string
}
type getHlCVRequest struct {
	JobTitle       string `json:"job_title"`
	JobDetailsPath string `json:"job_details_path"`
	PdfPath        string `json:"pdf_path"`
	EvalRefPath    string `json:"evaluation_path"`
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

		var req getHlCVRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		//jobtile = string of req.JobTitle
		jobTitle := req.JobTitle
		//read jobdetails string from req.JobDetailsPath
		jobDetailsBytes, err := os.ReadFile(req.JobDetailsPath)
		if err != nil {
			log.Printf("Failed to read job details file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read job details file"})
			return
		}
		jobDetails := string(jobDetailsBytes)

		pdfpath := req.PdfPath
		//pdfpath, jobtile, jobdetailspath, evaluationrefencepath
		// Get the uploaded file

		// // Get job title and details
		// jobTitle := c.PostForm("jobTitle")
		// jobDetails := c.PostForm("jobDetails")
		// evaluationReferenceStr := c.PostForm("evaluationReference")

		// if jobTitle == "" || jobDetails == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Job title and details are required"})
		// 	return
		// }

		// Parse evaluation reference if provided
		// var evaluationReference map[string]any
		// if evaluationReferenceStr != "" {
		// 	if err := json.Unmarshal([]byte(evaluationReferenceStr), &evaluationReference); err != nil {
		// 		log.Printf("Failed to parse evaluation reference: %v", err)
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation reference format"})
		// 		return
		// 	}
		// }

		// Create a unique filename
		// timestamp := time.Now().Format("20060102_150405")
		// filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(pdfpath))
		// fmt.Println("Filename:", filename)
		// pdfPath := filepath.Join(s.uploadDir, filename)

		// Copy the PDF file from pdfpath to pdfPath
		// sourceFile, err := os.Open(pdfpath)
		// if err != nil {
		// 	log.Printf("Failed to open source PDF file: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open source PDF file"})
		// 	return
		// }
		// defer sourceFile.Close()

		// destFile, err := os.Create(pdfPath)
		// if err != nil {
		// 	log.Printf("Failed to create destination PDF file: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create destination PDF file"})
		// 	return
		// }
		// defer destFile.Close()

		// if _, err := io.Copy(destFile, sourceFile); err != nil {
		// 	log.Printf("Failed to copy PDF file: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy PDF file"})
		// 	return
		// }

		//read evaluation reference as map[string]any from req.EvalRefPath
		evaluationReferenceBytes, err := os.ReadFile(req.EvalRefPath)
		if err != nil {
			log.Printf("Failed to read evaluation reference file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read evaluation reference file"})
			return
		}
		var evaluationReference map[string]any
		if err := json.Unmarshal(evaluationReferenceBytes, &evaluationReference); err != nil {
			log.Printf("Failed to parse evaluation reference JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evaluation reference JSON format"})
			return
		}

		// // Save the uploaded file
		// if err := c.SaveUploadedFile(file, pdfPath); err != nil {
		// 	log.Printf("Failed to save file: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		// 	return
		// }

		areas, err := FindAreas(pdfpath, jobTitle, jobDetails, "http://localhost:8081", evaluationReference)
		if err != nil {
			log.Printf("Failed to find areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze CV"})
			return
		}

		// Count strong and weak areas
		strongCount := 0
		weakCount := 0
		for _, area := range areas {
			if area.Type == "strong" {
				strongCount++
			} else {
				weakCount++
			}
		}

		// Save areas to JSON file
		areasJSON, err := json.MarshalIndent(areas, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process areas"})
			return
		}

		// Use consistent file name for areas JSON
		areasPath := filepath.Join(s.uploadDir, "areas.json")
		// Write with UTF-8 encoding
		if err := os.WriteFile(areasPath, areasJSON, 0644); err != nil {
			log.Printf("Failed to write areas file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save areas"})
			return
		}

		// Create highlight client and highlight the PDF
		highlightClient := NewClient("http://localhost:8083")
		highlightResp, err := highlightClient.HighlightPDF(pdfpath, areas)
		if err != nil {
			log.Printf("Failed to highlight PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to highlight PDF"})
			return
		}

		// Return success response with download link
		c.JSON(http.StatusOK, gin.H{
			"message":              fmt.Sprintf("CV analyzed successfully. Found %d strong areas and %d weak areas.", strongCount, weakCount),
			"highlighted_pdf_path": highlightResp.HighlightedPDFPath,
		})
	})

	// Start the server
	fmt.Println("Web server running at http://localhost:4000")
	r.Run(":4000")
}
