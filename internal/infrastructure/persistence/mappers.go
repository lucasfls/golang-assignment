package persistence

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ToDomainProduct converts a persistence Product model to a domain Product.
func (p *Product) ToDomainProduct() *product.Product {
	variants := make([]product.Variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = v.ToDomainVariant()
	}

	return &product.Product{
		ID:         p.ID,
		Code:       p.Code,
		Price:      p.Price,
		CategoryID: p.CategoryID,
		Variants:   variants,
	}
}

// ToDomainVariant converts a persistence Variant model to a domain Variant.
func (v *Variant) ToDomainVariant() product.Variant {
	return product.Variant{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		SKU:       v.SKU,
		Price:     v.Price,
	}
}

// ToDomainCategory converts a persistence Category model to a domain Category.
func (c *Category) ToDomainCategory() *product.Category {
	return &product.Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}

// FromDomainProduct converts a domain Product to a persistence Product model.
func FromDomainProduct(p *product.Product) *Product {
	variants := make([]Variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = FromDomainVariant(v)
	}

	return &Product{
		ID:         p.ID,
		Code:       p.Code,
		Price:      p.Price,
		CategoryID: p.CategoryID,
		Variants:   variants,
	}
}

// FromDomainVariant converts a domain Variant to a persistence Variant model.
func FromDomainVariant(v product.Variant) Variant {
	return Variant{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		SKU:       v.SKU,
		Price:     v.Price,
	}
}

// FromDomainCategory converts a domain Category to a persistence Category model.
func FromDomainCategory(c *product.Category) *Category {
	return &Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
