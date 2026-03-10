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

// FindAll retrieves products with pagination and optional filters support.
// Returns slice of products, total count, and error.
func (r *ProductsRepository) FindAll(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
	var persistenceProducts []Product
	var total int64

	// Start building query with joined categories table
	query := r.db.WithContext(ctx).Joins("JOIN categories ON products.category_id = categories.id")

	// Apply filters if provided
	if len(filters) > 0 && filters[0].CategoryCode != "" {
		query = query.Where("categories.code = ?", filters[0].CategoryCode)
	}

	if len(filters) > 0 && !filters[0].MaxPrice.IsZero() {
		query = query.Where("products.price < ?", filters[0].MaxPrice)
	}

	// Get total count with filters applied
	if err := query.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with filters
	if err := query.
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

// FindByCode retrieves a product by its code.
// Returns the product or an error if not found.
func (r *ProductsRepository) FindByCode(ctx context.Context, code string) (*product.Product, error) {
	var persistenceProduct Product

	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		Preload("Category").
		Preload("Variants").
		First(&persistenceProduct).Error; err != nil {
		return nil, err
	}

	domainProduct := persistenceProduct.ToDomainProduct()
	return domainProduct, nil
}
