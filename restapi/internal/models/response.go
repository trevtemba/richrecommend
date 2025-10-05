package models

type Response struct {
	Recommendaton map[string][]string `json:"recommendation"`
}
