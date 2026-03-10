package category

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

type CreateCategory struct {
	repo category.Repository
}

func NewCreateCategory(repo category.Repository) *CreateCategory {
	return &CreateCategory{repo: repo}
}

func (uc *CreateCategory) Execute(ctx context.Context, code, name string) (*category.Category, error) {
	cat, err := category.New(code, name)
	if err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, cat.Code, cat.Name)
}
