package category

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

type ListCategories struct {
	repo category.Repository
}

func NewListCategories(repo category.Repository) *ListCategories {
	return &ListCategories{repo: repo}
}

func (lc *ListCategories) List(ctx context.Context) ([]category.Category, error) {
	return lc.repo.FindAll(ctx)
}
