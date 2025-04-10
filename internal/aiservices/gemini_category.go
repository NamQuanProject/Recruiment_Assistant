package aiservices

import (
	"encoding/json"
	"fmt"
	"strings"
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
	You are an expert HR analyst working for a company in the field of """` + job_type + `""" recruiting.
	
	You are given:
	1. The company's official job description (main JD)
	2. A guided job suggestion (candidate JD)
	
	Your task is to extract the most relevant categories for evaluation and organize them into:
	- Main Categories: absolutely essential for the role
	- Subcategories: desirable or supporting qualities
	
	Each category must include:
	- A name
	- A scoring scale (1–10 for main, 1–5 for sub)
	- A comprehensive explanation:
	  - Why it's important
	  - How it relates to the job
	  - If it came from the main JD or candidate JD
	
	Make sure your analysis is fair and unbiased. Do NOT include anything based on gender, age, race, or personal traits. Evaluate purely based on job content.
	
	Official JD:
	"""` + main_jd + `"""
	
	Suggested JD:
	"""` + accountData + `"""
	
	The output format must strictly follow the following JSON structure:
	` + structurePrompt + `
	
	Return only a single top-level JSON object.
	`

	resp := agent.CallChatGemini(finalStructurePrompt)
	agent.Model.ResponseMIMEType = "application/json"

	fmt.Println("Parsed Response:", resp)
	return resp, nil
}

func HandleCategoryPrompt(structure_jd map[string]any) string {
	description := ""
	if val, ok := structure_jd["description"].(string); ok {
		description = val
	}

	objectivesStr := ""
	if jobDesc, ok := structure_jd["job_description"].(map[string]any); ok {
		if objectives, ok := jobDesc["Objectives of this role"].([]any); ok && len(objectives) > 0 {
			objectivesStr = " This role involves " + objectives[0].(string)
			for i := 1; i < len(objectives); i++ {
				objectivesStr += ", " + objectives[i].(string)
			}
			objectivesStr += "."
		}
	}

	skillsStr := ""
	if skills, ok := structure_jd["skills_requirements"].([]any); ok && len(skills) > 0 {
		skillsStr = " Ideal candidates should possess skills such as " + skills[0].(string)
		for i := 1; i < len(skills); i++ {
			skillsStr += ", " + skills[i].(string)
		}
		skillsStr += "."
	}

	insightStr := ""
	if questions, ok := structure_jd["interview_questions"].([]any); ok && len(questions) > 0 {
		insightStr = " During the interview, candidates are often asked about topics like "
		first := true
		for _, q := range questions {
			if qmap, ok := q.(map[string]any); ok {
				if questionText, ok := qmap["question"].(string); ok {
					if !first {
						insightStr += ", "
					}
					insightStr += strings.ToLower(questionText)
					first = false
				}
			}
		}
		insightStr += ". This reflects the importance of practical knowledge and strong communication skills in this role."
	}

	fullString := description + objectivesStr + skillsStr + insightStr
	return fullString
}
