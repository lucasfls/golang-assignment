package category

import (
	"context"
	"errors"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

// CreateCategoryUseCase handles the business logic for creating a new category.
type CreateCategoryUseCase struct {
	repo category.Repository
}

// NewCreateCategoryUseCase creates a new instance of CreateCategoryUseCase.
func NewCreateCategoryUseCase(repo category.Repository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		repo: repo,
	}
}

// Execute creates a new category with the provided code and name.
func (uc *CreateCategoryUseCase) Execute(ctx context.Context, code, name string) (*category.Category, error) {
	// Validate input
	cat := category.NewCategory(code, name)
	if !cat.IsValid() {
		return nil, errors.New("invalid category: code and name are required")
	}

	return uc.repo.Create(ctx, code, name)
}
