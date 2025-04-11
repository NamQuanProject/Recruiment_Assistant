package aiservices

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type AgentFactory struct {
	Config Config
}

func (f *AgentFactory) CreateAgent(cvID string) (*AIAgent, error) {
	log.Printf("Creating an Agent with ID: %s.\n", cvID)
	agent, err := GetAIAgent(cvID, f.Config)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	agent.Id = cvID

	return agent, nil
}

type ChatBot struct {
	AgentFactory *AgentFactory
	evaluationID string
	agentsCache  map[string]*AIAgent // stores initialized agents
	mu           sync.Mutex          // for thread safety
}

func (cb *ChatBot) AddAgent(cv_id string, agent *AIAgent) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.agentsCache == nil {
		cb.agentsCache = make(map[string]*AIAgent)
	}

	agent.Model.ResponseMIMEType = "text/plain"

	cb.agentsCache[cv_id] = agent
	log.Printf("Added AI agent with ID: %s to Chatbot", cv_id)
	// cb.SaveHistoryToFile()
}

func (cb *ChatBot) Ask(cvID string, question string) (string, error) {
	log.Printf("Asking \"%s\" on agent ID: \"%s\".\n", question, cvID)
	cb.mu.Lock()
	agent, ok := cb.agentsCache[cvID]
	cb.mu.Unlock()

	// Lazy load agent if not cached
	if !ok {
		var err error
		agent, err = cb.AgentFactory.CreateAgent(cvID)
		if err != nil {
			log.Print(err)
			return "", err
		}
		cb.mu.Lock()
		cb.agentsCache[cvID] = agent
		cb.mu.Unlock()
	} else {
		log.Printf("Chatbot for cv_id: %s exists", cvID)
	}

	historyStr := agent.GetHistory()

	evalStr, err := EvaluationToString(cb.evaluationID, cvID)

	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	evalStr = fmt.Sprintf("You are an AI assistant analyzing a candidate's CV. This is your evaluation of a CV structured into categories:\n\"%s\"\nI will ask you several questions relating to the CV\n", evalStr)
	historyStr = fmt.Sprintf("This is the history of our chat:\n\"%s\"\n", historyStr)

	constructedPrompt := fmt.Sprintf(
		"%s%sPlease give concise answer together with explanation (appropriate length) for my question in plain text format: %s",
		evalStr,
		historyStr,
		question,
	)

	result := agent.CallChatBotGemini(constructedPrompt)

	response, ok := result["Response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response from agent")
	}

	agent.AddToHistory(question, response)
	if err := cb.SaveHistoryToFile(); err != nil {
		log.Printf("Failed to save history: %v", err)
	} else {
		fmt.Println("History saved successfully.")
	}

	return response, nil
}

func GetChatBot(evalID string, factory *AgentFactory) (*ChatBot, error) {
	log.Printf("Getting ChatBot with evalID: %s.\n", evalID)
	cb := ChatBot{
		AgentFactory: factory,
		evaluationID: evalID,
		agentsCache:  make(map[string]*AIAgent),
	}

	err := cb.loadExistingHistories()

	return &cb, err
}

func (cb *ChatBot) loadExistingHistories() error {
	historyFolder := filepath.Join("storage", fmt.Sprintf("evaluation_%s", cb.evaluationID), "agents_history")
	files, err := os.ReadDir(historyFolder)
	if err != nil {
		if os.IsNotExist(err) {
			// Folder doesn't exist — no histories to load
			return nil
		}
		log.Printf("Error reading storage folder: %v", err)
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Extract cvID from filename: e.g., "agent_123.json" → "123"
		cvID := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "agent_"), ".json")

		// Load history
		history, err := LoadHistoryFromFile(path.Join(historyFolder, file.Name()))
		if err != nil {
			log.Printf("Failed to load history for %s: %v", cvID, err)
			continue
		}

		// Create agent and cache it
		agent, err := cb.AgentFactory.CreateAgent(cvID)
		if err != nil {
			log.Printf("Failed to create agent %s: %v", cvID, err)
			continue
		}

		agent.History = history
		cb.agentsCache[cvID] = agent
	}

	return nil
}

func (cb *ChatBot) SaveHistoryToFile() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	historyFolder := filepath.Join("storage", fmt.Sprintf("evaluation_%s", cb.evaluationID), "agents_history")

	log.Print("Saving Chatbot History to file...\n")
	// Create folder if it doesn't exist
	if err := os.MkdirAll(historyFolder, 0755); err != nil {
		return fmt.Errorf("failed to create history folder: %w", err)
	}

	for cvID, agent := range cb.agentsCache {
		historyFile := filepath.Join(historyFolder, fmt.Sprintf("agent_%s.json", cvID))
		data, err := json.MarshalIndent(agent.History, "", "  ")
		if err != nil {
			log.Printf("failed to marshal history for agent %s: %v", cvID, err)
			continue
		}
		if err := os.WriteFile(historyFile, data, 0644); err != nil {
			log.Printf("failed to write history file for agent %s: %v", cvID, err)
			continue
		}
	}

	return nil
}
