package product

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

func TestListCatalogUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context) ([]product.Product, error) {
		return GetSampleProducts(), nil
	})

	uc := NewListCatalogUseCase(mockRepo)
	req := ListCatalogRequest{}

	// Act
	res, err := uc.Execute(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Products))
	assert.Equal(t, uint(3), res.Total)
	assert.Equal(t, "PROD001", res.Products[0].Code)
	assert.Equal(t, float64(29.99), res.Products[0].Price)
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
	assert.Equal(t, "Clothing", res.Products[0].Category.Name)
}

func TestListCatalogUseCase_Execute_WithVariants(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context) ([]product.Product, error) {
		return GetSampleProducts(), nil
	})

	uc := NewListCatalogUseCase(mockRepo)
	req := ListCatalogRequest{}

	// Act
	res, err := uc.Execute(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Products[0].Variants))
	assert.Equal(t, "Small", res.Products[0].Variants[0].Name)
	assert.Equal(t, float64(25.99), res.Products[0].Variants[0].Price)
	assert.Equal(t, float64(29.99), res.Products[0].Variants[1].Price) // Inherits product price
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
}
