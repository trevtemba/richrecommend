package recommendation

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	recHelpers "github.com/trevtemba/richrecommend/internal/helpers/recommendation"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
)

// func GenerateWithBaseParams(params models.RecommendationParams, requestId string) (models.RecommendationResponse, error) {

// 	var recommendation models.RecommendationResponse
// 	var recommendationMap map[string][]string
// 	var ProductRecommendationResponseSchema map[string]any = recHelpers.GenerateSchema(params.Categories)

// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		return recommendation, fmt.Errorf("OPENAI_API_KEY not set")
// 	}

// 	client := openai.NewClient(
// 		option.WithAPIKey(apiKey),
// 	)
// 	ctx := context.Background()

// 	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
// 		Name:        "product_recommendations",
// 		Description: openai.String("Conditioners, shampoos, and leave-in conditioners that are recommended for the user's hair profile"),
// 		Schema:      ProductRecommendationResponseSchema,
// 		Strict:      openai.Bool(true),
// 	}

// 	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
// 		Messages: []openai.ChatCompletionMessageParamUnion{
// 			openai.SystemMessage("params.SystemPrompt"),
// 		},
// 		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
// 			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
// 		},
// 		Seed:  openai.Int(1),
// 		Model: openai.ChatModelGPT4oMini2024_07_18,
// 	})
// 	if err != nil {
// 		return recommendation, fmt.Errorf("could not initiate chat with ai: %w", err)
// 	}

// 	recommendationMap, err = recHelpers.ParseChatResponse(chat.Choices[0].Message.Content, params.Categories)

// 	if err != nil {
// 		return recommendation, fmt.Errorf("could not parse chat response: %w", err)
// 	}

// 	recommendation.Recommendation = recommendationMap

// 	return recommendation, nil
// }

func GenerateWithAdvParams(params models.RecommendationParams, key string, requestId string) (models.RecommendationResponse, error) {

	logger.Log(logger.LogTypeAgentStart, logger.LevelInfo, "Recommendation agent started", "request_id", requestId)
	var recommendation models.RecommendationResponse
	logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Generating product recommendation schema...", "request_id", requestId)
	var ProductRecommendationResponseSchema map[string]any = recHelpers.GenerateSchema(params.Categories)

	apiKey := key
	if apiKey == "" {
		return recommendation, fmt.Errorf("no key was provided")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	ctx := context.Background()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "product_recommendations",
		Description: openai.String("products that are recommended to the user given their client context"),
		Schema:      ProductRecommendationResponseSchema,
		Strict:      openai.Bool(true),
	}

	logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Generating system message...", "request_id", requestId)
	sysMsg, err := recHelpers.GenerateSystemMessage(params.SystemPromptParams, params.ContextSchema.Name, params.Categories, params.RecommendationsPerCategory)
	if err != nil {
		return recommendation, fmt.Errorf("could not generate system prompt: %w", err)
	}

	logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Generating user message...", "request_id", requestId)
	usrMsg, err := recHelpers.GenerateUserMessage(params.ContextSchema)
	if err != nil {
		return recommendation, fmt.Errorf("could not generate user prompt: %w", err)
	}

	logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Chatting with LLM...", "request_id", requestId)
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sysMsg),
			openai.UserMessage(usrMsg),
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

	logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Parsing LLM response...", "request_id", requestId)
	recommendation, err = recHelpers.ParseChatResponse(chat.Choices[0].Message.Content, params.Categories)
	if err != nil {
		return recommendation, fmt.Errorf("could not parse chat response: %w", err)
	}

	return recommendation, nil
}
