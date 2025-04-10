package aiservices

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadJsonStructure(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteJsonStructure(filename string, data map[string]any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func LoadHistoryFromFile(filePath string) ([]History, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var history []History
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func EvaluationToString(evaluationID, cvID string) (string, error) {
	evaluationPath := filepath.Join("storage", fmt.Sprintf("evaluation_%s", evaluationID), "evaluation", cvID+".json")

	data, err := os.ReadFile(evaluationPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	evals, ok := raw["Evaluation"].([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid or missing 'Evaluation' field")
	}

	authenticity, ok := raw["Authenticity"]
	if !ok {
		return "", fmt.Errorf("invalid or missing 'Authenticity' field")
	}

	var builder strings.Builder
	for _, item := range evals {
		evalMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		category, _ := evalMap["category"].(string)
		score, _ := evalMap["score"].(string)
		explanation, _ := evalMap["explanation"].(string)

		builder.WriteString(fmt.Sprintf("Category: %s\n", category))
		builder.WriteString(fmt.Sprintf("Score: %s\n", score))
		builder.WriteString(fmt.Sprintf("Explanation: %s\n\n", explanation))
	}

	builder.WriteString(fmt.Sprintf("Authenticity: %s\n\n", authenticity))
	return builder.String(), nil
}
