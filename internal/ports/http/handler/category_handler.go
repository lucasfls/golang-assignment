package handler

import (
	"encoding/json"
	"net/http"

	appcategory "github.com/mytheresa/go-hiring-challenge/internal/application/category"
	httputil "github.com/mytheresa/go-hiring-challenge/internal/ports/http"
)

// CategoryHandler handles HTTP requests for category operations.
type CategoryHandler struct {
	listCategoriesUC *appcategory.ListCategoriesUseCase
	createCategoryUC *appcategory.CreateCategoryUseCase
}

// NewCategoryHandler creates a new instance of CategoryHandler.
func NewCategoryHandler(listCategoriesUC *appcategory.ListCategoriesUseCase, createCategoryUC *appcategory.CreateCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		listCategoriesUC: listCategoriesUC,
		createCategoryUC: createCategoryUC,
	}
}

// HandleListCategories handles GET /categories requests.
// Returns a list of all categories.
func (h *CategoryHandler) HandleListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.listCategoriesUC.Execute(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteSuccess(w, categories)
}

// CreateCategoryRequest represents the request body for creating a category.
type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// HandleCreateCategory handles POST /categories requests.
// Accepts a JSON body with category details and creates a new category.
func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.createCategoryUC.Execute(r.Context(), req.Code, req.Name)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteCreated(w, category)
}
