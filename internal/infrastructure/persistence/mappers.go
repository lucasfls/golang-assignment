package persistence

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

func (p *Product) ToDomainProduct() *product.Product {
	variants := make([]product.Variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = v.ToDomainVariant()
	}

	var domainCategory *category.Category
	if p.Category.ID != 0 {
		domainCategory = p.Category.ToDomainCategory()
	}

	return &product.Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    p.Price,
		Category: domainCategory,
		Variants: variants,
	}
}

func (v *Variant) ToDomainVariant() product.Variant {
	return product.Variant{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		SKU:       v.SKU,
		Price:     v.Price,
	}
}

func (c *Category) ToDomainCategory() *category.Category {
	return &category.Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
