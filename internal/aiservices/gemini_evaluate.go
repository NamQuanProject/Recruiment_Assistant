package aiservices

import (
	"encoding/json"
	"fmt"
)

func GeminiEvaluateScoring(jobType string, mainCategory string, CV string) (map[string]any, error) {
	structure, jsonErr := ReadJsonStructure("./internal/aiservices/evaluate_structure.json")
	if jsonErr != nil {
		return nil, jsonErr
	}

	structureBytes, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return nil, err
	}
	structurePrompt := string(structureBytes)

	mainCategoryStr := mainCategory
	// subCategoryStr := strings.Join(subCategory, ", ")

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	finalPrompt := fmt.Sprintf(`
	You are an experienced recruiter for the field of "%s".

	🎯 Your task:
	Evaluate the following CV **fairly and objectively**, using only information in the document.

	- You must provide scores for each main category (1–10) and subcategory (1–5)
	- You must scoring with full of the category provided.
	- Give a **comprehensive explanation** per category — highlighting strong areas, weak areas, missing elements, and alignment with the job.
	- Avoid any assumptions based on gender, name, race, religion, appearance, or background. Be absolutely unbiased.
	- Also, if the information provided in the CV has proof for it, then evaluate an authenticity score for the whole CV — this is the reliability point.


	📁 Main and Sub Categories: %s  

	📄 Candidate CV:
	"""%s"""

	📋 Output:
	Return a single valid JSON object formatted like this:
	%s
`, jobType, mainCategoryStr, CV, structurePrompt)

	agent.Model.ResponseMIMEType = "application/json"

	resp := agent.CallChatGemini(finalPrompt)

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}
