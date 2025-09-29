package models

type RecommendationParams struct {
	Provider                   string        `json:"provider"`
	Model                      string        `json:"model"`
	SystemPrompt               SystemPrompt  `json:"system_prompt"`
	UserPrompt                 string        `json:"user_prompt"`
	Categories                 []string      `json:"categories"`
	RecommendationsPerCategory int           `json:"recommendatins_per_category"`
	ContextSchema              ContextSchema `json:"context_schema"`
}
