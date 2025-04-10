package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHlCVHandler(c *gin.Context) {

	//Get the integer index from the query parameter
	indexStr := c.Query("index")
	if indexStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing intdex"})
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid intdex"})
		return
	}

	// Get the current path from the file

	currentPathBytes, err := os.ReadFile("storage/current.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read storage/current.txt"})
		return
	}
	basePath := strings.TrimSpace(string(currentPathBytes))
	// read the file filename as a string form basepath/parse/jobname.txt
	// Construct the path to the jobname.txt file
	jobNamePath := filepath.Join(basePath, "parse", "jobname.txt")

	// Check if the jobname.txt file exists
	if _, err := os.Stat(jobNamePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "jobname.txt not found"})
		return
	}

	// Read the jobname.txt file
	jobNameBytes, err := os.ReadFile(jobNamePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read jobname.txt"})
		return
	}

	// Convert the file content to a string
	jobName := strings.TrimSpace(string(jobNameBytes))

	// opne basePath/finaloutput.json
	finalOutputPath := filepath.Join("internal", "backend", "output", "output.json")
	if _, err := os.Stat(finalOutputPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "finaloutput.json not found"})
		return
	}
	// Read the finaloutput.json file
	finalOutputFile, err := os.Open(finalOutputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open finaloutput.json"})
		return
	}
	defer finalOutputFile.Close()
	// Read the file content and save it into var list in json
	var list []string
	if err := json.NewDecoder(finalOutputFile).Decode(&list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode finaloutput.json"})
		return
	}
	// Check if the index is within bounds
	if index < 0 || index >= len(list) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index out of bounds"})
		return
	}

	// Get the item at the specified index
	item := list[index]

	// Parse the item to extract pathToCV and pathToEval
	var itemData struct {
		PathToCV   string `json:"pathToCV"`
		PathToEval string `json:"pathToEval"`
	}
	if err := json.Unmarshal([]byte(item), &itemData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot parse item data"})
		return
	}

	// Save the paths in variables
	pathtocv := itemData.PathToCV
	pathtoeval := itemData.PathToEval

	request := struct {
		JobTitle       string `json:"job_title"`
		JobDetailsPath string `json:"job_details_path"`
		PdfPath        string `json:"pdf_path"`
		EvalRefPath    string `json:"evaluation_path"`
	}{
		JobTitle:       jobName,
		JobDetailsPath: filepath.Join(basePath, "parse", "jd.txt"),
		PdfPath:        pathtocv,
		EvalRefPath:    pathtoeval,
	}

	//call server local:4000
	// Marshal the request into JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot marshal request body"})
		return
	}

	// Make a POST request to the server at localhost:4000
	resp, err := http.Post("http://localhost:4000", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call server"})
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server returned an error", "status": resp.StatusCode})
		return
	}

	// Decode the response body into a map
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode server response"})
		return
	}

	// Return the paths as a response
	c.JSON(http.StatusOK, gin.H{"highlighted_pdf_path": responseBody["highlighted_pdf_path"]})
}
