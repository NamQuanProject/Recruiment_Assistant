package aiservices

import (
	"encoding/json"
	"fmt"
	"log"
)

func GeminiEvaluateScoring(jobType, mainCategory, CV, cv_id string) (map[string]any, error) {
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

	// Get chatbot to update
	cb, err := GetChatBotInstance()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	agent, err := NewAIAgent(Config{}, true)
	if err != nil {
		return nil, err
	}
	// defer agent.Close()

	finalPrompt := fmt.Sprintf(`
	You are an experienced recruiter in the field of "%s".
	
	🎯 **Your Task:**  
	Evaluate the following CV **fairly and objectively**, using only the information provided in the document.
	
	- Assign a **score for EACH main category (1–10)** and for **each subcategory (1–5)** listed below.
	- ⚠️ **IMPORTANT**: You **must evaluate and score** *every single category and subcategory* listed in the "Evaluation Categories" section. **Do not skip any.**
	- Provide a **detailed explanation per category**: highlight candidate strengths, weaknesses, what's missing, and how well the CV aligns with the job.
	- Stay unbiased: **do NOT** make assumptions about gender, ethnicity, nationality, religion, or personal appearance.
	- If the CV provides **verifiable links or certifications**, give an **Authenticity Score (1–10)** explaining how trustworthy the information appears.
	- If the provide CV is not really a CV just return 0 for all score of all category with the explanation the provided CV is not a CV at all.
	
	📁 **Evaluation Categories:**  
	%s
	
	📄 **Candidate CV:**  
	"""%s"""
	
	📋 **Output Format:**  
	Return a **single JSON object** in the following structure:  
	%s
	`, jobType, mainCategoryStr, CV, structurePrompt)

	agent.Model.ResponseMIMEType = "application/json"

	resp := agent.CallChatGemini(finalPrompt)

	cb.AddAgent(cv_id, agent)
	// fmt.Println("Parsed Response:", resp)
	return resp, nil
}
