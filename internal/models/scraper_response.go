package models

type ScraperResponse struct {
	ProductsScraped []map[string]any `json:"products_scraped"`
	FailedProducts  []string         `json:"failed_products"`
}
