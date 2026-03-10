package persistence

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) FindAll(ctx context.Context) ([]category.Category, error) {
	var models []Category

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]category.Category, len(models))
	for i, m := range models {
		categories[i] = *m.ToDomainCategory()
	}

	return categories, nil
}

func (r *CategoriesRepository) Create(ctx context.Context, code, name string) (*category.Category, error) {
	model := Category{
		Code: code,
		Name: name,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return model.ToDomainCategory(), nil
}
