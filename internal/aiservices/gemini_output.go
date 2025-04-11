package aiservices

import (
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
)

// OUTPUT SETTINGS
func (agent *AIAgent) SetOutputStructure(structure map[string]any) {
	agent.Model.ResponseMIMEType = "application/json"
	agent.Model.ResponseSchema = &genai.Schema{
		Type:       genai.TypeObject,
		Properties: StructureToProperties(structure),
	}
}

func StructureToProperties(structure map[string]any) map[string]*genai.Schema {
	properties := make(map[string]*genai.Schema)
	for key, value := range structure {
		switch v := value.(type) {
		case string:
			if v == "string" {
				properties[key] = &genai.Schema{Type: genai.TypeString}
			} else if v == "int" || v == "float" {
				properties[key] = &genai.Schema{Type: genai.TypeNumber}
			} else if v == "bool" {
				properties[key] = &genai.Schema{Type: genai.TypeBoolean}
			}
		case map[string]any:
			properties[key] = &genai.Schema{
				Type:       genai.TypeObject,
				Properties: StructureToProperties(v),
			}
		case []any:
			if len(v) > 0 {
				var itemSchema *genai.Schema

				switch firstItem := v[0].(type) {
				case string:
					itemSchema = &genai.Schema{Type: genai.TypeString}
				case int, int8, int16, int32, int64, float32, float64:
					itemSchema = &genai.Schema{Type: genai.TypeNumber}
				case bool:
					itemSchema = &genai.Schema{Type: genai.TypeBoolean}
				case map[string]any:
					itemSchema = &genai.Schema{
						Type:       genai.TypeObject,
						Properties: StructureToProperties(firstItem),
					}
				default:
					itemSchema = &genai.Schema{Type: genai.TypeString}
				}

				properties[key] = &genai.Schema{
					Type:  genai.TypeArray,
					Items: itemSchema,
				}
			} else {
				properties[key] = &genai.Schema{
					Type:  genai.TypeArray,
					Items: &genai.Schema{Type: genai.TypeString},
				}
			}
		}
	}
	return properties
}

func DefaultGeminiStructure() map[string]any {
	return map[string]any{
		"Answer": "string",
	}
}
func (agent *AIAgent) AddToHistory(question string, response string) {
	agent.History = append(agent.History, History{
		Question: question,
		Response: response,
		Date:     toString(time.Now()),
	})

	if len(agent.History) > 20 {
		first := agent.History[0]
		last := agent.History[len(agent.History)-19:] // keep last 4
		agent.History = append([]History{first}, last...)
	}
}

func (agent *AIAgent) GetHistory() string {
	// Convert the history to a string
	history_prompt := ""
	for _, history := range agent.History {
		history_prompt += fmt.Sprintf("Question: %s\nAnswer: %s\n", history.Question, history.Response)
	}
	return history_prompt
}

// func (agent *AIAgent) ClearHistory() {
// 	agent.History = []History{}
// }
// func (agent *AIAgent) SetHistory(history []History) {
// 	agent.History = history
// }
