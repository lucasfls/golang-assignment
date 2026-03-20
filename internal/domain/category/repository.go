package category

import "context"

// Repository defines the interface for persisting and retrieving categories.
// This is a domain interface - implementations live in the infrastructure layer.
type Repository interface {
	// FindAll returns all categories.
	FindAll(ctx context.Context) ([]Category, error)

	// Create creates a new category.
	Create(ctx context.Context, code, name string) (*Category, error)
}
