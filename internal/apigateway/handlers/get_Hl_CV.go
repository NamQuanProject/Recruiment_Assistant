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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing index"})
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	// Get the current path from the file

	currentPathBytes, err := os.ReadFile("storage/current.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read storage/current.txt"})
		return
	}
	basePath := strings.TrimSpace(string(currentPathBytes))
	// opne basePath/finaloutput.json
	finalOutputPath := filepath.Join(basePath, "finaloutput.json")
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

	// Return the paths as a response
	c.JSON(http.StatusOK, gin.H{"pathToCV": pathtocv, "pathToEval": pathtoeval})
}
