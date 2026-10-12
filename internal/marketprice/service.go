package marketprice

import (
	"context"
	"log"

	"samplegadai-backend/internal/dto"
	"samplegadai-backend/internal/marketprice/providers"
	"samplegadai-backend/internal/model"

	"gorm.io/gorm"
)

type MarketPriceService struct {
	db         *gorm.DB
	providers  []MarketPriceProvider
	normalizer *Normalizer
	calculator *PriceCalculator
}

func NewMarketPriceService(db *gorm.DB) *MarketPriceService {
	mockProv := providers.NewMockMarketProvider()
	return &MarketPriceService{
		db:         db,
		providers:  []MarketPriceProvider{mockProv},
		normalizer: NewNormalizer(),
		calculator: NewPriceCalculator(),
	}
}

func (s *MarketPriceService) FetchMarketPrice(ctx context.Context, query dto.MarketPriceQuery) (dto.MarketPriceResult, error) {
	var rawComparables []dto.MarketComparableResponse

	// Fetch from all registered providers
	for _, provider := range s.providers {
		if provider.HealthCheck(ctx) {
			res, err := provider.Search(ctx, query)
			if err != nil {
				log.Printf("Provider %s error: %v", provider.Name(), err)
				continue
			}
			rawComparables = append(rawComparables, res...)
		}
	}

	// Compute similarity scores for all comparables
	for i := range rawComparables {
		rawComparables[i].SimilarityScore = s.normalizer.CalculateSimilarity(query, rawComparables[i])
	}

	// Fetch internal historical price from PostgreSQL database appraisals
	var historicalPrice int64
	var pastAppraisals []model.Appraisal
	err := s.db.Joins("JOIN assets ON assets.id = appraisals.asset_id").
		Where("LOWER(assets.brand) = LOWER(?) AND LOWER(assets.model) = LOWER(?)", query.Brand, query.Model).
		Order("appraisals.appraised_at DESC").
		Limit(5).
		Find(&pastAppraisals).Error

	if err == nil && len(pastAppraisals) > 0 {
		var sum int64
		for _, app := range pastAppraisals {
			sum += app.AppraiserPrice
		}
		historicalPrice = sum / int64(len(pastAppraisals))
	}

	// Process comparables: filter duplicates & low similarity, remove outliers, compute statistics
	result := s.calculator.ProcessAndCalculate(rawComparables, historicalPrice)
	return result, nil
}
