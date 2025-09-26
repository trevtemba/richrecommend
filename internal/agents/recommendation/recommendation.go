package recommendation

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/trevtemba/richrecommend/internal/helpers"
	"github.com/trevtemba/richrecommend/internal/models"
)

func GenerateWithParams(params models.RecommendationParams) (models.RecommendationResponse, error) {

	var recommendation models.RecommendationResponse
	var recommendationMap map[string][]string
	var ProductRecommendationResponseSchema map[string]any = helpers.GenerateSchema(params.Categories)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return recommendation, fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	ctx := context.Background()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "product_recommendations",
		Description: openai.String("Conditioners, shampoos, and leave-in conditioners that are recommended for the user's hair profile"),
		Schema:      ProductRecommendationResponseSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(params.SystemPrompt),
			openai.UserMessage(params.UserPrompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4oMini2024_07_18,
	})
	if err != nil {
		return recommendation, fmt.Errorf("could not initiate chat with ai: %w", err)
	}

	recommendationMap, err = helpers.ParseChatResponse(chat.Choices[0].Message.Content, params.Categories)

	if err != nil {
		return recommendation, fmt.Errorf("could not parse chat response: %w", err)
	}

	recommendation.Recommendation = recommendationMap

	return recommendation, nil
}
