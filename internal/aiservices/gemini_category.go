package aiservices

import (
	"encoding/json"
	"fmt"
)

func GeminiQuieriaExtract(job_type string, sub_jd string, main_jd string) (map[string]any, error) {
	// STRUCTURE THE OUTPUT
	structure, err := ReadJsonStructure("./internal/aiservices/category_structure.json")
	if err != nil {
		return nil, err
	}
	structureBytes, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return nil, err
	}
	structurePrompt := string(structureBytes)

	// SUBCATEGORY
	fullCategory, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
	if err != nil {
		return nil, err
	}
	accountDataRaw, exists := fullCategory[job_type]
	if !exists {
		return nil, fmt.Errorf("job_type %s not found in jobs_guideds.json", job_type)
	}

	accountDataBytes, err := json.MarshalIndent(accountDataRaw, "", "  ")
	if err != nil {
		return nil, err
	}
	accountData := string(accountDataBytes)

	// CREATE AGENT
	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	finalStructurePrompt := `
	You are an expert in Human Resources of a company with deep experience recruiting in the field of """` + job_type + `"""

	You are highly skilled at reading two job descriptions:
	1. A suggested job description from us (candidate JD).
	2. The official job description from your company.

	Your task is to identify:
	- Subcategories that may be considered a plus.
	- Main categories that are absolutely required.

	You are given the following:
	- Extracted job description (candidate JD):
	"""` + accountData + `"""

	- Official main job description:
	"""` + main_jd + `"""

	Please extract the relevant information into the following JSON structure:
	` + structurePrompt + `

	Return only a single top-level JSON object called "CV".
	`

	resp := agent.CallChatGemini(finalStructurePrompt)
	agent.Model.ResponseMIMEType = "application/json"

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}
