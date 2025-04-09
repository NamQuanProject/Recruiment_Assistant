package aiservices

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ParseAreasFromGeminiResponse parses the response from Gemini to extract areas
func ParseAreasFromGeminiResponse(response string) ([]Area, error) {
	// First, try to extract JSON from the response
	jsonRegex := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := jsonRegex.FindString(response)
	
	if jsonMatch == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}
	
	// Try to parse the JSON
	var result struct {
		Areas []Area `json:"areas"`
	}
	
	err := json.Unmarshal([]byte(jsonMatch), &result)
	if err == nil && len(result.Areas) > 0 {
		// Validate the areas
		validAreas := []Area{}
		for _, area := range result.Areas {
			// Skip areas with empty text
			if strings.TrimSpace(area.Text) == "" {
				continue
			}
			
			// Skip areas with invalid coordinates
			if area.X < 0 || area.Y < 0 || area.Width <= 0 || area.Height <= 0 {
				continue
			}
			
			validAreas = append(validAreas, area)
		}
		
		if len(validAreas) > 0 {
			return validAreas, nil
		}
	}
	
	// If JSON parsing failed or no valid areas found, try to extract areas using regex
	log.Println("Failed to parse JSON or no valid areas found, trying regex extraction")
	
	// This is a fallback method that tries to extract areas from the text
	// It's less reliable but can work if the AI doesn't return proper JSON
	areas := []Area{}
	
	// Look for patterns like "Text: ... Page: ... X: ... Y: ... Width: ... Height: ... Description: ... Type: ..."
	textRegex := regexp.MustCompile(`Text:\s*([^\n]+)\s*Page:\s*(\d+)\s*X:\s*([\d.]+)\s*Y:\s*([\d.]+)\s*Width:\s*([\d.]+)\s*Height:\s*([\d.]+)\s*Description:\s*([^\n]+)\s*Type:\s*([^\n]+)`)
	matches := textRegex.FindAllStringSubmatch(response, -1)
	
	for _, match := range matches {
		if len(match) >= 9 {
			area := Area{
				Text:        strings.TrimSpace(match[1]),
				Page:        parseInt(match[2]),
				X:           parseFloat(match[3]),
				Y:           parseFloat(match[4]),
				Width:       parseFloat(match[5]),
				Height:      parseFloat(match[6]),
				Description: strings.TrimSpace(match[7]),
				Type:        strings.TrimSpace(match[8]),
			}
			
			// Skip areas with empty text
			if strings.TrimSpace(area.Text) == "" {
				continue
			}
			
			// Skip areas with invalid coordinates
			if area.X < 0 || area.Y < 0 || area.Width <= 0 || area.Height <= 0 {
				continue
			}
			
			// Set default type based on the description content
			if area.Type == "" {
				// Check if the description contains keywords indicating strength
				lowerDesc := strings.ToLower(area.Description)
				if strings.Contains(lowerDesc, "strong") || 
				   strings.Contains(lowerDesc, "excellent") || 
				   strings.Contains(lowerDesc, "good") || 
				   strings.Contains(lowerDesc, "impressive") || 
				   strings.Contains(lowerDesc, "relevant") {
					area.Type = "strong"
				} else {
					area.Type = "weak"
				}
			}
			
			areas = append(areas, area)
		}
	}
	
	if len(areas) == 0 {
		// If we still couldn't extract areas, log an error
		log.Println("Could not extract any valid areas from the AI response")
		return nil, fmt.Errorf("could not extract any valid areas from the AI response")
	}
	
	return areas, nil
}

// Helper functions to parse strings to numbers
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
} 