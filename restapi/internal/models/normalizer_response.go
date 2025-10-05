package models

type NormalizerResponse struct {
	Recommendations map[string]any `json:"recommendations"`
}
