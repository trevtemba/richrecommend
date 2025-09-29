package orchestrator

import (
	"github.com/trevtemba/richrecommend/internal/agents/recommendation"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
)

func RunAdvPipelineWithParams(params models.OrchestratorParams, key string, requestId string) (any, error) {
	// Step 1: Recommendation Agent

	var recommendationParams models.RecommendationParams

	recommendationParams.Provider = params.Provider
	recommendationParams.Model = params.Model
	recommendationParams.SystemPrompt = params.SystemPrompt
	recommendationParams.UserPrompt = params.UserPrompt
	recommendationParams.Categories = params.Categories
	recommendationParams.RecommendationsPerCategory = params.RecommendationsPerCategory
	recommendationParams.ContextSchema = params.ContextSchema

	recommendedProducts, err := recommendation.GenerateWithAdvParams(recommendationParams, key, requestId)
	if err != nil {
		logger.Log(logger.LogTypeAgentAbort, logger.LevelError, "Agent aborted due to error", "request_id", requestId)
		return nil, err
	}
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Recommendation agent finished", "request_id", requestId)
	// // Step 2: Scraper Agents
	// rawProducts := scraper.ScrapeProducts(recommendedProducts)

	// // Step 3: Normalizer Agent
	// normalized := normalizer.NormalizeProducts(rawProducts)

	// // Step 4: Transform to JSON-like structure based on categories
	// result := make(map[string][]map[string]any)
	// for _, p := range normalized {
	// 	if _, ok := result[p.Category]; !ok {
	// 		result[p.Category] = []map[string]any{}
	// 	}
	// 	result[p.Category] = append(result[p.Category], map[string]any{
	// 		"title": p.Title,
	// 		"price": p.Price,
	// 		"link":  p.Link,
	// 		"image": p.Image,
	// 	})
	// }s

	return recommendedProducts, nil
}

// func RunBasePipelineWithParams(params models.OrchestratorParams, requestId string) (any, error) {
// 	// Step 1: Recommendation Agent

// 	var recommendationParams models.RecommendationParams

// 	recommendationParams.SystemPrompt = params.SystemPrompt
// 	recommendationParams.UserPrompt = params.UserPrompt
// 	recommendationParams.Categories = params.Categories
// 	recommendationParams.RecommendationsPerCategory = params.RecommendationsPerCategory
// 	recommendationParams.ContextSchema = params.ContextSchema

// 	recommendedProducts, err := recommendation.GenerateWithBaseParams(recommendationParams, requestId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// // Step 2: Scraper Agents
// 	// rawProducts := scraper.ScrapeProducts(recommendedProducts)

// 	// // Step 3: Normalizer Agent
// 	// normalized := normalizer.NormalizeProducts(rawProducts)

// 	// // Step 4: Transform to JSON-like structure based on categories
// 	// result := make(map[string][]map[string]any)
// 	// for _, p := range normalized {
// 	// 	if _, ok := result[p.Category]; !ok {
// 	// 		result[p.Category] = []map[string]any{}
// 	// 	}
// 	// 	result[p.Category] = append(result[p.Category], map[string]any{
// 	// 		"title": p.Title,
// 	// 		"price": p.Price,
// 	// 		"link":  p.Link,
// 	// 		"image": p.Image,
// 	// 	})
// 	// }s

// 	return recommendedProducts, nil
// }
