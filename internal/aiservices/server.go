package aiservices

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	r.GET("/ai", func(c *gin.Context) {
		agent, err := NewAIAgent(Config{}, true)
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
		agent.Close()
	})

	r.GET("/ai/jd_category/", func(c *gin.Context) {
		fmt.Println("Route /ai/category is hit")

		structure, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		jobData := make(map[string]interface{})

		for jobCategory, accountDataRaw := range structure {
			accountData, ok := accountDataRaw.(map[string]interface{})
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid job data format"})
				return
			}
			jobData[jobCategory] = accountData
		}

		c.JSON(http.StatusOK, jobData)
	})

	// New endpoint to find weak areas in a CV
	r.POST("/ai/find_weak_areas", func(c *gin.Context) {
		var req FindWeakAreasRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON"})
			return
		}

		if _, err := os.Stat(req.CVPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "CV file does not exist"})
			return
		}

		agent, err := NewAIAgent(Config{}, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI agent"})
			return
		}
		defer agent.Close()

		textBlocksStr := ""
		for _, block := range req.TextBlocks {
			textBlocksStr += fmt.Sprintf("Page %d: '%s' at position (%.2f, %.2f) with size %.2f x %.2f\n",
				block.Page, block.Text, block.X, block.Y, block.Width, block.Height)
		}

		// Construct the prompt for finding weak areas
		prompt := fmt.Sprintf(`Analyze the CV for the job title "%s" with the following job details: "%s".

		Here are the text blocks extracted from the CV:
		%s
		
		Identify weak areas in the CV that could be improved to better match the job requirements. For each weak area, provide:
		1. The exact text from the CV (must match one of the text blocks above)
		2. The page number where it appears
		3. The x, y coordinates and dimensions of the text on the page (must match the position of the text block)
		4. A description of why this area is weak and how it could be improved

		Return the results in the following JSON format:
		{
		"weak_areas": [
			{
			"text": "Weak area text",
			"page": 1,
			"x": 100,
			"y": 200,
			"width": 200,
			"height": 50,
			"description": "This area is weak because..."
			}
		]
		}

		Make sure to use the exact text and coordinates from the text blocks provided.`, req.JobTitle, req.JobDetails, textBlocksStr)

		result := agent.CallChatGemini(prompt)
		response := result["Response"].(string)

		weakAreas, err := ParseWeakAreasFromGeminiResponse(response)
		if err != nil {
			log.Printf("Error parsing weak areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse weak areas from AI response"})
			return
		}

		// Return the weak areas
		c.JSON(http.StatusOK, FindWeakAreasResponse{
			WeakAreas: weakAreas,
		})
	})

	r.POST("/ai/parsing", func(c *gin.Context) {
		var requestBody struct {
			JobRawText string `json:"job_raw_text"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		prompt := "Parse the following CV: " + requestBody.JobRawText

		parsed_response, err := GeminiParsingRawCVText(prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Question": prompt,
			"Response": parsed_response,
		})
	})

	r.GET("/ai/jd_category/:job_name", func(c *gin.Context) {
		jobName := c.Param("job_name")
		fmt.Printf("Route /ai/jd_category/%s is hit\n", jobName)

		structure, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		accountDataRaw, exists := structure[jobName]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job category not found"})
			return
		}

		accountData, ok := accountDataRaw.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid job data format"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Response": accountData,
		})
	})

	r.POST("/ai/jd_criteria", func(c *gin.Context) {
		type JDRequest struct {
			JobName            string `json:"job_name"`
			CompanyDescription string `json:"company_jd"`
		}

		var request JDRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		fmt.Printf("Route /ai/jd_quiteria/%s is hit\n", request.JobName)

		structure, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse job structure"})
			return
		}

		accountDataRaw, exists := structure[request.JobName]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job category not found"})
			return
		}

		accountData, ok := accountDataRaw.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid job data structure"})
			return
		}

		subCategoryString := HandleCategoryPrompt(accountData)

		resp, err := GeminiQuieriaExtract(request.JobName, subCategoryString, request.CompanyDescription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract criteria"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"criteria": resp})
	})

	r.Run(":8081")
}
