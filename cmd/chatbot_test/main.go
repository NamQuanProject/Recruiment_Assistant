package main

import (
	"fmt"
	"log"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
)

func main() {
	evalID := "20250410_165023"
	cvID := "20250410_013723_0065"
	question := "List all questions that I asked you please."

	// config := aiservices.Config{
	// 	APIKey:      "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok", // Replace with env var in production!
	// 	ModelName:   "gemini-2.0-flash",
	// 	Temperature: 0.0,
	// 	Name:        "DefaultAgent",
	// }

	aiservices.InitChatBot(evalID)

	chatbot, err := aiservices.GetChatBotInstance()
	if err != nil {
		log.Fatalf("Get Chatbot failed: %v", err)
	}

	resp, err := chatbot.Ask(cvID, question)
	if err != nil {
		log.Fatalf("Ask failed: %v", err)
	}

	fmt.Printf("Chatbot response: %s\n", resp)

}
