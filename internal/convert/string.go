package convert

import (
	"api-chi/cmd/models"
	"strings"
)

func StringToBlogtagSlice(input string) []models.BlogTag {
	result := []models.BlogTag{}

	parts := strings.Split(input, ";")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if part != "" {
			result = append(result, models.BlogTag{Name: part})
		}
	}

	return result
}
