package models

type Retailer struct {
	Name    string  `json:"name"`
	Link    string  `json:"link"`
	Rating  float64 `json:"rating"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

type ProductData struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Thumbnail   string     `json:"thumbnail"`
	Ingredients []string   `json:"ingredients"`
	Retailers   []Retailer `json:"retailers"`
}
