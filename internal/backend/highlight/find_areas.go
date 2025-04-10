package highlight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type FindAreasRequest struct {
	CVPath              string        `json:"cv_path"`
	JobTitle            string        `json:"job_title"`
	JobDetails          string        `json:"job_details"`
	TextBlocks          []PDFTextBlock `json:"text_blocks"`
	EvaluationReference map[string]any  `json:"evaluation_reference"`
}


type FindAreasResponse struct {
	Areas []Area `json:"areas"`
}


func FindAreas(cvPath, jobTitle, jobDetails, aiServerURL string, evaluationReference map[string]any) ([]Area, error) {
	// Extract text from the PDF
	textBlocks, err := ExtractTextFromPDF(cvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %v", err)
	}

	// Create the request body
	reqBody := FindAreasRequest{
		CVPath:              cvPath,
		JobTitle:            jobTitle,
		JobDetails:          jobDetails,
		TextBlocks:          textBlocks,
		EvaluationReference: evaluationReference,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", aiServerURL+"/ai/analyze-cv-areas", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("request failed: %s", errorResp.Error)
	}

	// Unmarshal the response
	var findAreasResp FindAreasResponse
	if err := json.Unmarshal(body, &findAreasResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return findAreasResp.Areas, nil
}	