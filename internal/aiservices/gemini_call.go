package aiservices

import (
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
)

func (agent *AIAgent) CallChatGemini(prompt string) map[string]any {

	history_prompt := ""
	for _, history := range agent.History {
		history_prompt += fmt.Sprintf("Question: %s\nAnswer: %s\n", history.Question, history.Response)
	}

	final_prompt := history_prompt + fmt.Sprintf("Question: %s\n", prompt)

	resp, err := agent.Model.GenerateContent(agent.ctx, genai.Text(final_prompt))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	// if resp.Candidates[0].Content != nil {
	// 	part := resp.Candidates[0].Content.Parts[0]

	// 	if txt, ok := part.(genai.Text); ok {
	// 		var result map[string]interface{}

	// 		if err := json.Unmarshal([]byte(txt), &result); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Println("Answer:", result["Answer"])
	// 	}
	// }
	string_ouput := resp.Candidates[0].Content.Parts[0].(genai.Text)

	// fmt.Println("Answer:", resp.Candidates[0].Content.Parts[0])

	agent.AddToHistory(prompt, string(string_ouput))

	return map[string]any{
		"Response": string(string_ouput),
	}
}

func (agent *AIAgent) CallChatBotGemini(prompt string) map[string]any {

	resp, err := agent.Model.GenerateContent(agent.ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	string_ouput := resp.Candidates[0].Content.Parts[0].(genai.Text)

	return map[string]any{
		"Response": string(string_ouput),
	}
}

func (agent *AIAgent) CallGeminiStructure(prompt string, structure map[string]any) map[string]any {
	agent.SetOutputStructure(structure)
	history_prompt := ""
	for _, history := range agent.History {
		history_prompt += fmt.Sprintf("Question: %s\nAnswer: %s\n", history.Question, history.Response)
	}

	final_prompt := "This is the history context:" + history_prompt + "\n" + "And this is the current problems" + fmt.Sprintf("Question: %s\n", prompt)
	// Set the final prompt

	resp, err := agent.Model.GenerateContent(agent.ctx, genai.Text(final_prompt))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	string_ouput := resp.Candidates[0].Content.Parts[0].(genai.Text)

	fmt.Println("Answer:", resp.Candidates[0].Content.Parts[0])

	agent.AddToHistory(prompt, string(string_ouput))

	return map[string]any{
		"Response": string(string_ouput),
	}

}
