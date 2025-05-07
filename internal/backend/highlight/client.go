package highlight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents a client for the highlight server
type Client struct {
	BaseURL string
}

// NewClient creates a new highlight client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// HighlightPDF sends a request to highlight a PDF with areas
func (c *Client) HighlightPDF(pdfPath string, areas []Area) (*HighlightResponse, error) {
	// Create the request body
	reqBody := HighlightRequest{
		PDFPath: pdfPath,
		Areas:   areas,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.BaseURL+"/highlight", bytes.NewBuffer(jsonData))
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
	var highlightResp HighlightResponse
	if err := json.Unmarshal(body, &highlightResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &highlightResp, nil
}
