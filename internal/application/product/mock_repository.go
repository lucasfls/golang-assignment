package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
)

// MockRepository is a mock implementation of product.Repository for testing.
type MockRepository struct {
	FindAllFunc func(ctx context.Context) ([]product.Product, error)
}

// NewMockRepository creates a new MockRepository with default behavior.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		FindAllFunc: func(ctx context.Context) ([]product.Product, error) {
			return []product.Product{}, nil
		},
	}
}

// FindAll delegates to the mocked function.
func (m *MockRepository) FindAll(ctx context.Context) ([]product.Product, error) {
	return m.FindAllFunc(ctx)
}

// WithFindAll sets the behavior for FindAll method.
func (m *MockRepository) WithFindAll(fn func(ctx context.Context) ([]product.Product, error)) *MockRepository {
	m.FindAllFunc = fn
	return m
}

// GetSampleProducts returns sample products for testing.
func GetSampleProducts() []product.Product {
	return []product.Product{
		{
			ID:    1,
			Code:  "PROD001",
			Price: mustDecimal("29.99"),
			Variants: []product.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Small",
					SKU:       "PROD001-S",
					Price:     mustDecimal("25.99"),
				},
				{
					ID:        2,
					ProductID: 1,
					Name:      "Large",
					SKU:       "PROD001-L",
					Price:     decimal.Zero,
				},
			},
		},
		{
			ID:    2,
			Code:  "PROD002",
			Price: mustDecimal("49.99"),
			Variants: []product.Variant{
				{
					ID:        3,
					ProductID: 2,
					Name:      "Red",
					SKU:       "PROD002-R",
					Price:     decimal.Zero,
				},
			},
		},
		{
			ID:       3,
			Code:     "PROD003",
			Price:    mustDecimal("19.99"),
			Variants: []product.Variant{},
		},
	}
}

// mustDecimal converts a string to decimal, panics on error.
func mustDecimal(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}
