package apigateway

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/handlers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/initializers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/user"
	"github.com/gin-contrib/cors" //chau added this
	"github.com/gin-gonic/gin"
)

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

// HighlightRequest represents the request to highlight a CV
type HighlightRequest struct {
	PDFPath string `json:"pdf_path" binding:"required"`
	Areas   []Area `json:"areas" binding:"required"`
}

// Area represents an area in the CV that needs to be highlighted
type Area struct {
	Text        string  `json:"text"`
	Page        int     `json:"page"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "weak" or "strong"
}

// HighlightResponse represents the response from the highlight server
type HighlightResponse struct {
	HighlightedPDFPath string `json:"highlighted_pdf_path"`
	Message            string `json:"message"`
}

type ParseRequest struct {
	InputPath  string `json:"input_path" binding:"required"`
	OutputPath string `json:"output_path" binding:"required"`
}

type JDRequest struct {
	JobName                string `json:"job_name"`
	CompanyDescriptionPath string `json:"company_jd"`
	TxtPath                string `json:"txt_path"`
	JsonPath               string `json:"json_path"`
}

func Init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func RunServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",              // local frontend
			"https://frontend-eqtg.onrender.com", // deployed frontend
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})) // Chau added this

	r.Static("/storage", "./storage")
	r.Static("/internal", "./internal")

	// Authentication routes
	r.POST("/signup", user.Signup)
	r.POST("/login", user.Login)
	r.POST("/logout", user.Logout)

	// Serve HTML form upload
	//r.GET("/", func(c *gin.Context) {
	//c.HTML(http.StatusOK, "upload.html", nil)
	//})

	// Create a route for testing the upload without auth (temporary)
	//r.GET("/upload-test", func(c *gin.Context) {
	//c.HTML(http.StatusOK, "upload.html", nil)
	//})

	// Job description routes
	r.POST("/submitJD", handlers.SubmitJDHandler)
	r.POST("/submitCVs", handlers.SubmitCVsHandler)
	r.POST("/getHlCV", handlers.GetHlCVHandler)

	// parsing routes

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

		err := handlers.ExtractCategoriesFromJDText(req.JobName, req.CompanyDescriptionPath, req.TxtPath, req.JsonPath)
		if err != nil {
			log.Printf("Failed to extract and send JD: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "JD processed and sent for AI extraction successfully",
		})
	})

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
		if ext != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{
				"file":  req.InputPath,
				"error": "Only PDF files are allowed",
			})
			return
		}

		log.Printf("Detected PDF: %s", req.InputPath)
		err := handlers.ExtractTextFromPDF(req.InputPath, req.OutputPath)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Print("ExtractJsonFromText:", req.OutputPath)

		c.JSON(http.StatusOK, gin.H{"message": "PDF processed successfully"})

	})

	// evaluation routes
	r.POST("/evaluate", handlers.EvaluateJobHandler)
	// r.POST("/evaluate", handlers.evaluateJobHandler)

	// output routes
	r.POST("/output", handlers.OutputHandler)

	// highlight routes
	r.POST("/highlight", func(c *gin.Context) {
		var req HighlightRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		if _, err := os.Stat(req.PDFPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", req.PDFPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist"})
			return
		}

		// Create a timestamp for unique file naming
		timestamp := time.Now().Format("20060102_150405")
		baseName := strings.TrimSuffix(filepath.Base(req.PDFPath), ".pdf")

		// Create output directory
		outputDir := filepath.Join("storage", "highlighted_pdfs", "highlight_"+timestamp)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			log.Printf("Failed to create output directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output directory"})
			return
		}

		// Copy the original PDF to the output directory
		copiedPDFPath := filepath.Join(outputDir, baseName+".pdf")
		if err := copyFile(req.PDFPath, copiedPDFPath); err != nil {
			log.Printf("Failed to copy PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy PDF"})
			return
		}

		// Save areas to JSON file
		areasJSON, err := json.MarshalIndent(req.Areas, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process areas"})
			return
		}

		// Use consistent file name for areas JSON
		areasPath := filepath.Join(outputDir, "areas.json")
		if err := os.WriteFile(areasPath, areasJSON, 0644); err != nil {
			log.Printf("Failed to write areas file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save areas"})
			return
		}

		// Run the Python script to highlight the PDF
		pythonScriptPath := filepath.Join("internal", "backend", "highlight", "highlight_pdf.py")
		highlightedPDFPath := filepath.Join(outputDir, baseName+"_highlighted.pdf")

		cmd := exec.Command("python3", pythonScriptPath, copiedPDFPath, areasPath, highlightedPDFPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing Python script: %v\nOutput: %s", err, string(output))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to highlight PDF"})
			return
		}

		// Return the path to the highlighted PDF
		c.JSON(http.StatusOK, HighlightResponse{
			HighlightedPDFPath: highlightedPDFPath,
			Message:            "PDF highlighted successfully",
		})
	})

	// r.Static("/storage", "./storage")

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

		newURL := os.Getenv("AI_URL")
		if newURL == "" {
			newURL = "http://localhost:8081"
		}
		// newURL := "https://aiservice-service.onrender.com"

		areas, err := highlight.FindAreas(pdfpath, jobTitle, jobDetails, newURL, evaluationReference)
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
		areasPath := filepath.Join("storage", "uploads", "areas.json")
		// Write with UTF-8 encoding
		if err := os.WriteFile(areasPath, areasJSON, 0644); err != nil {
			log.Printf("Failed to write areas file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save areas"})
			return
		}

		// Create highlight client and highlight the PDF
		URL := os.Getenv("HIGHLIGHT_URL")
		if URL == "" {
			URL = "http://localhost:8080" // Default URL for local testing
		}
		// URL := "https://highlight-service.onrender.com"
		highlightClient := highlight.NewClient(URL)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback khi chạy local
	}

	r.Run(":" + port)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
