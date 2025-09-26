package models

type RecommendationResponse struct {
	Recommendation map[string][]string `json:"recommendation"`
}
