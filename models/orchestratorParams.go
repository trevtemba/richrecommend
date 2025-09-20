package models

type OrchestratorParams struct {
	SystemPrompt               string
	UserPrompt                 string
	Categories                 []string
	RecommendationsPerCategory int
	RecommendationSchema       map[string]interface{}
}
