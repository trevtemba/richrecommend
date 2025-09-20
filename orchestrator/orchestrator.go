package orchestrator

import (
	"github.com/trevtemba/richrecommend/models"
	"github.com/trevtemba/richrecommend/normalizer"
	"github.com/trevtemba/richrecommend/recommendation"
	"github.com/trevtemba/richrecommend/scraper"
)

func RunPipelineWithParams(params models.OrchestratorParams) (interface{}, error) {
	// Step 1: Recommendation Agent
	recommendedProducts, err := recommendation.GenerateWithParams(params)
	if err != nil {
		return nil, err
	}

	// Step 2: Scraper Agents
	rawProducts := scraper.ScrapeProducts(recommendedProducts)

	// Step 3: Normalizer Agent
	normalized := normalizer.NormalizeProducts(rawProducts)

	// Step 4: Transform to JSON-like structure based on categories
	result := make(map[string][]map[string]interface{})
	for _, p := range normalized {
		if _, ok := result[p.Category]; !ok {
			result[p.Category] = []map[string]interface{}{}
		}
		result[p.Category] = append(result[p.Category], map[string]interface{}{
			"title": p.Title,
			"price": p.Price,
			"link":  p.Link,
			"image": p.Image,
		})
	}

	return result, nil
}
