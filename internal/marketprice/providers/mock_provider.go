package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"samplegadai-backend/internal/dto"
)

type MockMarketProvider struct{}

func NewMockMarketProvider() *MockMarketProvider {
	return &MockMarketProvider{}
}

func (m *MockMarketProvider) Name() string {
	return "Verified Marketplace Data Aggregator"
}

func (m *MockMarketProvider) HealthCheck(ctx context.Context) bool {
	return true
}

func (m *MockMarketProvider) Search(ctx context.Context, query dto.MarketPriceQuery) ([]dto.MarketComparableResponse, error) {
	now := time.Now()
	brand := strings.ToLower(query.Brand)
	model := strings.ToLower(query.Model)

	var items []dto.MarketComparableResponse

	if strings.Contains(brand, "apple") || strings.Contains(model, "iphone") {
		items = []dto.MarketComparableResponse{
			{
				Source:        "Tokopedia Verified",
				ExternalID:    "TOK-101",
				Title:         "iPhone 15 Pro 256GB Natural Titanium Mulus Fullset",
				URL:           "https://tokopedia.com/item/101",
				Price:         14500000,
				Location:      "Jakarta Selatan",
				Condition:     "Bekas - Mulus",
				SellerType:    "Power Merchant",
				Brand:         "Apple",
				Model:         "iPhone 15 Pro",
				Variant:       "256GB",
				Specification: "Battery Health 98%, Garansi Resmi iBox",
				FetchedAt:     now.Add(-2 * 24 * time.Hour),
			},
			{
				Source:        "OLX Official",
				ExternalID:    "OLX-202",
				Title:         "iPhone 15 Pro 256GB Blue Titanium Garansi Aktif",
				URL:           "https://olx.co.id/item/202",
				Price:         14800000,
				Location:      "Jakarta Barat",
				Condition:     "Bekas - Mulus",
				SellerType:    "Individual",
				Brand:         "Apple",
				Model:         "iPhone 15 Pro",
				Variant:       "256GB",
				Specification: "Battery Health 96%, Fullset Original",
				FetchedAt:     now.Add(-1 * 24 * time.Hour),
			},
			{
				Source:        "Shopee Mall",
				ExternalID:    "SHP-303",
				Title:         "Apple iPhone 15 Pro 256 GB Second Like New",
				URL:           "https://shopee.co.id/item/303",
				Price:         14700000,
				Location:      "Bandung",
				Condition:     "Bekas - Like New",
				SellerType:    "Star Seller",
				Brand:         "Apple",
				Model:         "iPhone 15 Pro",
				Variant:       "256GB",
				Specification: "Battery Health 99%, No Minus",
				FetchedAt:     now.Add(-12 * time.Hour),
			},
			{
				Source:        "Tokopedia Verified",
				ExternalID:    "TOK-104",
				Title:         "iPhone 15 Pro 256GB Black Titanium Ex iBox",
				URL:           "https://tokopedia.com/item/104",
				Price:         14200000,
				Location:      "Surabaya",
				Condition:     "Bekas - Good",
				SellerType:    "Power Merchant",
				Brand:         "Apple",
				Model:         "iPhone 15 Pro",
				Variant:       "256GB",
				Specification: "Battery Health 90%, Casing agak lecet halus",
				FetchedAt:     now.Add(-5 * 24 * time.Hour),
			},
			{
				Source:        "Bukalapak Verified",
				ExternalID:    "BLK-505",
				Title:         "iPhone 15 Pro 256GB Titanium Gray Mulus",
				URL:           "https://bukalapak.com/item/505",
				Price:         15200000,
				Location:      "Jakarta Pusat",
				Condition:     "Bekas - Mulus",
				SellerType:    "Recommended Seller",
				Brand:         "Apple",
				Model:         "iPhone 15 Pro",
				Variant:       "256GB",
				Specification: "Battery Health 100%, Garansi Aktif",
				FetchedAt:     now.Add(-3 * 24 * time.Hour),
			},
		}
	} else if strings.Contains(brand, "samsung") || strings.Contains(model, "s24") {
		items = []dto.MarketComparableResponse{
			{
				Source:        "Tokopedia Verified",
				ExternalID:    "TOK-SAM1",
				Title:         "Samsung Galaxy S24 Ultra 5G 12/256GB Garansi Sein",
				URL:           "https://tokopedia.com/item/sam1",
				Price:         16200000,
				Location:      "Jakarta Barat",
				Condition:     "Bekas - Mulus",
				SellerType:    "Power Merchant",
				Brand:         "Samsung",
				Model:         "Galaxy S24 Ultra",
				Variant:       "256GB",
				Specification: "Garansi SEIN Aktif, Fullset Box",
				FetchedAt:     now.Add(-1 * 24 * time.Hour),
			},
			{
				Source:        "OLX Official",
				ExternalID:    "OLX-SAM2",
				Title:         "Samsung S24 Ultra 256GB Titanium Grey Mulus",
				URL:           "https://olx.co.id/item/sam2",
				Price:         16500000,
				Location:      "Tangerang",
				Condition:     "Bekas - Like New",
				SellerType:    "Individual",
				Brand:         "Samsung",
				Model:         "Galaxy S24 Ultra",
				Variant:       "256GB",
				Specification: "Fullset Ori, No Minus",
				FetchedAt:     now.Add(-2 * 24 * time.Hour),
			},
		}
	} else if strings.Contains(brand, "emas") || strings.Contains(model, "antam") {
		items = []dto.MarketComparableResponse{
			{
				Source:        "Harga Emas Pegadaian Spot",
				ExternalID:    "GOLD-01",
				Title:         "Emas Antam Logam Mulia 10 Gram Certicard",
				URL:           "https://logammulia.com/item/10g",
				Price:         13200000,
				Location:      "Jakarta",
				Condition:     "Baru / Mint",
				SellerType:    "Official Distributor",
				Brand:         "Antam",
				Model:         "Logam Mulia",
				Variant:       "10 Gram",
				Specification: "Kadar 999.9 Certicard Redefine",
				FetchedAt:     now,
			},
			{
				Source:        "Tokopedia Gold Market",
				ExternalID:    "GOLD-02",
				Title:         "LM Antam 10 Gram Press Terbaru Redefine",
				URL:           "https://tokopedia.com/gold/10g",
				Price:         13150000,
				Location:      "Jakarta Selatan",
				Condition:     "Baru / Mint",
				SellerType:    "Official Store",
				Brand:         "Antam",
				Model:         "Logam Mulia",
				Variant:       "10 Gram",
				Specification: "Kadar 99.99%",
				FetchedAt:     now.Add(-3 * time.Hour),
			},
		}
	} else {
		// Generic fallback comparable listings
		basePrice := int64(3500000)
		for i := 1; i <= 4; i++ {
			items = append(items, dto.MarketComparableResponse{
				Source:        "Marketplace Aggregator",
				ExternalID:    fmt.Sprintf("GEN-%d", i),
				Title:         fmt.Sprintf("%s %s Second Mulus", query.Brand, query.Model),
				URL:           "https://marketplace.co.id",
				Price:         basePrice + int64(i*200000),
				Location:      "Jakarta",
				Condition:     "Bekas - Baik",
				SellerType:    "Verified Seller",
				Brand:         query.Brand,
				Model:         query.Model,
				Variant:       query.Variant,
				Specification: "Fullset Box Original",
				FetchedAt:     now.Add(-time.Duration(i) * 24 * time.Hour),
			})
		}
	}

	return items, nil
}
