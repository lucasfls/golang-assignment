package product

import (
	"github.com/shopspring/decimal"
)

// Product is the aggregate root for the product bounded context.
// It represents a product in the catalog with its variants and category.
type Product struct {
	ID         uint
	Code       string
	Price      decimal.Decimal
	CategoryID uint
	Variants   []Variant
}

// GetEffectivePrice returns the product price, with fallback to variant price if available.
func (p *Product) GetEffectivePrice(variantID uint) decimal.Decimal {
	for _, v := range p.Variants {
		if v.ID == variantID && !v.Price.IsZero() {
			return v.Price
		}
	}
	return p.Price
}

// NewProduct creates a new Product aggregate root.
func NewProduct(code string, price decimal.Decimal, categoryID uint) *Product {
	return &Product{
		Code:       code,
		Price:      price,
		CategoryID: categoryID,
		Variants:   make([]Variant, 0),
	}
}
