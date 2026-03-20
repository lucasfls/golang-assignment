package product

import (
	"context"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

type GetProductDetail struct {
	repo product.Repository
}

func NewGetProductDetail(repo product.Repository) *GetProductDetail {
	return &GetProductDetail{repo: repo}
}

type DetailResponse struct {
	Code     string            `json:"code"`
	Price    string            `json:"price"`
	Category *CategoryResponse `json:"category"`
	Variants []VariantResponse `json:"variants"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VariantResponse struct {
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Price string `json:"price"`
}

func (gpd *GetProductDetail) Get(ctx context.Context, code string) (*DetailResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("product code cannot be empty")
	}

	prod, err := gpd.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get product detail: %w", err)
	}

	variants := make([]VariantResponse, len(prod.Variants))
	for i, v := range prod.Variants {
		price := v.Price
		if price.IsZero() {
			price = prod.Price
		}
		variants[i] = VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.String(),
		}
	}

	var catResp *CategoryResponse
	if prod.Category != nil {
		catResp = &CategoryResponse{
			Code: prod.Category.Code,
			Name: prod.Category.Name,
		}
	}

	return &DetailResponse{
		Code:     prod.Code,
		Price:    prod.Price.String(),
		Category: catResp,
		Variants: variants,
	}, nil
}
