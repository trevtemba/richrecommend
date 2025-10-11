package normalizer

import (
	"github.com/trevtemba/richrecommend/internal/ai"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
)

func NormalizeProducts(loadedData models.ScraperResponse, includedFields []string, requestId string) (models.NormalizerResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Normalizer agent started", "request_id", requestId)
	var normalizerResponse models.NormalizerResponse

	parserAgentClient, err := ai.NewClient("localhost:50051")
	if err != nil {
		return normalizerResponse, err
	}
	resp, err := parserAgentClient.GetProductData(map[string]any{"test": "hello", "test2": "helloagain"})
	if err != nil {
		return normalizerResponse, err
	}

	parserResponses := make([]map[string]any, len(loadedData.ProductsScraped))

	for _, productData := range resp {
		parserResponses = append(parserResponses, map[string]any{
			productData.Name: map[string]any{
				"data": productData,
			},
		})
	}

	normalizerResponse.Recommendations = parserResponses
	return normalizerResponse, nil
}
