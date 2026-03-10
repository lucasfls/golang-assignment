package product

import (
	"context"
	"errors"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// GetProductDetailUseCase orchestrates retrieving product details by code.
type GetProductDetailUseCase struct {
	productRepo product.Repository
}

// NewGetProductDetailUseCase creates a new instance of GetProductDetailUseCase.
func NewGetProductDetailUseCase(productRepo product.Repository) *GetProductDetailUseCase {
	return &GetProductDetailUseCase{
		productRepo: productRepo,
	}
}

// DetailResponse represents the product detail response.
type DetailResponse struct {
	Code     string              `json:"code"`
	Price    string              `json:"price"`
	Category *CategoryResponse   `json:"category"`
	Variants []VariantResponse   `json:"variants"`
}

// CategoryResponse represents a category in the detail response.
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// VariantResponse represents a variant in the detail response.
type VariantResponse struct {
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Price string `json:"price"`
}

// Execute retrieves product details by code, enriching variants with inherited prices.
func (uc *GetProductDetailUseCase) Execute(ctx context.Context, code string) (*DetailResponse, error) {
	if code == "" {
		return nil, errors.New("product code cannot be empty")
	}

	prod, err := uc.productRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if prod == nil {
		return nil, errors.New("product not found")
	}

	// Build response with enriched variants
	variants := make([]VariantResponse, len(prod.Variants))
	for i, v := range prod.Variants {
		price := v.Price
		// If variant has no price, inherit from product
		if price.IsZero() {
			price = prod.Price
		}
		variants[i] = VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.String(),
		}
	}

	catResp := (*CategoryResponse)(nil)
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
