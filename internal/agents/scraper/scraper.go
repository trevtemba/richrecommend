package scraper

import (
	"fmt"
	"os"

	g "github.com/serpapi/google-search-results-golang"
	"github.com/trevtemba/richrecommend/internal/logger"
	"github.com/trevtemba/richrecommend/internal/models"
	"golang.org/x/sync/errgroup"
)

func ScrapeProducts(recommendedProducts models.RecommendationResponse, requestId string) (map[string]any, error) {
	logger.Log(logger.LogTypeAgentFinish, logger.LevelInfo, "Scraper agent started", "request_id", requestId)

	params := map[string]string{
		"q":             "",
		"location":      "Austin, Texas, United States",
		"hl":            "en",
		"gl":            "us",
		"google_domain": "google.com",
	}

	eg := new(errgroup.Group)
	c := make(chan map[string]any, recommendedProducts.ItemCount)

	for _, productList := range recommendedProducts.Recommendation {
		for _, productName := range productList {
			eg.Go(func() error {
				logger.Log(logger.LogTypeAgentWork, logger.LevelDebug, fmt.Sprintf("Scraping data for %s...", productName), "request_id", requestId)
				params["q"] = productName

				search := g.NewGoogleSearch(params, os.Getenv("SERP_API_KEY"))
				results, err := search.GetJSON()
				if err != nil {
					return fmt.Errorf("serpAPI search for %s failed: %w", productName, err)
				}

				c <- results

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	loadedData := <-c
	return loadedData, nil
}
