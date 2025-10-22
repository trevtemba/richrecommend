package models

type RecommendationParams struct {
	Provider                   string             `json:"provider"`
	Model                      string             `json:"model"`
	SystemPromptParams         SystemPromptParams `json:"system_prompt_params"`
	Categories                 []string           `json:"categories"`
	RecommendationsPerCategory int                `json:"recommendations_per_category"`
	ContextSchema              ContextSchema      `json:"context_schema"`
}
