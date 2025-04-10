package aiservices

import (
	"context"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func NewAIAgent(config Config, default_agent bool) (*AIAgent, error) {
	if default_agent {
		config.APIKey = "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok"
		config.ModelName = "gemini-2.0-flash"
		config.Temperature = 0.0
		config.Name = "DefaultAgent"
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	model := client.GenerativeModel(config.ModelName)

	setting := DefaultSafetySettings()
	history := []History{}

	return &AIAgent{
		Name:          config.Name,
		Client:        client,
		Model:         model,
		SafetySetting: setting,
		History:       history,
		APIKey:        config.APIKey,
		ModelName:     config.ModelName,
		MaxTokens:     config.MaxTokens,
		Temperature:   config.Temperature,
		ctx:           ctx,
	}, nil
}

func GetAIAgent(id string, config Config) (*AIAgent, error) {
	config.APIKey = "AIzaSyB22ThtcCvZuXual9uaT_6v4Bo5R6oBdok"
	config.ModelName = "gemini-2.0-flash"
	config.Temperature = 0.0
	config.Name = "DefaultAgent"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel(config.ModelName)
	setting := DefaultSafetySettings()
	history := HandleHistoryGet(id)

	return &AIAgent{
		Name:          config.Name,
		Client:        client,
		Model:         model,
		SafetySetting: setting,
		History:       history,
		APIKey:        config.APIKey,
		ModelName:     config.ModelName,
		MaxTokens:     config.MaxTokens,
		Temperature:   config.Temperature,
		ctx:           ctx,
	}, nil
}
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	default:
		return ""
	}
}
func HandleHistoryGet(id string) []History {
	historyData, err := ReadJsonStructure("./internal/aiservices/data/history.json")
	if err != nil {
		return []History{}
	}

	currentModelHistory, ok := historyData[id]
	if !ok {
		return []History{}
	}

	rawEntries, ok := currentModelHistory.([]interface{})
	if !ok {
		return []History{}
	}

	finalResult := make([]History, 0, len(rawEntries))
	for _, raw := range rawEntries {
		entryMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		finalResult = append(finalResult, History{
			Question: toString(entryMap["Question"]),
			Response: toString(entryMap["Response"]),
			Date:     toString(entryMap["Date"]),
		})
	}

	return finalResult
}

// SETTINGS AND INITIALIZATION

func DefaultSafetySettings() []*genai.SafetySetting {
	safety := []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	return safety
}

func (agent *AIAgent) SetSafetySettings(settings []*genai.SafetySetting) {
	agent.SafetySetting = settings
}

func (a *AIAgent) GetName() string {
	return a.Name
}
func (a *AIAgent) GetModel() string {
	return a.ModelName
}
func (a *AIAgent) GetMaxTokens() int {
	return a.MaxTokens
}
func (a *AIAgent) GetTemperature() float32 {
	return a.Temperature
}

func (agent *AIAgent) SetName(name string) {
	agent.Name = name
}
func (agent *AIAgent) SetAPIKey(apiKey string) {
	agent.APIKey = apiKey
	agent.Client, _ = genai.NewClient(agent.ctx, option.WithAPIKey(apiKey))
	agent.Model = agent.Client.GenerativeModel(agent.ModelName)
}

func (agent *AIAgent) SetModel(modelName string) {
	agent.ModelName = modelName
	agent.Model = agent.Client.GenerativeModel(modelName)
}
func (agent *AIAgent) SetTemperature(temp float32) {
	agent.Temperature = temp
}

func (agent *AIAgent) SetMaxTokens(maxTokens int) {
	agent.MaxTokens = maxTokens
}

func (agent *AIAgent) Close() {
	if agent.Client != nil {
		agent.Client.Close()
	}
}
