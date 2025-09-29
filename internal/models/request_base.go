package models

// Request payload for dynamic recommendations
type RequestBase struct {
	SystemPrompt               SystemPrompt  `json:"system_prompt"`
	UserPrompt                 string        `json:"user_prompt"`
	Categories                 []string      `json:"categories"`
	RecommendationsPerCategory int           `json:"recommendations_per_category"`
	ContextSchema              ContextSchema `json:"context_schema"`
	Include                    []string      `json:"include"`
}

// type HairProfile struct {
// 	CurlType       string `json:"curl_type"`
// 	Porosity       string `json:"porosity"`
// 	Volume         string `json:"volume"`
// 	DesiredOutcome string `json:"desired_outcome"`
// }
