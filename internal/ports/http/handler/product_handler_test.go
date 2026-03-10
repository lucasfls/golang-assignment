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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 10, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 2, offset)
		assert.Equal(t, 10, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 5, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Handler should cap at 100
		assert.Equal(t, 100, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Handler should coerce invalid limit to minimum of 1
		assert.Greater(t, limit, 0)
		assert.LessOrEqual(t, limit, 100)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Negative offset should be treated as 0
		assert.GreaterOrEqual(t, offset, 0)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		// Invalid offset should default to 0
		assert.Equal(t, 0, offset)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filters ...product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 5, offset)
		assert.Equal(t, 25, limit)
		return appproduct.GetSampleProducts(), 100, nil
	})
	uc := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(uc, detailUC)

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

// TestProductHandler_HandleGetProductDetail tests GET /catalog/:code endpoint
func TestProductHandler_HandleGetProductDetail(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		assert.Equal(t, "PROD001", code)
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, "PROD001", res.Code)
	assert.Equal(t, "29.99", res.Price)
	assert.NotNil(t, res.Category)
	assert.Equal(t, "CLOTHING", res.Category.Code)
}

// TestProductHandler_HandleGetProductDetail_WithVariants tests product with variants
func TestProductHandler_HandleGetProductDetail_WithVariants(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Variants))
	assert.Equal(t, "Small", res.Variants[0].Name)
	assert.Equal(t, "25.99", res.Variants[0].Price)
}

// TestProductHandler_HandleGetProductDetail_VariantInheritsPrice tests variant price inheritance
func TestProductHandler_HandleGetProductDetail_VariantInheritsPrice(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	// Variant at index 1 should inherit product price
	assert.Equal(t, "Large", res.Variants[1].Name)
	assert.Equal(t, "29.99", res.Variants[1].Price)
}

// TestProductHandler_HandleGetProductDetail_EmptyCode tests empty product code
func TestProductHandler_HandleGetProductDetail_EmptyCode(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository()
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/", nil)
	req.SetPathValue("code", "")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestProductHandler_HandleGetProductDetail_NotFound tests product not found
func TestProductHandler_HandleGetProductDetail_NotFound(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		return nil, assert.AnError
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/NONEXISTENT", nil)
	req.SetPathValue("code", "NONEXISTENT")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestProductHandler_HandleGetProductDetail_WithoutVariants tests product without variants
func TestProductHandler_HandleGetProductDetail_WithoutVariants(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[2]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD003", nil)
	req.SetPathValue("code", "PROD003")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, "PROD003", res.Code)
	assert.Equal(t, 0, len(res.Variants))
}

// TestProductHandler_HandleGetProductDetail_DifferentCategory tests product with different category
func TestProductHandler_HandleGetProductDetail_DifferentCategory(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[1]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD002", nil)
	req.SetPathValue("code", "PROD002")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotNil(t, res.Category)
	assert.Equal(t, "SHOES", res.Category.Code)
	assert.Equal(t, "Shoes", res.Category.Name)
}

// TestProductHandler_HandleGetProductDetail_IncludesAllFields tests all response fields
func TestProductHandler_HandleGetProductDetail_IncludesAllFields(t *testing.T) {
	// Arrange
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalogUseCase(mockRepo)
	detailUC := appproduct.NewGetProductDetailUseCase(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	// Act
	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Code)
	assert.NotEmpty(t, res.Price)
	assert.NotNil(t, res.Category)
	assert.NotEmpty(t, res.Category.Code)
	assert.NotEmpty(t, res.Category.Name)
}
