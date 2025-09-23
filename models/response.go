package models

type RecommendationResponse struct {
	Recommendaton map[string][]string `json:"recommendation"`
}
