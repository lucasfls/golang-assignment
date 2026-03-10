package persistence

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ToDomainProduct converts a persistence Product model to a domain Product.
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
func (c *Category) ToDomainCategory() *category.Category {
	return &category.Category{
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

	var persistenceCategory Category
	if p.Category != nil {
		persistenceCategory = *FromDomainCategory(p.Category)
	}

	return &Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    p.Price,
		Category: persistenceCategory,
		Variants: variants,
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
func FromDomainCategory(c *category.Category) *Category {
	return &Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
