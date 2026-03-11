package handler

import (
	"encoding/json"
	"net/http"

	appcategory "github.com/mytheresa/go-hiring-challenge/internal/application/category"
	"github.com/mytheresa/go-hiring-challenge/internal/ports/http/response"
)

type CategoryHandler struct {
	list   *appcategory.ListCategories
	create *appcategory.CreateCategory
}

func NewCategoryHandler(list *appcategory.ListCategories, create *appcategory.CreateCategory) *CategoryHandler {
	return &CategoryHandler{
		list:   list,
		create: create,
	}
}

func (h *CategoryHandler) HandleListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.list.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, categories)
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	category, err := h.create.Create(r.Context(), req.Code, req.Name)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.JSON(w, http.StatusCreated, category)
}
