package normalizer

import (
	"context"
	"fmt"
	"time"

	"github.com/trevtemba/richrecommend/internal/ai"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
	"golang.org/x/sync/errgroup"
)

func NormalizeProducts(loadedData models.ScraperResponse, includedFields []string, requestId string) (models.NormalizerResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Normalizer agent started", "request_id", requestId)
	var normalizerResponse models.NormalizerResponse

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	agent, err := ai.NewClient("localhost:50051")
	if err != nil {
		return normalizerResponse, err
	}

	ch := make(chan map[string]any, len(loadedData.ProductsScraped))
	failedCh := make(chan string, len(loadedData.ProductsScraped))

	for _, productMap := range loadedData.ProductsScraped {
		for productName, scrapedData := range productMap {
			pn := productName
			eg.Go(func() error {
				reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				resCh := make(chan map[string]any, 1)
				failCh := make(chan string, 1)
				go func() {
					productData := scrapedData["product_result"]
					productDataTyped := productData.(map[string]any)
					relevantProductData, err := agent.GetProductData(productDataTyped)
					if err != nil {
						failCh <- pn
					} else {
						resCh <- map[string]any{productName: relevantProductData}
					}
				}()

				select {
				case <-reqCtx.Done():
					failCh <- pn
					logger.Log(logger.LogTypeAgentDebug, logger.LevelDebug, fmt.Sprintf("Scraper for %s timed out", pn), "request_id", requestId)
					return nil
				case fail := <-failCh:
					//todo add a retry so can attempt to correct!
					logger.Log(logger.LogTypeAgentDebug, logger.LevelDebug, fmt.Sprintf("Scraper for %s failed", pn), "request_id", requestId)
					failCh <- fail
					return nil
				case res := <-resCh:
					ch <- res
					return nil
				}
			})

		}
	}

	if err := eg.Wait(); err != nil {
		return normalizerResponse, fmt.Errorf("%w", err)
	}
	close(ch)
	close(failedCh)

	for productData := range ch {
		normalizerResponse.Recommendations = append(normalizerResponse.Recommendations, productData)
	}
	for failedProducts := range failedCh {
		normalizerResponse.FailedProducts = append(normalizerResponse.FailedProducts, failedProducts)
	}
	return normalizerResponse, nil
}
