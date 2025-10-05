package normalizer

import (
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
)

func NormalizeProducts(loadedData models.ScraperResponse, includedFields []string, requestId string) (models.NormalizerResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Normalizer agent started", "request_id", requestId)
	var normalizedRecommendations models.NormalizerResponse

	normalizedRecommendationsMap := make(map[string]any)

	for _, productMap := range loadedData.ProductsScraped {
		for productName, scrapedData := range productMap {
			productData := scrapedData["product_result"]

			normalizedRecommendationsMap[productName] = productData
			// for _, fieldName := range includedFields {
			// 	normalizedRecommendationsMap[productName] = productData[fieldName]
			// }
		}
	}
	normalizedRecommendations.Recommendations = normalizedRecommendationsMap
	return normalizedRecommendations, nil
}
