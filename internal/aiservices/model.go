package aiservices

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

type History struct {
	Question string
	Response string
	Date     string
}

type AIAgent struct {
	Id            string
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

type PersonalInfo struct {
	Name     string
	Age      int
	Location string
	Phone    string
	Email    string
	Website  string
	Github   string
	Linkedin string
	Facebook string
}

type EducationBackground struct {
	University      string
	Degree          string
	Major           string
	GraduationDate  string
	CurrentGPA      float32
	RelevantCourses []string
	Award           []string
	Scholarships    []string
}

type WorkExperience struct {
	CompanyName      string
	JobTitle         string
	JobDescription   string
	TypeJob          string
	Location         string
	StartDate        string
	EndDate          string
	Description      string
	Responsibilities string
	Skills           []Skills
	Accomplishments  []string
	Projects         []string
	Technologies     []string
	References       []string
}

type Skills struct {
	SkillName         string
	SkillLevel        string
	YearsOfExperience int
}

type Project struct {
	ProjectName string
	Description string
	Reference   string
}

type CV struct {
	JobApply            string
	PersonalInfo        PersonalInfo
	EducationBackground []EducationBackground
	WorkExperience      []WorkExperience
	WorkingSkills       []Skills
	Projects            []Project
	Activitys           []string
}

// AnalyzeCVRequest represents the request to analyze areas in a CV
type AnalyzeCVRequest struct {
	CVPath     string `json:"cv_path" binding:"required"`
	JobTitle   string `json:"job_title" binding:"required"`
	JobDetails string `json:"job_details" binding:"required"`
	TextBlocks []struct {
		Text   string  `json:"text"`
		Page   int     `json:"page"`
		X      float64 `json:"x"`
		Y      float64 `json:"y"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"text_blocks"`
	EvaluationReference map[string]any `json:"evaluation_reference"`
}

// Area represents an area in the CV that needs to be highlighted
type Area struct {
	Text        string  `json:"text"`
	Page        int     `json:"page"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "weak" or "strong"
}

// AnalyzeCVResponse represents the response from the AI server
type AnalyzeCVResponse struct {
	Areas []Area `json:"areas"`
}

type JD struct {
	JobTitle               string
	JobType                string
	ExperienceRequirements string
	SkillsRequired         []string
	Location               string
}
