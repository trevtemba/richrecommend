package models

type NormalizerResponse struct {
	Recommendations []map[string]map[string]any `json:"recommendations"`
	FailedProducts  []string                    `json:"failed_products"`
}
