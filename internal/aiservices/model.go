package aiservices

import (
	"context"
	"time"

	"github.com/google/generative-ai-go/genai"
)

type History struct {
	Question string
	Response string
	Date     time.Time
}

type AIAgent struct {
	Name          string
	Client        *genai.Client
	Model         *genai.GenerativeModel
	SafetySetting []*genai.SafetySetting
	History       []History
	APIKey        string
	ModelName     string
	MaxTokens     int
	Temperature   float32
	ctx           context.Context
}

type Config struct {
	Name        string
	APIKey      string
	ModelName   string
	MaxTokens   int
	Temperature float32
}
