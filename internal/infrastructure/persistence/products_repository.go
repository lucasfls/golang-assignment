package persistence

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{db: db}
}

func (r *ProductsRepository) FindAll(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
	var models []Product
	var total int64

	query := r.db.WithContext(ctx).Joins("JOIN categories ON products.category_id = categories.id")

	if filter.CategoryCode != "" {
		query = query.Where("categories.code = ?", filter.CategoryCode)
	}

	if !filter.MaxPrice.IsZero() {
		query = query.Where("products.price < ?", filter.MaxPrice)
	}

	if err := query.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Preload("Category").
		Preload("Variants").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	products := make([]product.Product, len(models))
	for i, m := range models {
		products[i] = *m.ToDomainProduct()
	}

	return products, total, nil
}

func (r *ProductsRepository) FindByCode(ctx context.Context, code string) (*product.Product, error) {
	var model Product

	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		Preload("Category").
		Preload("Variants").
		First(&model).Error; err != nil {
		return nil, err
	}

	return model.ToDomainProduct(), nil
}
