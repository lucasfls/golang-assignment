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

func TestProductHandler_HandleListCatalog_DefaultPagination(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 10, limit)
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

func TestProductHandler_HandleListCatalog_WithOffset(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 2, offset)
		assert.Equal(t, 10, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?offset=2", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

func TestProductHandler_HandleListCatalog_WithLimit(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		assert.Equal(t, 5, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?limit=5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Total)
}

func TestProductHandler_HandleListCatalog_LimitCapAt100(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 100, limit)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?limit=200", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_HandleListCatalog_LimitMinimum1(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Greater(t, limit, 0)
		assert.LessOrEqual(t, limit, 100)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?limit=0", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_HandleListCatalog_NegativeOffsetTreatedAsZero(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.GreaterOrEqual(t, offset, 0)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?offset=-5", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_HandleListCatalog_InvalidOffsetIgnored(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 0, offset)
		return appproduct.GetSampleProducts(), 3, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?offset=invalid", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_HandleListCatalog_OffsetAndLimitTogether(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindAll(func(ctx context.Context, offset, limit int, filter product.FindAllFilter) ([]product.Product, int64, error) {
		assert.Equal(t, 5, offset)
		assert.Equal(t, 25, limit)
		return appproduct.GetSampleProducts(), 100, nil
	})
	uc := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(uc, detailUC)

	req := httptest.NewRequest("GET", "/catalog?offset=5&limit=25", nil)
	w := httptest.NewRecorder()
	handler.HandleListCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res appproduct.Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), res.Total)
}

func TestProductHandler_HandleGetProductDetail(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		assert.Equal(t, "PROD001", code)
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

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

func TestProductHandler_HandleGetProductDetail_WithVariants(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Variants))
	assert.Equal(t, "Small", res.Variants[0].Name)
	assert.Equal(t, "25.99", res.Variants[0].Price)
}

func TestProductHandler_HandleGetProductDetail_VariantInheritsPrice(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, "Large", res.Variants[1].Name)
	assert.Equal(t, "29.99", res.Variants[1].Price)
}

func TestProductHandler_HandleGetProductDetail_EmptyCode(t *testing.T) {
	mockRepo := appproduct.NewMockRepository()
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/", nil)
	req.SetPathValue("code", "")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_HandleGetProductDetail_NotFound(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		return nil, assert.AnError
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/NONEXISTENT", nil)
	req.SetPathValue("code", "NONEXISTENT")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_HandleGetProductDetail_WithoutVariants(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[2]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD003", nil)
	req.SetPathValue("code", "PROD003")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, "PROD003", res.Code)
	assert.Equal(t, 0, len(res.Variants))
}

func TestProductHandler_HandleGetProductDetail_DifferentCategory(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[1]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD002", nil)
	req.SetPathValue("code", "PROD002")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res appproduct.DetailResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotNil(t, res.Category)
	assert.Equal(t, "SHOES", res.Category.Code)
	assert.Equal(t, "Shoes", res.Category.Name)
}

func TestProductHandler_HandleGetProductDetail_IncludesAllFields(t *testing.T) {
	mockRepo := appproduct.NewMockRepository().WithFindByCode(func(ctx context.Context, code string) (*product.Product, error) {
		prod := appproduct.GetSampleProducts()[0]
		return &prod, nil
	})
	listUC := appproduct.NewListCatalog(mockRepo)
	detailUC := appproduct.NewGetProductDetail(mockRepo)
	handler := NewProductHandler(listUC, detailUC)

	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()
	handler.HandleGetProductDetail(w, req)

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
