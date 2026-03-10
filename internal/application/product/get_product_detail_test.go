package product

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

func TestGetProductDetailUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := GetSampleProducts()[0]
		return &prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "PROD001", res.Code)
	assert.Equal(t, "29.99", res.Price)
	assert.NotNil(t, res.Category)
	assert.Equal(t, "CLOTHING", res.Category.Code)
	assert.Equal(t, "Clothing", res.Category.Name)
}

func TestGetProductDetailUseCase_Execute_WithVariants(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := GetSampleProducts()[0]
		return &prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD001")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Variants))
	assert.Equal(t, "Small", res.Variants[0].Name)
	assert.Equal(t, "PROD001-S", res.Variants[0].SKU)
	assert.Equal(t, "25.99", res.Variants[0].Price)
}

func TestGetProductDetailUseCase_Execute_VariantInheritsProductPrice(t *testing.T) {
	// Arrange - PROD001 has variant B with null price
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := GetSampleProducts()[0]
		return &prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD001")

	// Assert
	assert.NoError(t, err)
	// Variant B (index 1) has no price, should inherit from product (29.99)
	assert.Equal(t, "Large", res.Variants[1].Name)
	assert.Equal(t, "29.99", res.Variants[1].Price) // Inherited from product price
}

func TestGetProductDetailUseCase_Execute_EmptyCode(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository()
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "product code cannot be empty", err.Error())
}

func TestGetProductDetailUseCase_Execute_ProductNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		return nil, assert.AnError
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "NONEXISTENT")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetProductDetailUseCase_Execute_WithoutVariants(t *testing.T) {
	// Arrange - PROD003 has no variants
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := GetSampleProducts()[2]
		return &prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD003")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Variants))
	assert.Equal(t, "PROD003", res.Code)
}

func TestGetProductDetailUseCase_Execute_WithDifferentCategory(t *testing.T) {
	// Arrange - PROD002 is SHOES category
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := GetSampleProducts()[1]
		return &prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD002")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res.Category)
	assert.Equal(t, "SHOES", res.Category.Code)
	assert.Equal(t, "Shoes", res.Category.Name)
}

func TestGetProductDetailUseCase_Execute_AllVariantsHavePrices(t *testing.T) {
	// Arrange - Create a product where all variants have prices
	mockRepo := NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		clothingCategory := &category.Category{
			ID:   1,
			Code: "CLOTHING",
			Name: "Clothing",
		}
		prod := &product.Product{
			ID:       1,
			Code:     "PROD001",
			Price:    MustDecimal("29.99"),
			Category: clothingCategory,
			Variants: []product.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Small",
					SKU:       "PROD001-S",
					Price:     MustDecimal("25.99"),
				},
				{
					ID:        2,
					ProductID: 1,
					Name:      "Large",
					SKU:       "PROD001-L",
					Price:     MustDecimal("35.99"),
				},
			},
		}
		return prod, nil
	})
	uc := NewGetProductDetailUseCase(mockRepo)

	// Act
	res, err := uc.Execute(context.Background(), "PROD001")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Variants))
	assert.Equal(t, "25.99", res.Variants[0].Price)
	assert.Equal(t, "35.99", res.Variants[1].Price)
}
