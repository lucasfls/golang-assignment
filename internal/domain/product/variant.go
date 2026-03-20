package product

import (
	"github.com/shopspring/decimal"
)

type Variant struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Name      string          `json:"name"`
	SKU       string          `json:"sku"`
	Price     decimal.Decimal `json:"price"`
}
