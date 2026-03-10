package category

import "fmt"

type Category struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func New(code, name string) (*Category, error) {
	if code == "" {
		return nil, fmt.Errorf("category code cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("category name cannot be empty")
	}

	return &Category{
		Code: code,
		Name: name,
	}, nil
}
