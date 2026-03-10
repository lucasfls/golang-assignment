package product

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

func TestListCatalogUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		return GetSampleProducts(), 3, nil
	})

	uc := NewListCatalogUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), 0, 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Products))
	assert.Equal(t, int64(3), res.Total)
	assert.Equal(t, "PROD001", res.Products[0].Code)
	assert.Equal(t, "29.99", res.Products[0].Price.String())
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
	assert.Equal(t, "Clothing", res.Products[0].Category.Name)
}

func TestListCatalogUseCase_Execute_WithVariants(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		return GetSampleProducts(), 3, nil
	})

	uc := NewListCatalogUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), 0, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Products[0].Variants))
	assert.Equal(t, "Small", res.Products[0].Variants[0].Name)
	assert.Equal(t, "25.99", res.Products[0].Variants[0].Price.String())
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
}

func TestListCatalogUseCase_Execute_WithOffset(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Verify offset and limit are passed correctly
		assert.Equal(t, 1, offset)
		assert.Equal(t, 10, limit)
		// Return only 2 products (simulating offset=1)
		products := GetSampleProducts()
		return products[1:], 3, nil
	})

	uc := NewListCatalogUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Products))
	assert.Equal(t, int64(3), res.Total) // Total should still be 3
	assert.Equal(t, "PROD002", res.Products[0].Code)
}

func TestListCatalogUseCase_Execute_WithCustomLimit(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Verify custom limit is passed
		assert.Equal(t, 0, offset)
		assert.Equal(t, 2, limit)
		// Return limited products
		products := GetSampleProducts()
		return products[:2], 3, nil
	})

	uc := NewListCatalogUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), 0, 2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Products))
	assert.Equal(t, int64(3), res.Total)
}

func TestListCatalogUseCase_Execute_WithMaxLimit(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Use case receives already-validated limit from handler
		return GetSampleProducts(), 3, nil
	})

	uc := NewListCatalogUseCase(mockRepo)

	// Act - Handler validates before calling use case, so max of 100 is enforced there
	res, err := uc.Execute(context.Background(), 0, 100)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Products))
}
