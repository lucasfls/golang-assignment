package product

// Category represents a product category.
// Categories are used to organize and filter products in the catalog.
type Category struct {
	ID   uint
	Code string
	Name string
}

// NewCategory creates a new Category.
func NewCategory(code, name string) *Category {
	return &Category{
		Code: code,
		Name: name,
	}
}

// IsValid checks if the category has required fields.
func (c *Category) IsValid() bool {
	return c.Code != "" && c.Name != ""
}
