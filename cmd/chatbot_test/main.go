package main

import (
	"fmt"
	"log"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
)

func main() {
	evalID := "20250410_144613"

	config := aiservices.Config{
		APIKey:      "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok", // Replace with env var in production!
		ModelName:   "gemini-2.0-flash",
		Temperature: 0.0,
		Name:        "DefaultAgent",
	}

	factory := &aiservices.AgentFactory{Config: config}
	chatbot := aiservices.GetChatBot(evalID, factory)

	resp, err := chatbot.Ask("1", "List all asked questions please.")
	if err != nil {
		log.Fatalf("Ask failed: %v", err)
	}

	fmt.Printf("Chatbot response: %s\n", resp)

	// Save conversation history to file
	if err := chatbot.SaveHistoryToFile(); err != nil {
		log.Printf("Failed to save history: %v", err)
	} else {
		fmt.Println("History saved successfully.")
	}
}
