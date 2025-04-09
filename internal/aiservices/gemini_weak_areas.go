package aiservices

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ParseWeakAreasFromGeminiResponse parses the response from Gemini to extract weak areas
func ParseWeakAreasFromGeminiResponse(response string) ([]WeakArea, error) {
	jsonRegex := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := jsonRegex.FindString(response)

	if jsonMatch == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var result struct {
		WeakAreas []WeakArea `json:"weak_areas"`
	}

	err := json.Unmarshal([]byte(jsonMatch), &result)
	if err == nil && len(result.WeakAreas) > 0 {
		validWeakAreas := []WeakArea{}
		for _, area := range result.WeakAreas {
			if strings.TrimSpace(area.Text) == "" {
				continue
			}

			// Skip weak areas with invalid coordinates
			if area.X < 0 || area.Y < 0 || area.Width <= 0 || area.Height <= 0 {
				continue
			}

			validWeakAreas = append(validWeakAreas, area)
		}

		if len(validWeakAreas) > 0 {
			return validWeakAreas, nil
		}
	}

	log.Println("Failed to parse JSON or no valid weak areas found, trying regex extraction")

	weakAreas := []WeakArea{}

	textRegex := regexp.MustCompile(`Text:\s*([^\n]+)\s*Page:\s*(\d+)\s*X:\s*([\d.]+)\s*Y:\s*([\d.]+)\s*Width:\s*([\d.]+)\s*Height:\s*([\d.]+)\s*Description:\s*([^\n]+)`)
	matches := textRegex.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) >= 8 {
			weakArea := WeakArea{
				Text:        strings.TrimSpace(match[1]),
				Page:        parseInt(match[2]),
				X:           parseFloat(match[3]),
				Y:           parseFloat(match[4]),
				Width:       parseFloat(match[5]),
				Height:      parseFloat(match[6]),
				Description: strings.TrimSpace(match[7]),
			}

			// Skip weak areas with empty text
			if strings.TrimSpace(weakArea.Text) == "" {
				continue
			}

			// Skip weak areas with invalid coordinates
			if weakArea.X < 0 || weakArea.Y < 0 || weakArea.Width <= 0 || weakArea.Height <= 0 {
				continue
			}

			weakAreas = append(weakAreas, weakArea)
		}
	}

	if len(weakAreas) == 0 {
		log.Println("Could not extract any valid weak areas from the AI response")
		return nil, fmt.Errorf("could not extract any valid weak areas from the AI response")
	}

	return weakAreas, nil
}

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
