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

func HandleCategoryPrompt(structure_jd map[string]any) string {
	// Lấy mô tả chính nếu có
	description := ""
	if val, ok := structure_jd["description"].(string); ok {
		description = val
	}

	// Lấy mục tiêu chính trong phần job_description
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

	// Lấy skill requirements nếu có
	skillsStr := ""
	if skills, ok := structure_jd["skills_requirements"].([]any); ok && len(skills) > 0 {
		skillsStr = " Ideal candidates should possess skills such as " + skills[0].(string)
		for i := 1; i < len(skills); i++ {
			skillsStr += ", " + skills[i].(string)
		}
		skillsStr += "."
	}

	// Tận dụng interview_questions như là điểm nổi bật ứng viên cần chuẩn bị
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

	// Kết hợp các phần
	fullString := description + objectivesStr + skillsStr + insightStr
	return fullString
}
