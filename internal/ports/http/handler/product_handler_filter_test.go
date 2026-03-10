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

func TestProductHandler_HandleListCatalog_FilterByCategory(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, "CLOTHING", filter.CategoryCode)
		clothingProducts := []product.Product{appproduct.GetSampleProducts()[0]}
		return clothingProducts, 1, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?category=CLOTHING", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "CLOTHING", res.Products[0].Category.Code)
}

func TestProductHandler_HandleListCatalog_FilterByMaxPrice(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.True(t, filter.MaxPrice.Equal(appproduct.MustDecimal("20.00")))
		cheapProducts := []product.Product{appproduct.GetSampleProducts()[0]}
		return cheapProducts, 1, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?maxPrice=20.00", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
}

func TestProductHandler_HandleListCatalog_FilterByCategoryAndPrice(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, "SHOES", filter.CategoryCode)
		assert.True(t, filter.MaxPrice.Equal(appproduct.MustDecimal("100.00")))
		return []product.Product{appproduct.GetSampleProducts()[1]}, 1, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?category=SHOES&maxPrice=100.00", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, "SHOES", res.Products[0].Category.Code)
}

func TestProductHandler_HandleListCatalog_FilterWithPagination(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 1, offset)
		assert.Equal(t, 5, limit)
		assert.Equal(t, "CLOTHING", filter.CategoryCode)
		return appproduct.GetSampleProducts()[1:], 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?category=CLOTHING&offset=1&limit=5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

func TestProductHandler_HandleListCatalog_InvalidMaxPrice(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, true, filter.MaxPrice.IsZero())
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?maxPrice=invalid", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

func TestProductHandler_HandleListCatalog_EmptyFilters(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, "", filter.CategoryCode)
		assert.Equal(t, true, filter.MaxPrice.IsZero())
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
	assert.Equal(t, 3, len(res.Products))
}
