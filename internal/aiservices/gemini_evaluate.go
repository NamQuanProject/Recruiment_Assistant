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
You are an experienced recruiter in the field of "%s".

🎯 **Your Task:**  
Evaluate the following CV **fairly and objectively**, using only the information provided in the document.

- You must assign scores for each **main category (1–10)** and **subcategory (1–5)**.
- ⚠️ **IMPORTANT**: You must evaluate and score **every single category and subcategory** listed below. **Do not skip any.**
- For each category, provide a **comprehensive explanation** — highlighting strengths, weaknesses, missing elements, and how well the candidate aligns with the role.
- Be completely unbiased: do **not** make assumptions based on gender, name, race, religion, appearance, or background.
- If the CV contains verifiable evidence (e.g., links, certifications, official documents), include an **Authenticity Score (1–10)** to reflect the reliability of the CV.

📁 **Evaluation Categories:**  
%s

📄 **Candidate CV:**  
"""%s"""

📋 **Output Format:**  
Return a single valid JSON object in the following structure:  
%s
`, jobType, mainCategoryStr, CV, structurePrompt)

	agent.Model.ResponseMIMEType = "application/json"

	resp := agent.CallChatGemini(finalPrompt)

	// fmt.Println("Parsed Response:", resp)
	return resp, nil
}
