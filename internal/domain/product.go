package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID       uint
	Code     string
	Price    decimal.Decimal
	Variants []Variant
}

type Variant struct {
	ID        uint
	ProductID uint
	Name      string
}

type ProductRepository interface {
	FindAll(ctx context.Context) ([]Product, error)
}
