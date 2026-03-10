package product

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ListCatalogRequest represents the input for listing products.
type ListCatalogRequest struct {
}

// ListCatalogResponse represents the output for listing products.
type ListCatalogResponse struct {
	Products []ProductDTO `json:"products"`
	Total    uint         `json:"total"`
}

// ProductDTO is the data transfer object for a product.
type ProductDTO struct {
	Code     string       `json:"code"`
	Price    float64      `json:"price"`
	Variants []VariantDTO `json:"variants"`
}

// VariantDTO is the data transfer object for a variant.
type VariantDTO struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

// ListCatalogUseCase orchestrates listing products with filters and pagination.
type ListCatalogUseCase struct {
	productRepo product.Repository
}

// NewListCatalogUseCase creates a new instance of ListCatalogUseCase.
func NewListCatalogUseCase(
	productRepo product.Repository,
) *ListCatalogUseCase {
	return &ListCatalogUseCase{
		productRepo: productRepo,
	}
}

// Execute performs the listing operation.
func (uc *ListCatalogUseCase) Execute(ctx context.Context, req ListCatalogRequest) (*ListCatalogResponse, error) {
	// Fetch all products
	products, err := uc.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Build response DTOs
	productDTOs := make([]ProductDTO, len(products))
	for i, p := range products {
		variants := make([]VariantDTO, len(p.Variants))
		productPrice := p.Price.InexactFloat64()
		for j, v := range p.Variants {
			price := productPrice // Default to product price
			if !v.Price.IsZero() {
				price = v.Price.InexactFloat64()
			}
			variants[j] = VariantDTO{
				Name:  v.Name,
				SKU:   v.SKU,
				Price: price,
			}
		}

		productDTOs[i] = ProductDTO{
			Code:     p.Code,
			Price:    productPrice,
			Variants: variants,
		}
	}

	return &ListCatalogResponse{
		Products: productDTOs,
		Total:    uint(len(productDTOs)),
	}, nil
}
