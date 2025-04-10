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

func (cb *ChatBot) Ask(cvID string, question string) (string, error) {
	log.Printf("Asking %s on ID: %s.\n", question, cvID)
	cb.mu.Lock()
	agent, ok := cb.agentsCache[cvID]
	cb.mu.Unlock()

	// Lazy load agent if not cached
	if !ok {
		var err error
		agent, err = cb.AgentFactory.CreateAgent(cvID)
		if err != nil {
			fmt.Print(err)
			return "", err
		}
		cb.mu.Lock()
		cb.agentsCache[cvID] = agent
		cb.mu.Unlock()
	}

	// ✅ Step 1: Build prompt from history
	historyStr := agent.GetHistory()

	// ✅ Step 2: Add evaluation summary
	evalStr := ""
	// if agent.CVEvaluationSummary != "" {
	// 	evalStr = fmt.Sprintf("Here is the evaluation summary of the CV:\n%s\n\n", agent.CVEvaluationSummary)
	// }

	// ✅ Step 3: Add new question
	constructedPrompt := fmt.Sprintf(
		"%s%sUser question: %s",
		evalStr,
		historyStr,
		question,
	)

	// ✅ Step 4: Call Gemini with full prompt
	result := agent.CallChatGemini(constructedPrompt)

	response, ok := result["Response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response from agent")
	}

	// ✅ Step 5: Save to history
	agent.AddToHistory(question, response)

	return response, nil
}

func GetChatBot(evalID string, factory *AgentFactory) *ChatBot {
	log.Printf("Getting ChatBot with evalID: %s.\n", evalID)
	cb := ChatBot{
		AgentFactory: factory,
		evaluationID: evalID,
		agentsCache:  make(map[string]*AIAgent),
	}

	cb.loadExistingHistories()

	return &cb
}

func (cb *ChatBot) loadExistingHistories() {
	historyFolder := filepath.Join("storage", fmt.Sprintf("evaluation_%s", cb.evaluationID), "agents_history")
	files, err := os.ReadDir(historyFolder)
	if err != nil {
		if os.IsNotExist(err) {
			// Folder doesn't exist — no histories to load
			return
		}
		log.Printf("Error reading storage folder: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Extract cvID from filename: e.g., "agent_123.json" → "123"
		cvID := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "agent_"), ".json")

		// Load history
		history, err := loadHistoryFromFile(path.Join(historyFolder, file.Name()))
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
}

func loadHistoryFromFile(filePath string) ([]History, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var history []History
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (cb *ChatBot) SaveHistoryToFile() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	historyFolder := filepath.Join("storage", fmt.Sprintf("evaluation_%s", cb.evaluationID), "agents_history")

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
