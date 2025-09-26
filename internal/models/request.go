package models

// Request payload for dynamic recommendations
type Request struct {
	SystemPrompt               string         `json:"system_prompt"`
	UserPrompt                 string         `json:"user_prompt"`
	Categories                 []string       `json:"categories"`
	RecommendationsPerCategory int            `json:"recommendations_per_category"`
	ResponseFormat             map[string]any `json:"response_format"`
	ClientContext              map[string]any `json:"client_context"`
}

// type HairProfile struct {
// 	CurlType       string `json:"curl_type"`
// 	Porosity       string `json:"porosity"`
// 	Volume         string `json:"volume"`
// 	DesiredOutcome string `json:"desired_outcome"`
// }
