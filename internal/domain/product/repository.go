package product

import (
	"context"

	"github.com/shopspring/decimal"
)

// FindAllFilter contains optional filter criteria for product queries.
type FindAllFilter struct {
	CategoryCode string          // Filter by category code (optional)
	MaxPrice     decimal.Decimal // Filter for products with price less than this (optional)
}

// Repository defines the interface for persisting and retrieving products.
// This is a domain interface - implementations live in the infrastructure layer.
type Repository interface {
	// FindAll returns all products with pagination and optional filters.
	// Returns slice of products, total count, and error.
	FindAll(ctx context.Context, offset, limit int, filters ...FindAllFilter) ([]Product, int64, error)

	// FindByCode returns a product by its code.
	// Returns the product or an error if not found.
	FindByCode(ctx context.Context, code string) (*Product, error)
}
