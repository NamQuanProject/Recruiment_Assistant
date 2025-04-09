package highlight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FindWeakAreasRequest represents the request to find weak areas in a CV
type FindWeakAreasRequest struct {
	CVPath     string        `json:"cv_path"`
	JobTitle   string        `json:"job_title"`
	JobDetails string        `json:"job_details"`
	TextBlocks []PDFTextBlock `json:"text_blocks"`
}

// FindWeakAreasResponse represents the response from the AI server
type FindWeakAreasResponse struct {
	WeakAreas []WeakArea `json:"weak_areas"`
}

// FindWeakAreas finds weak areas in a CV by calling the AI server
func FindWeakAreas(cvPath, jobTitle, jobDetails, aiServerURL string) ([]WeakArea, error) {
	// Extract text from the PDF
	textBlocks, err := ExtractTextFromPDF(cvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %v", err)
	}

	// Create the request body
	reqBody := FindWeakAreasRequest{
		CVPath:     cvPath,
		JobTitle:   jobTitle,
		JobDetails: jobDetails,
		TextBlocks: textBlocks,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", aiServerURL+"/ai/find_weak_areas", bytes.NewBuffer(jsonData))
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
	var findWeakAreasResp FindWeakAreasResponse
	if err := json.Unmarshal(body, &findWeakAreasResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return findWeakAreasResp.WeakAreas, nil
}	