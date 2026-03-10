package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ListCatalogUseCase orchestrates listing products with pagination.
type ListCatalogUseCase struct {
	productRepo product.Repository
}

// NewListCatalogUseCase creates a new instance of ListCatalogUseCase.
func NewListCatalogUseCase(productRepo product.Repository) *ListCatalogUseCase {
	return &ListCatalogUseCase{
		productRepo: productRepo,
	}
}

// Response represents the catalog listing response with products and pagination info.
type Response struct {
	Products []product.Product `json:"products"`
	Total    int64             `json:"total"`
}

// Execute retrieves products with pagination.
func (uc *ListCatalogUseCase) Execute(ctx context.Context, offset, limit int) (*Response, error) {
	products, total, err := uc.productRepo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return &Response{
		Products: products,
		Total:    total,
	}, nil
}
