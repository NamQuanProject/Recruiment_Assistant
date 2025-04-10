// chatbot_singleton.go

package aiservices

import (
	"fmt"
	"sync"
)

var (
	currentChatbot *ChatBot
	mu             sync.Mutex
)

// InitChatBot creates a new chatbot instance and replaces the current one.
func InitChatBot(evaluationID string) error {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Initializing chatbot for evaluation ID:", evaluationID)

	config := Config{
		APIKey:      "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok", // Replace with env var in production!
		ModelName:   "gemini-2.0-flash",
		Temperature: 0.0,
		Name:        "DefaultAgent",
	}

	factory := AgentFactory{
		Config: config,
	}

	cb, err := GetChatBot(evaluationID, &factory)
	if err != nil {
		return err
	}

	currentChatbot = cb
	return nil
}

// GetChatBotInstance returns the current chatbot instance.
func GetChatBotInstance() (*ChatBot, error) {
	mu.Lock()
	defer mu.Unlock()

	if currentChatbot == nil {
		return nil, fmt.Errorf("chatbot not initialized")
	}

	return currentChatbot, nil
}
