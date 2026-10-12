package marketprice

import (
	"math"
	"strings"

	"samplegadai-backend/internal/dto"
)

type Normalizer struct{}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) CalculateSimilarity(query dto.MarketPriceQuery, item dto.MarketComparableResponse) float64 {
	score := 0.0

	// Brand similarity (20%)
	if strings.EqualFold(strings.TrimSpace(query.Brand), strings.TrimSpace(item.Brand)) {
		score += 20.0
	} else if strings.Contains(strings.ToLower(item.Brand), strings.ToLower(query.Brand)) {
		score += 10.0
	}

	// Model similarity (30%)
	if strings.EqualFold(strings.TrimSpace(query.Model), strings.TrimSpace(item.Model)) {
		score += 30.0
	} else if strings.Contains(strings.ToLower(item.Model), strings.ToLower(query.Model)) {
		score += 20.0
	}

	// Variant similarity (15%)
	if query.Variant != "" {
		if strings.EqualFold(strings.TrimSpace(query.Variant), strings.TrimSpace(item.Variant)) {
			score += 15.0
		} else if strings.Contains(strings.ToLower(item.Variant), strings.ToLower(query.Variant)) {
			score += 10.0
		}
	} else {
		score += 15.0
	}

	// Specification similarity (15%)
	if query.Specification != "" {
		if strings.Contains(strings.ToLower(item.Specification), strings.ToLower(query.Specification)) {
			score += 15.0
		} else {
			score += 5.0
		}
	} else {
		score += 15.0
	}

	// Condition similarity (10%)
	if query.Condition != "" {
		if strings.EqualFold(strings.TrimSpace(query.Condition), strings.TrimSpace(item.Condition)) {
			score += 10.0
		} else {
			score += 5.0
		}
	} else {
		score += 10.0
	}

	// Year similarity (10%)
	if query.Year > 0 && item.FetchedAt.Year() > 0 {
		diff := math.Abs(float64(query.Year - item.FetchedAt.Year()))
		if diff == 0 {
			score += 10.0
		} else if diff == 1 {
			score += 5.0
		}
	} else {
		score += 10.0
	}

	return score
}
