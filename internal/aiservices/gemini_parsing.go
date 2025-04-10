package aiservices

import (
	"encoding/json"
	"fmt"
)

func GeminiParsingRawCVText(CVRawText string) (map[string]any, error) {
	structure, jsonErr := ReadJsonStructure("./internal/aiservices/parsing_structure.json")
	if jsonErr != nil {
		return nil, jsonErr
	}
	structureBytes, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return nil, err
	}
	structurePrompt := string(structureBytes)

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	// Step 3: Create detailed prompt
	FinalstructurePrompt := `
	You are an expert in Human Resources with deep experience recruiting in the field of Computer Science and Technology.

	You are highly skilled at reading CVs and extracting useful information to evaluate a candidate.
	
	You are given a .txt file containing extracted raw text from a CV. Here is the raw text from the CV:
	"""` + CVRawText + `"""

	Please extract the relevant information into the following JSON structure to format the response:
	` + string(structurePrompt) + `
	
	Please return only a single top-level JSON object called following the structure.
	`

	resp := agent.CallChatGemini(FinalstructurePrompt)
	agent.Model.ResponseMIMEType = "application/json"

	fmt.Println("Parsed Response:", resp)

	return resp, nil
}

func GeminiParsingRawJDText(CVRawText string) (map[string]any, error) {
	structure, jsonErr := ReadJsonStructure("./internal/aiservices/parsing_structure.json")
	if jsonErr != nil {
		return nil, jsonErr
	}
	structureBytes, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return nil, err
	}
	structurePrompt := string(structureBytes)

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	FinalstructurePrompt := `
	You are an expert in Human Resources with deep experience recruiting in the field of Computer Science and Technology.

	You are highly skilled at reading JDs and extracting useful information to evaluate a candidate.
	
	You are given a .txt file containing extracted raw text from a CV. Here is the raw text from the CV:
	"""` + CVRawText + `"""

	You MUST extract the relevant information into the following JSON structure to format the response:
	` + string(structurePrompt) + `
	
	Make sure to return only a single top-level JSON object following the structure.
	`

	resp := agent.CallChatGemini(FinalstructurePrompt)
	agent.Model.ResponseMIMEType = "application/json"

	fmt.Println("Parsed Response:", resp)

	return resp, nil
}
