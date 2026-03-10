package product

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID       uint               `json:"id"`
	Code     string             `json:"code"`
	Price    decimal.Decimal    `json:"price"`
	Category *category.Category `json:"category"`
	Variants []Variant          `json:"variants"`
}
