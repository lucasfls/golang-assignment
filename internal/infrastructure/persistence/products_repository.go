package persistence

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"gorm.io/gorm"
)

// ProductsRepository implements the product.Repository interface.
type ProductsRepository struct {
	db *gorm.DB
}

// NewProductsRepository creates a new instance of ProductsRepository.
func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// FindAll retrieves all products with their variants.
func (r *ProductsRepository) FindAll(ctx context.Context) ([]product.Product, error) {
	var products []Product

	err := r.db.WithContext(ctx).
		Preload("Variants").
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	// Convert persistence models to domain models
	domainProducts := make([]product.Product, len(products))
	for i, p := range products {
		domainProducts[i] = *p.ToDomainProduct()
	}

	return domainProducts, nil
}
