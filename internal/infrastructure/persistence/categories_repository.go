package persistence

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"gorm.io/gorm"
)

// CategoriesRepository implements the category.Repository interface.
type CategoriesRepository struct {
	db *gorm.DB
}

// NewCategoriesRepository creates a new instance of CategoriesRepository.
func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

// FindAll retrieves all categories.
// Returns slice of categories or an error.
func (r *CategoriesRepository) FindAll(ctx context.Context) ([]category.Category, error) {
	var persistenceCategories []Category

	if err := r.db.WithContext(ctx).Find(&persistenceCategories).Error; err != nil {
		return nil, err
	}

	// Convert to domain models
	domainCategories := make([]category.Category, len(persistenceCategories))
	for i, c := range persistenceCategories {
		domainCat := c.ToDomainCategory()
		domainCategories[i] = *domainCat
	}

	return domainCategories, nil
}
