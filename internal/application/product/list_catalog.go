package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

type ListCatalog struct {
	repo product.Repository
}

func NewListCatalog(repo product.Repository) *ListCatalog {
	return &ListCatalog{repo: repo}
}

type Response struct {
	Products []product.Product `json:"products"`
	Total    int64             `json:"total"`
}

func (lc *ListCatalog) List(ctx context.Context, offset, limit int, filter product.FindAllFilter) (*Response, error) {
	products, total, err := lc.repo.FindAll(ctx, offset, limit, filter)
	if err != nil {
		return nil, err
	}

	return &Response{
		Products: products,
		Total:    total,
	}, nil
}
