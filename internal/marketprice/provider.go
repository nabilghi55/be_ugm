package marketprice

import (
	"context"
	"samplegadai-backend/internal/dto"
)

type MarketPriceProvider interface {
	Name() string
	Search(ctx context.Context, query dto.MarketPriceQuery) ([]dto.MarketComparableResponse, error)
	HealthCheck(ctx context.Context) bool
}
