package marketprice

import (
	"math"
	"sort"

	"samplegadai-backend/internal/dto"
)

type PriceCalculator struct{}

func NewPriceCalculator() *PriceCalculator {
	return &PriceCalculator{}
}

func (c *PriceCalculator) ProcessAndCalculate(comparables []dto.MarketComparableResponse, historicalPrice int64) dto.MarketPriceResult {
	if len(comparables) == 0 {
		recPrice := historicalPrice
		if recPrice == 0 {
			recPrice = 5000000
		}
		return dto.MarketPriceResult{
			RecommendedPrice: recPrice,
			MinimumPrice:     recPrice,
			MaximumPrice:     recPrice,
			MedianPrice:      recPrice,
			TrimmedMeanPrice: recPrice,
			ConfidenceScore:  50,
			ComparableCount:  0,
			HistoricalPrice:  historicalPrice,
			Comparables:      []dto.MarketComparableResponse{},
		}
	}

	// 1. Filter out low similarity items (< 40%)
	var filtered []dto.MarketComparableResponse
	for _, item := range comparables {
		if item.SimilarityScore >= 40.0 {
			filtered = append(filtered, item)
		}
	}

	if len(filtered) == 0 {
		filtered = comparables
	}

	// 2. Sort by price
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Price < filtered[j].Price
	})

	// 3. Remove outlier prices via Trimmed IQR method
	prices := make([]int64, len(filtered))
	for i, item := range filtered {
		prices[i] = item.Price
	}

	cleanPrices := removeOutliers(prices)

	// 4. Compute statistics
	minPrice := cleanPrices[0]
	maxPrice := cleanPrices[len(cleanPrices)-1]
	medianPrice := calculateMedian(cleanPrices)
	trimmedMeanPrice := calculateTrimmedMean(cleanPrices)

	// Recommended price is median weighted with trimmed mean
	recommendedPrice := int64(float64(medianPrice)*0.7 + float64(trimmedMeanPrice)*0.3)

	// 5. Calculate Confidence Score (0-100)
	confidence := calculateConfidence(cleanPrices, filtered, historicalPrice)

	return dto.MarketPriceResult{
		RecommendedPrice: recommendedPrice,
		MinimumPrice:     minPrice,
		MaximumPrice:     maxPrice,
		MedianPrice:      medianPrice,
		TrimmedMeanPrice: trimmedMeanPrice,
		ConfidenceScore:  confidence,
		ComparableCount:  len(filtered),
		HistoricalPrice:  historicalPrice,
		Comparables:      filtered,
	}
}

func removeOutliers(prices []int64) []int64 {
	n := len(prices)
	if n < 4 {
		return prices
	}

	q1 := float64(prices[n/4])
	q3 := float64(prices[(n*3)/4])
	iqr := q3 - q1

	lowerBound := int64(q1 - 1.5*iqr)
	upperBound := int64(q3 + 1.5*iqr)

	var result []int64
	for _, p := range prices {
		if p >= lowerBound && p <= upperBound {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		return prices
	}
	return result
}

func calculateMedian(prices []int64) int64 {
	n := len(prices)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return prices[n/2]
	}
	return (prices[n/2-1] + prices[n/2]) / 2
}

func calculateTrimmedMean(prices []int64) int64 {
	n := len(prices)
	if n <= 2 {
		return calculateMedian(prices)
	}
	// Trim top and bottom 10%
	trimCount := int(math.Floor(float64(n) * 0.1))
	if trimCount < 1 {
		trimCount = 1
	}

	var sum int64
	count := 0
	for i := trimCount; i < n-trimCount; i++ {
		sum += prices[i]
		count++
	}

	if count == 0 {
		return calculateMedian(prices)
	}

	return sum / int64(count)
}

func calculateConfidence(cleanPrices []int64, filtered []dto.MarketComparableResponse, historicalPrice int64) int {
	score := 50.0

	// Quantity factor (+up to 25)
	countFactor := math.Min(float64(len(filtered))*2.5, 25.0)
	score += countFactor

	// Similarity score average (+up to 15)
	var simSum float64
	for _, item := range filtered {
		simSum += item.SimilarityScore
	}
	if len(filtered) > 0 {
		score += (simSum / float64(len(filtered))) * 0.15
	}

	// Variance factor (- penalty for high variance)
	minP := float64(cleanPrices[0])
	maxP := float64(cleanPrices[len(cleanPrices)-1])
	if minP > 0 {
		spread := (maxP - minP) / minP
		if spread < 0.15 {
			score += 10.0
		} else if spread > 0.40 {
			score -= 10.0
		}
	}

	// Historical alignment factor
	if historicalPrice > 0 {
		score += 5.0
	}

	finalScore := int(math.Round(score))
	if finalScore > 100 {
		return 100
	}
	if finalScore < 0 {
		return 0
	}
	return finalScore
}
