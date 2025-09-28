package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trevtemba/richrecommend/internal/models"
)

func GenerateSchema(categories []string) map[string]any {
	schema := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"required":             categories,
		"additionalProperties": false,
		"DoNotReference":       true,
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

func GenerateSystemMessage(systemPrompt models.SystemPrompt, contextSchemaName string, categories []string, recsPerCategory int) (string, error) {

	var message strings.Builder

	message.WriteString(fmt.Sprintf("You are a %s for %s ", systemPrompt.Role, systemPrompt.Clientele))

	if systemPrompt.Domain != "" {
		message.WriteString(fmt.Sprintf("in %s ", systemPrompt.Domain))
	}

	categoryStr := strings.Join(categories, ", ")

	message.WriteString(fmt.Sprintf("Given a %s, recommend up to %d products per category (%s). ", contextSchemaName, recsPerCategory, categoryStr))
	message.WriteString("Return the response in structured JSON format.")

	return message.String(), nil
}

func GenerateUserMessage(contextSchema models.ContextSchema) (string, error) {

	var message strings.Builder
	var ctxSchemaStr strings.Builder

	for field, val := range contextSchema.Content {
		strVal, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("could not turn context value into a string")
		}
		ctxSchemaStr.WriteString(fmt.Sprintf("\n%s: %s", field, strVal))
	}

	message.WriteString(fmt.Sprintf("%s:%s", contextSchema.Name, ctxSchemaStr.String()))

	message.WriteString("\n\nRecommend products.")

	return message.String(), nil
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
