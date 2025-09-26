package models

type RecommendationParams struct {
	SystemPrompt               string         `json:"system_prompt"`
	UserPrompt                 string         `json:"user_prompt"`
	Categories                 []string       `json:"categories"`
	RecommendationsPerCategory int            `json:"recommendatins_per_category"`
	Context                    map[string]any `json:"context"`
}
