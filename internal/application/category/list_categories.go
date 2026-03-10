package category

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

// ListCategoriesUseCase handles the business logic for listing all categories.
type ListCategoriesUseCase struct {
	repo category.Repository
}

// NewListCategoriesUseCase creates a new instance of ListCategoriesUseCase.
func NewListCategoriesUseCase(repo category.Repository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		repo: repo,
	}
}

// Execute retrieves all categories.
func (uc *ListCategoriesUseCase) Execute(ctx context.Context) ([]category.Category, error) {
	return uc.repo.FindAll(ctx)
}
