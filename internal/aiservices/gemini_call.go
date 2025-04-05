package aiservices

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
)

func (agent *AIAgent) CallChatGemini(prompt string) {
	structure := DefaultGeminiStructure()
	agent.SetOutputStructure(structure)

	resp, err := agent.Model.GenerateContent(agent.ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	if resp.Candidates[0].Content != nil {
		part := resp.Candidates[0].Content.Parts[0]

		if txt, ok := part.(genai.Text); ok {
			var result map[string]interface{}

			if err := json.Unmarshal([]byte(txt), &result); err != nil {
				log.Fatal(err)
			}

			prettyJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(prettyJSON))
		}
	}
}
