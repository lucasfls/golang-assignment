package handler

import (
	"encoding/json"
	"net/http"

	appproduct "github.com/mytheresa/go-hiring-challenge/internal/application/product"
)

// ProductHandler handles HTTP requests for product operations.
type ProductHandler struct {
	listCatalogUC *appproduct.ListCatalogUseCase
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(listCatalogUC *appproduct.ListCatalogUseCase) *ProductHandler {
	return &ProductHandler{
		listCatalogUC: listCatalogUC,
	}
}

// HandleListCatalog handles GET /catalog requests.
func (h *ProductHandler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	// Execute use case
	req := appproduct.ListCatalogRequest{}

	res, err := h.listCatalogUC.Execute(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
