package scraper

import (
	"context"
	"fmt"
	"maps"
	"os"
	"time"

	g "github.com/serpapi/google-search-results-golang"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
	"golang.org/x/sync/errgroup"
)

func ScrapeProducts(recommendedProducts models.RecommendationResponse, requestId string) (models.ScraperResponse, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Scraper agent started", "request_id", requestId)
	var scraperResponse models.ScraperResponse

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	ch := make(chan map[string]map[string]any, recommendedProducts.ItemCount)
	failedCh := make(chan string, recommendedProducts.ItemCount)

	params := map[string]string{
		"q":             "",
		"location":      "Austin, Texas, United States",
		"hl":            "en",
		"gl":            "us",
		"google_domain": "google.com",
	}

	for _, productList := range recommendedProducts.Recommendation {
		for _, productName := range productList {
			pn := productName
			eg.Go(func() error {
				logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, fmt.Sprintf("Scraping data for %s...", pn), "request_id", requestId)

				localParams := maps.Clone(params)
				localParams["q"] = pn

				reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				failCh := make(chan string, 1)
				resCh := make(chan map[string]map[string]any, 1)

				go func() {

					search := g.NewGoogleSearch(localParams, os.Getenv("SERP_API_KEY"))
					results, err := search.GetJSON()

					if err != nil {
						failCh <- pn
					} else {
						resCh <- map[string]map[string]any{pn: results}
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
		return scraperResponse, fmt.Errorf("%w", err)
	}
	close(ch)
	close(failedCh)

	for productData := range ch {
		scraperResponse.ProductsScraped = append(scraperResponse.ProductsScraped, productData)
	}
	for failedProducts := range failedCh {
		scraperResponse.FailedProducts = append(scraperResponse.FailedProducts, failedProducts)
	}

	return scraperResponse, nil
}
