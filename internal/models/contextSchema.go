package models

type ContextSchema struct {
	Name    string         `json:"name"`
	Content map[string]any `json:"content"`
}
