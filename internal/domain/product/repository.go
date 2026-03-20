package product

import (
	"context"

	"github.com/shopspring/decimal"
)

type FindAllFilter struct {
	CategoryCode string
	MaxPrice     decimal.Decimal
}

type Repository interface {
	FindAll(ctx context.Context, offset, limit int, filter FindAllFilter) ([]Product, int64, error)
	FindByCode(ctx context.Context, code string) (*Product, error)
}
