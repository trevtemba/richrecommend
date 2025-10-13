package normalizer

import (
	"context"
	"time"

	"github.com/trevtemba/richrecommend/internal/ai"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
	"golang.org/x/sync/errgroup"
)

func NormalizeProducts(loadedData models.ScraperResponse, includedFields []string, requestId string) (models.NormalizerResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Normalizer agent started", "request_id", requestId)
	var normalizerResponse models.NormalizerResponse

	parserAgentClient, err := ai.NewClient("localhost:50051")
	if err != nil {
		return normalizerResponse, err
	}

	batchSize := 3
	numProducts := len(loadedData.ProductsScraped)
	numBatches := (numProducts + batchSize - 1) / batchSize
	parserBatches := make([][]map[string]map[string]any, 0, numBatches)

	for i := 0; i < numProducts; i += batchSize {
		end := i + batchSize
		if end > numProducts {
			end = numProducts
		}
		batch := loadedData.ProductsScraped[i:end]
		parserBatches = append(parserBatches, batch)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	ch := make(chan []models.ProductData)
	failedCh := make(chan int)

	for batchNum, batch := range parserBatches {
		bn := batchNum
		b := batch
		eg.Go(func() error {
			reqCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
			defer cancel()
			logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, "Sending scraped data to parsing agent for pre-normalization...", "request_id", requestId)

			resp, err := parserAgentClient.GetProductData(reqCtx, b)
			if err != nil {
				failedCh <- bn
				logger.Log(logger.LogTypeAgentDebug, logger.LevelDebug, "batch #%d failed...", "request_id", bn, requestId)
				return nil
			}
			ch <- resp
			return nil
		})
	}
	go func() {
		if err := eg.Wait(); err != nil {
			logger.Log(logger.LogTypeAgentDebug, logger.LevelDebug, "err group wait failed...", "request_id", requestId, "error", err)
		}
		close(ch)
		close(failedCh)
	}()

	var parserResponses []map[string]models.ProductData
	for productDataBatch := range ch {
		for _, productData := range productDataBatch {
			parserResponses = append(parserResponses, map[string]models.ProductData{
				productData.Name: productData,
			})
		}
	}

	// Temporary! Just going to convert this to map[string]any so I can debug parser agent before filtering fields.
	var normalizedRecs []map[string]map[string]any
	for _, response := range parserResponses {
		var productName string
		var productDataMap models.ProductData
		for key, value := range response {
			productName = key
			productDataMap = value
			break
		}
		normalizedRecs = append(normalizedRecs, map[string]map[string]any{
			productName: {
				"description": productDataMap.Description,
				"thumbnail":   productDataMap.Thumbnail,
				"ingredients": productDataMap.Ingredients,
				"retailers":   productDataMap.Retailers,
			},
		})
	}
	normalizerResponse.Recommendations = normalizedRecs
	normalizerResponse.FailedProducts = loadedData.FailedProducts
	return normalizerResponse, nil
}
