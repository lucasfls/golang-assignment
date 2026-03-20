package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
)

// MockRepository is a mock implementation of product.Repository for testing.
type MockRepository struct {
	FindAllFunc    func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error)
	FindByCodeFunc func(ctx context.Context, code string) (*product.Product, error)
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		FindAllFunc: func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
			return []product.Product{}, 0, nil
		},
		FindByCodeFunc: func(ctx context.Context, code string) (*product.Product, error) {
			return nil, nil
		},
	}
}

func (m *MockRepository) FindAll(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
	return m.FindAllFunc(ctx, offset, limit, filter)
}

func (m *MockRepository) FindByCode(ctx context.Context, code string) (*product.Product, error) {
	return m.FindByCodeFunc(ctx, code)
}

func (m *MockRepository) WithFindAll(fn func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error)) *MockRepository {
	m.FindAllFunc = fn
	return m
}

// WithFindByCode sets the behavior for FindByCode method.
func (m *MockRepository) WithFindByCode(fn func(ctx context.Context, code string) (*product.Product, error)) *MockRepository {
	m.FindByCodeFunc = fn
	return m
}

// GetSampleProducts returns sample products for testing.
func GetSampleProducts() []product.Product {
	clothingCategory := &category.Category{
		ID:   1,
		Code: "CLOTHING",
		Name: "Clothing",
	}
	shoesCategory := &category.Category{
		ID:   2,
		Code: "SHOES",
		Name: "Shoes",
	}

	return []product.Product{
		{
			ID:       1,
			Code:     "PROD001",
			Price:    mustDecimal("29.99"),
			Category: clothingCategory,
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
			ID:       2,
			Code:     "PROD002",
			Price:    mustDecimal("49.99"),
			Category: shoesCategory,
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
			Category: clothingCategory,
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

// MustDecimal is an exported version of mustDecimal for use in tests.
func MustDecimal(s string) decimal.Decimal {
	return mustDecimal(s)
}
