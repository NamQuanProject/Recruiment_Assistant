package aiservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GeminiEvaluateScoring(jobType string, mainCategory []string, subCategory []string, CV string) (map[string]any, error) {
	// Load the target structure for evaluation
	structure, jsonErr := ReadJsonStructure("./internal/aiservices/evaluate_structure.json")
	if jsonErr != nil {
		return nil, jsonErr
	}

	structureBytes, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return nil, err
	}
	structurePrompt := string(structureBytes)

	mainCategoryStr := strings.Join(mainCategory, ", ")
	subCategoryStr := strings.Join(subCategory, ", ")

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	// flexibleGuide := `
	// 🧠 Flexible Evaluation Guide:

	// Score and explain each category using a comprehensive explanation — include evidence found, strengths, weaknesses, and relevance to the job.

	// 🔹 Main Categories (Score: 1–10)
	// 1. Category related to skills
	// - 9–10: Deep, modern, and varied technical knowledge with real-world application and relevance to the role.
	// - 6–8: Strong but slightly limited or generic skills.
	// - 3–5: Partial or outdated skills, unclear usage.
	// - 1–2: No evidence or unrelated tools.

	// Marking:

	// `

	finalPrompt := fmt.Sprintf(`
	You are an experienced recruiter for the field of "%s".

	🎯 Your task:
	Evaluate the following CV **fairly and objectively**, using only information in the document.

	- Provide scores for each main category (1–10) and subcategory (1–5)
	- Give a **comprehensive explanation** per category — highlighting strong areas, weak areas, missing elements, and alignment with the job.
	- Avoid any assumptions based on gender, name, race, religion, appearance, or background. Be absolutely unbiased.
	- Also, if the information provided in the CV has proof for it, then evaluate an authenticity score for the whole CV — this is the reliability point.



	📁 Main Categories: %s  
	📂 Subcategories: %s

	📄 Candidate CV:
	"""%s"""

	📋 Output:
	Return a single valid JSON object formatted like this:
	%s
`, jobType, mainCategoryStr, subCategoryStr, CV, structurePrompt)

	agent.Model.ResponseMIMEType = "application/json"

	resp := agent.CallChatGemini(finalPrompt)

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}
