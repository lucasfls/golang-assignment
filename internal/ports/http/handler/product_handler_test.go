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

// TestProductHandler_HandleListCatalog_DefaultPagination tests default pagination values
func TestProductHandler_HandleListCatalog_DefaultPagination(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 10, limit)
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

// TestProductHandler_HandleListCatalog_WithOffset tests custom offset
func TestProductHandler_HandleListCatalog_WithOffset(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		assert.Equal(t, 2, offset)
		assert.Equal(t, 10, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?offset=2", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

// TestProductHandler_HandleListCatalog_WithLimit tests custom limit
func TestProductHandler_HandleListCatalog_WithLimit(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 5, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?limit=5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

// TestProductHandler_HandleListCatalog_LimitCapAt100 tests that limit is capped at 100
func TestProductHandler_HandleListCatalog_LimitCapAt100(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Handler should cap at 100
		assert.Equal(t, 100, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?limit=200", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestProductHandler_HandleListCatalog_LimitMinimum1 tests that limit is at least 1
func TestProductHandler_HandleListCatalog_LimitMinimum1(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Handler should coerce invalid limit to minimum of 1
		assert.Greater(t, limit, 0)
		assert.LessOrEqual(t, limit, 100)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act - request with invalid limit
	req := httptest.NewRequest("GET", "/catalog?limit=0", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestProductHandler_HandleListCatalog_NegativeOffsetTreatedAsZero tests negative offset
func TestProductHandler_HandleListCatalog_NegativeOffsetTreatedAsZero(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Negative offset should be treated as 0
		assert.GreaterOrEqual(t, offset, 0)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?offset=-5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestProductHandler_HandleListCatalog_InvalidOffsetIgnored tests non-numeric offset
func TestProductHandler_HandleListCatalog_InvalidOffsetIgnored(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		// Invalid offset should default to 0
		assert.Equal(t, 0, offset)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?offset=invalid", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestProductHandler_HandleListCatalog_OffsetAndLimitTogether tests both params together
func TestProductHandler_HandleListCatalog_OffsetAndLimitTogether(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int) ([]product.Product, int64, error) {
		assert.Equal(t, 5, offset)
		assert.Equal(t, 25, limit)
		return appproduct.GetSampleProducts(), 100, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	handler := NewProductHandler(uc)

	// Act
	req := httptest.NewRequest("GET", "/catalog?offset=5&limit=25", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), res.Total)
}
