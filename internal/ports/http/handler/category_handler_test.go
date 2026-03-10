package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appcategory "github.com/mytheresa/go-hiring-challenge/internal/application/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
	"github.com/stretchr/testify/assert"
)

// TestCategoryHandler_HandleListCategories tests GET /categories endpoint
func TestCategoryHandler_HandleListCategories(t *testing.T) {
	// Arrange
	expectedCategories := []category.Category{
		{ID: 1, Code: "clothing", Name: "Clothing"},
		{ID: 2, Code: "shoes", Name: "Shoes"},
		{ID: 3, Code: "accessories", Name: "Accessories"},
	}

	mockRepo := appcategory.NewMockRepository().WithFindAll(func(ctx context.Context) ([]category.Category, error) {
		return expectedCategories, nil
	})

	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	handler.HandleListCategories(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var categories []category.Category
	err := json.NewDecoder(w.Body).Decode(&categories)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(categories))
	assert.Equal(t, "clothing", categories[0].Code)
	assert.Equal(t, "Clothing", categories[0].Name)
}

// TestCategoryHandler_HandleListCategories_EmptyList tests empty category list
func TestCategoryHandler_HandleListCategories_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository().WithFindAll(func(ctx context.Context) ([]category.Category, error) {
		return []category.Category{}, nil
	})

	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	handler.HandleListCategories(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var categories []category.Category
	err := json.NewDecoder(w.Body).Decode(&categories)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(categories))
}

// TestCategoryHandler_HandleCreateCategory tests POST /categories endpoint
func TestCategoryHandler_HandleCreateCategory(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository().WithCreate(func(ctx context.Context, code, name string) (*category.Category, error) {
		assert.Equal(t, "electronics", code)
		assert.Equal(t, "Electronics", name)
		return &category.Category{ID: 4, Code: code, Name: name}, nil
	})

	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	reqBody := CreateCategoryRequest{
		Code: "electronics",
		Name: "Electronics",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleCreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var createdCategory category.Category
	err := json.NewDecoder(w.Body).Decode(&createdCategory)
	assert.NoError(t, err)
	assert.Equal(t, uint(4), createdCategory.ID)
	assert.Equal(t, "electronics", createdCategory.Code)
	assert.Equal(t, "Electronics", createdCategory.Name)
}

// TestCategoryHandler_HandleCreateCategory_InvalidJSON tests invalid request body
func TestCategoryHandler_HandleCreateCategory_InvalidJSON(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository()
	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleCreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCategoryHandler_HandleCreateCategory_MissingCode tests missing code field
func TestCategoryHandler_HandleCreateCategory_MissingCode(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository()
	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	reqBody := CreateCategoryRequest{
		Code: "",
		Name: "Electronics",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleCreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCategoryHandler_HandleCreateCategory_MissingName tests missing name field
func TestCategoryHandler_HandleCreateCategory_MissingName(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository()
	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	reqBody := CreateCategoryRequest{
		Code: "electronics",
		Name: "",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleCreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCategoryHandler_HandleCreateCategory_EmptyBody tests empty request body
func TestCategoryHandler_HandleCreateCategory_EmptyBody(t *testing.T) {
	// Arrange
	mockRepo := appcategory.NewMockRepository()
	listUC := appcategory.NewListCategoriesUseCase(mockRepo)
	createUC := appcategory.NewCreateCategoryUseCase(mockRepo)
	handler := NewCategoryHandler(listUC, createUC)

	// Act
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleCreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
