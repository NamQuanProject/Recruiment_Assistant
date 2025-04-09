package aiservices

import (
	"fmt"
)

func GeminiEvaluate(job_type string, category []string, CV map[string]any) (map[string]any, error) {
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
	"""` + "" + `"""

	- Official main job description:
	"""` + "" + `"""

	Please extract the relevant information into the following JSON structure:
	` + "" + `

	Return only a single top-level JSON object called "CV".
	`

	resp := agent.CallChatGemini(finalStructurePrompt)
	agent.Model.ResponseMIMEType = "application/json"

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}
