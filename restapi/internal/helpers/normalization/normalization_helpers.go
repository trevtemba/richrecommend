package normalization

import (
	"slices"

	"github.com/trevtemba/richrecommend/internal/models"
)

func IsValidField(field string) bool {
	return slices.Contains([]string{"ingredients", "description", "thumbnail", "retailers"}, field)
}

func GetField(productData models.ProductData, field string) any {
	var fieldData any
	switch field {
	case "name":
		fieldData = productData.Name
	case "description":
		fieldData = productData.Description
	case "thumbnail":
		fieldData = productData.Thumbnail
	case "ingredients":
		fieldData = productData.Ingredients
	case "retailers":
		fieldData = productData.Retailers
	}
	return fieldData
}
