package normalizer

import (
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
)

func NormalizeProducts(loadedData models.ScraperResponse, includedFields []string, requestId string) (models.NormalizerResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Normalizer agent started", "request_id", requestId)
	var normalizedRecommendations models.NormalizerResponse
	return normalizedRecommendations, nil
}
