package handler

import (
	"encoding/json"
	"net/http"

	appcategory "github.com/mytheresa/go-hiring-challenge/internal/application/category"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.create.Create(r.Context(), req.Code, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}
