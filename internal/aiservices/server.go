package aiservices

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",              // local frontend
			"https://frontend-eqtg.onrender.com", // deployed frontend
		}, // Allow requests from your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                          // Headers exposed to the browser
		AllowCredentials: true,                                                // Allow cookies or authentication headers
	}))
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

	// Endpoint to analyze CV areas
	r.POST("/ai/analyze-cv-areas", func(c *gin.Context) {
		var req AnalyzeCVRequest

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

		// Get the prompt using the new function
		prompt := GetCVAnalysisPrompt(req.JobTitle, req.JobDetails, textBlocksStr, req.EvaluationReference)

		result := agent.CallChatGemini(prompt)
		response := result["Response"].(string)

		// Parse the response to extract areas
		areas, err := ParseAreasFromGeminiResponse(response)
		if err != nil {
			log.Printf("Error parsing CV areas: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse areas from AI response"})
			return
		}

		// Return the areas
		c.JSON(http.StatusOK, AnalyzeCVResponse{
			Areas: areas,
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
			EvaluationID       string `json:"evaluation_id"`
			CVID               string `json:"cv_id"`
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
		// subCategoryString := ""

		resp, err := GeminiQuieriaExtract(request.JobName, subCategoryString, request.CompanyDescription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract criteria"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"criteria": resp})
	})
	r.POST("/ai/evaluate", func(c *gin.Context) {
		type JDRequest struct {
			JobName        string `json:"job_name"`
			JDMainQuiteria string `json:"jd_main_quiteria"`
			CVRawText      string `json:"cv_raw_text"`
			EvaluationID   string `json:"evaluation_id"`
			CVID           string `json:"cv_id"`
		}
		var request JDRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		eval_id := request.EvaluationID
		cv_id := request.CVID

		err := InitChatBot(eval_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Chatbot"})
			return
		}

		resp, err := GeminiEvaluateScoring(request.JobName, request.JDMainQuiteria, request.CVRawText, cv_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract criteria"})
			return
		}

		cb, err := GetChatBotInstance()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Chatbot instance"})
			return
		}

		err = cb.SaveHistoryToFile()
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Chatbot History"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"evaluation": resp})
	})

	r.POST("/ai/chatbot/init", func(c *gin.Context) {
		type InitChatbotRequest struct {
			EvaluationID string `json:"eval_id"`
		}

		var request InitChatbotRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		evalID := request.EvaluationID
		if evalID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing eval_id"})
			return
		}

		err := InitChatBot(evalID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to initialize chatbot",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Chatbot initialized successfully",
			"evaluation_id": evalID,
		})
	})

	r.POST("/ai/chatbot/ask", func(c *gin.Context) {
		type ChatRequest struct {
			CV_ID    string `json:"cv_id"`
			Question string `json:"question"`
		}

		var request ChatRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		cvID := request.CV_ID
		question := request.Question

		if cvID == "" || question == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cv_id and question are required"})
			return
		}

		// Check if chatbot is initialized
		if currentChatbot == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Chatbot not initialized"})
			return
		}

		// Ask the chatbot
		resp, err := currentChatbot.Ask(cvID, question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   fmt.Sprintf("Failed to ask chatbot with ID: %s", cvID),
				"details": err.Error(),
			})
			return
		}

		// Success response
		c.JSON(http.StatusOK, gin.H{
			"answer": resp,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback khi chạy local
	}

	r.Run(":" + port)
}
