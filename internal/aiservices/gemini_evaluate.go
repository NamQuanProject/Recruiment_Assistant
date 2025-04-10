package aiservices

import (
	"encoding/json"
	"fmt"
)

func GeminiEvaluateScoring(jobType string, mainCategory string, CV string) (map[string]any, error) {
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

	mainCategoryStr := mainCategory
	// subCategoryStr := strings.Join(subCategory, ", ")

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	defer agent.Close()

	flexibleGuide := `
	🧠 Flexible Evaluation Guide:

	Score and explain each category using a comprehensive explanation — include evidence found, strengths, weaknesses, and relevance to the job.

	🔹 Main Categories (Score: 1–10)
	1. Technical Skills
	- 9–10: Deep, modern, and varied technical knowledge with real-world application and relevance to the role.
	- 6–8: Strong but slightly limited or generic skills.
	- 3–5: Partial or outdated skills, unclear usage.
	- 1–2: No evidence or unrelated tools.

	2. Work Experience
	- 9–10: Relevant roles, clear impact, progression, and achievements.
	- 6–8: Moderate experience or roles that are somewhat related.
	- 3–5: Weak or limited scope, missing impact.
	- 1–2: Little or no applicable experience.

	3. Education
	- 9–10: Prestigious or deeply relevant degrees.
	- 6–8: Standard education in relevant fields.
	- 3–5: General education, not directly related.
	- 1–2: No educational evidence or unrelated study.

	🔸 Subcategories (Score: 1–5)
	1. Leadership
	- 5: Strong leadership roles with measurable outcomes.
	- 3–4: Some leadership, unclear scope.
	- 1–2: No leadership activity mentioned.

	2. Communication
	- 5: Evidence of public speaking, writing, collaboration.
	- 3–4: Some indirect mention or soft evidence.
	- 1–2: Not demonstrated.

	3. Problem Solving
	- 5: Examples of overcoming challenges, innovation, analysis.
	- 3–4: Some signs but vague.
	- 1–2: No signs of problem-solving.
	`

	finalPrompt := fmt.Sprintf(`
	You are an experienced recruiter for the field of "%s".

	🎯 Your task:
	Evaluate the following CV **fairly and objectively**, using only information in the document.

	- Provide scores for each main category (1–10) and subcategory (1–5)
	- Give a **comprehensive explanation** per category — highlighting strong areas, weak areas, missing elements, and alignment with the job.
	- Avoid any assumptions based on gender, name, race, religion, appearance, or background. Be absolutely unbiased.
	- Also, if the information provided in the CV has proof for it, then evaluate an authenticity score for the whole CV — this is the reliability point.

	%s

	📁 Main and Sub Categories: %s  

	📄 Candidate CV:
	"""%s"""

	📋 Output:
	Return a single valid JSON object formatted like this:
	%s
`, jobType, flexibleGuide, mainCategoryStr, CV, structurePrompt)

	agent.Model.ResponseMIMEType = "application/json"

	resp := agent.CallChatGemini(finalPrompt)

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}
