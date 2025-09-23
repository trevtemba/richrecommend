package models

// Request payload for dynamic recommendations
type RecommendationRequest struct {
	SystemPrompt               string                 `json:"system_prompt"`
	UserPrompt                 string                 `json:"user_prompt"`
	RecommendationsCategories  []string               `json:"recommendations_categories"`
	RecommendationsPerCategory int                    `json:"recommendations_per_category"`
	ResponseSchema             map[string]interface{} `json:"response_schema"`
	ClientContext              map[string]interface{} `json:"client_context"`
}

// type HairProfile struct {
// 	CurlType       string `json:"curl_type"`
// 	Porosity       string `json:"porosity"`
// 	Volume         string `json:"volume"`
// 	DesiredOutcome string `json:"desired_outcome"`
// }
