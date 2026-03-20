package category

import (
	"context"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

type CreateCategory struct {
	repo category.Repository
}

func NewCreateCategory(repo category.Repository) *CreateCategory {
	return &CreateCategory{repo: repo}
}

func (cc *CreateCategory) Create(ctx context.Context, code, name string) (*category.Category, error) {
	cat, err := category.New(code, name)
	if err != nil {
		return nil, fmt.Errorf("validate category: %w", err)
	}

	created, err := cc.repo.Create(ctx, cat.Code, cat.Name)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return created, nil
}
