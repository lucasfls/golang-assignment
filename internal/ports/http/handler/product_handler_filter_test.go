package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appproduct "github.com/mytheresa/go-hiring-challenge/internal/application/product"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

// TestProductHandler_HandleListCatalog_FilterByCategory tests filtering by category
func TestProductHandler_HandleListCatalog_FilterByCategory(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Verify that filter is passed
		assert.Equal(t, 1, len(filters))
		assert.Equal(t, "CLOTHING", filters[0].CategoryCode)

		// Return only clothing products
		clothingProducts := []product.Product{appproduct.GetSampleProducts()[0]}
		return clothingProducts, 1, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?category=CLOTHING", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
}

// TestProductHandler_HandleListCatalog_FilterByMaxPrice tests filtering by max price
func TestProductHandler_HandleListCatalog_FilterByMaxPrice(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Verify that filter is passed
		assert.Equal(t, 1, len(filters))
		// Check that the max price is 20 (not 20.00, decimal may normalize)
		assert.True(t, filters[0].MaxPrice.Equal(appproduct.MustDecimal("20.00")))

		// Return products with price < 20.00
		cheapProducts := []product.Product{appproduct.GetSampleProducts()[0]}
		return cheapProducts, 1, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?maxPrice=20.00", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
}

// TestProductHandler_HandleListCatalog_FilterByCategoryAndPrice tests filtering by both category and price
func TestProductHandler_HandleListCatalog_FilterByCategoryAndPrice(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Verify that both filters are passed
		assert.Equal(t, 1, len(filters))
		assert.Equal(t, "SHOES", filters[0].CategoryCode)
		assert.True(t, filters[0].MaxPrice.Equal(appproduct.MustDecimal("100.00")))

		// Return filtered results
		return []product.Product{appproduct.GetSampleProducts()[1]}, 1, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?category=SHOES&maxPrice=100.00", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, "SHOES", res.Products[0].Category.Code)
}

// TestProductHandler_HandleListCatalog_FilterWithPagination tests filtering with pagination
func TestProductHandler_HandleListCatalog_FilterWithPagination(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Verify pagination and filter parameters
		assert.Equal(t, 1, offset)
		assert.Equal(t, 5, limit)
		assert.Equal(t, "CLOTHING", filters[0].CategoryCode)

		return appproduct.GetSampleProducts()[1:], 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?category=CLOTHING&offset=1&limit=5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

// TestProductHandler_HandleListCatalog_InvalidMaxPrice tests handling of invalid max price
func TestProductHandler_HandleListCatalog_InvalidMaxPrice(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Invalid price should result in zero filter
		assert.Equal(t, 1, len(filters))
		assert.Equal(t, true, filters[0].MaxPrice.IsZero())

		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?maxPrice=invalid", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

// TestProductHandler_HandleListCatalog_EmptyFilters tests no filters provided
func TestProductHandler_HandleListCatalog_EmptyFilters(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// No filters provided, both should be zero/empty
		assert.Equal(t, 1, len(filters))
		assert.Equal(t, "", filters[0].CategoryCode)
		assert.Equal(t, true, filters[0].MaxPrice.IsZero())

		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
	assert.Equal(t, 3, len(res.Products))
}
