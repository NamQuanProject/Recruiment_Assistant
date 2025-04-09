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

type JD struct {
	JobTitle               string
	JobType                string
	ExperienceRequirements string
	SkillsRequired         []string
	Location               string
}
