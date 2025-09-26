package models

type OrchestratorParams struct {
	SystemPrompt               string         `json:"system_prompt"`
	UserPrompt                 string         `json:"user_prompt"`
	Categories                 []string       `json:"categories"`
	RecommendationsPerCategory int            `json:"recommendatins_per_category"`
	ResponseFormat             map[string]any `json:"response_formatpackage agents"`
	Context                    map[string]any `json:"context"`
}
