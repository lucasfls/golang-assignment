package product

import "context"

// Repository defines the interface for persisting and retrieving products.
// This is a domain interface - implementations live in the infrastructure layer.
type Repository interface {
	// FindAll returns all products with pagination support.
	// Returns slice of products, total count, and error.
	FindAll(ctx context.Context, offset, limit int) ([]Product, int64, error)
}
