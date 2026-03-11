package category

import (
	"context"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

type ListCategories struct {
	repo category.Repository
}

func NewListCategories(repo category.Repository) *ListCategories {
	return &ListCategories{repo: repo}
}

func (lc *ListCategories) List(ctx context.Context) ([]category.Category, error) {
	categories, err := lc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}
