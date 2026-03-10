package product

import "context"

// Repository defines the interface for persisting and retrieving products.
// This is a domain interface - implementations live in the infrastructure layer.
type Repository interface {
	// FindAll returns all products in the catalog.
	FindAll(ctx context.Context) ([]Product, error)
}
