package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
)

// ListCatalogUseCase orchestrates listing products with pagination.
type ListCatalogUseCase struct {
	productRepo product.Repository
}

// NewListCatalogUseCase creates a new instance of ListCatalogUseCase.
func NewListCatalogUseCase(productRepo product.Repository) *ListCatalogUseCase {
	return &ListCatalogUseCase{
		productRepo: productRepo,
	}
}

// Response represents the catalog listing response with products and pagination info.
type Response struct {
	Products []product.Product `json:"products"`
	Total    int64             `json:"total"`
}

// Execute retrieves products with pagination and optional filters.
func (uc *ListCatalogUseCase) Execute(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) (*Response, error) {
	products, total, err := uc.productRepo.FindAll(ctx, offset, limit, filters...)
	if err != nil {
		return nil, err
	}

	return &Response{
		Products: products,
		Total:    total,
	}, nil
}

// FilterBuilder provides a fluent API for building filters.
type FilterBuilder struct {
	filter product.FindAllFilter
	err    error
}

// NewFilter creates a new FilterBuilder.
func NewFilter() *FilterBuilder {
	return &FilterBuilder{
		filter: product.FindAllFilter{},
	}
}

// SetCategory sets the category filter.
func (fb *FilterBuilder) SetCategory(code string) *FilterBuilder {
	if fb.err == nil {
		fb.filter.CategoryCode = code
	}
	return fb
}

// SetMaxPrice sets the maximum price filter.
func (fb *FilterBuilder) SetMaxPrice(price decimal.Decimal) *FilterBuilder {
	if fb.err == nil {
		fb.filter.MaxPrice = price
	}
	return fb
}

// ToFilterOptions returns the filter as a variadic argument option.
func (fb *FilterBuilder) ToFilterOptions() []product.FindAllFilter {
	if fb.err != nil {
		return []product.FindAllFilter{}
	}
	return []product.FindAllFilter{fb.filter}
}

// ParseDecimal parses a string to decimal, returning zero value on error.
func ParseDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}
