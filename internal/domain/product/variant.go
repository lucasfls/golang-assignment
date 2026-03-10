package product

import (
	"github.com/shopspring/decimal"
)

// Variant is an entity within the Product aggregate.
// It represents different configurations or options for a product.
type Variant struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Name      string          `json:"name"`
	SKU       string          `json:"sku"`
	Price     decimal.Decimal `json:"price"`
}

// NewVariant creates a new Variant entity.
func NewVariant(productID uint, name, sku string, price decimal.Decimal) *Variant {
	return &Variant{
		ProductID: productID,
		Name:      name,
		SKU:       sku,
		Price:     price,
	}
}
