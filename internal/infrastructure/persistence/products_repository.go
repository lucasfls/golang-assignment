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

// FindAll retrieves all products with pagination support.
// Returns slice of products, total count, and error.
func (r *ProductsRepository) FindAll(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
	var persistenceProducts []Product
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Preload("Category").
		Preload("Variants").
		Find(&persistenceProducts).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain models
	domainProducts := make([]product.Product, len(persistenceProducts))
	for i, p := range persistenceProducts {
		domainProducts[i] = *p.ToDomainProduct()
	}

	return domainProducts, total, nil
}
