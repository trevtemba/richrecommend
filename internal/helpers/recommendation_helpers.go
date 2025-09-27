package helpers

import (
	"encoding/json"
	"fmt"
)

func GenerateSchema(categories []string) map[string]any {
	schema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   categories,
	}

	props := schema["properties"].(map[string]any)

	for _, cat := range categories {
		props[cat] = map[string]any{
			"type":        "array",
			"items":       map[string]any{"type": "string"},
			"description": fmt.Sprintf("list of recommended %s", cat),
		}
	}

	return schema
}

func ParseChatResponse(content string, categories []string) (map[string][]string, error) {
	var data map[string]any
	var recommendedItems = make(map[string][]string)

	err := json.Unmarshal([]byte(content), &data)

	if err != nil {
		return recommendedItems, fmt.Errorf("could not parse chat response: %w", err)
	}

	for _, cat := range categories {
		recommendedItems[cat] = nil
		if key, ok := data[cat]; ok {
			if itemArr, ok := key.([]any); ok {
				for _, item := range itemArr {
					if itemName, ok := item.(string); ok {
						recommendedItems[cat] = append(recommendedItems[cat], itemName)
					}
				}
			}
		}
	}
	return recommendedItems, nil
}
