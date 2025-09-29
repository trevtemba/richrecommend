package scraper

import "github.com/trevtemba/richrecommend/internal/logger"

func ScrapeProducts(recommendedProducts map[string][]string, requestId string) (map[string]any, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Scraper agent started", "request_id", requestId)
	var loadedData map[string]any

	return loadedData, nil
}
