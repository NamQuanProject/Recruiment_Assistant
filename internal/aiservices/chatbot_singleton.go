// chatbot_singleton.go

package aiservices

import (
	"fmt"
	"log"
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

	// Avoid reinitializing if already initialized for the same evaluationID
	if currentChatbot != nil && currentChatbot.evaluationID == evaluationID {
		log.Println("Current Chatbot for evaluationID has already been initialized")
		return nil
	}

	log.Println("Initializing chatbot for evaluation ID:", evaluationID)

	config := Config{
		APIKey:      "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok", // TODO: Use env var in production
		ModelName:   "gemini-2.0-flash",
		Temperature: 0.0,
		Name:        "DefaultAgent",
	}

	factory := &AgentFactory{
		Config: config,
	}

	cb, err := GetChatBot(evaluationID, factory)
	if err != nil {
		return fmt.Errorf("failed to get chatbot: %w", err)
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
