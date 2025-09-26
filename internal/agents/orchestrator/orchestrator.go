package orchestrator

import (
	"github.com/trevtemba/richrecommend/internal/agents/recommendation"
	"github.com/trevtemba/richrecommend/internal/models"
)

func RunPipelineWithParams(params models.OrchestratorParams) (any, error) {
	// Step 1: Recommendation Agent
	recommendedProducts, err := recommendation.GenerateWithParams(params)
	if err != nil {
		return nil, err
	}

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
	// }

	return result, nil
}
