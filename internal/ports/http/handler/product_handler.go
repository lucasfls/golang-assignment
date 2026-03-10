package handler

import (
	"net/http"
	"strconv"

	appproduct "github.com/mytheresa/go-hiring-challenge/internal/application/product"
	httputil "github.com/mytheresa/go-hiring-challenge/internal/ports/http"
)

// ProductHandler handles HTTP requests for product operations.
type ProductHandler struct {
	listCatalogUC      *appproduct.ListCatalogUseCase
	getProductDetailUC *appproduct.GetProductDetailUseCase
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(listCatalogUC *appproduct.ListCatalogUseCase, getProductDetailUC *appproduct.GetProductDetailUseCase) *ProductHandler {
	return &ProductHandler{
		listCatalogUC:      listCatalogUC,
		getProductDetailUC: getProductDetailUC,
	}
}

// HandleListCatalog handles GET /catalog requests.
// Query parameters:
//   - offset: starting position (default 0)
//   - limit: number of items to return (default 10, max 100, min 1)
//   - category: filter by category code (optional)
//   - maxPrice: filter for products with price less than this value (optional)
func (h *ProductHandler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	// Parse and validate pagination parameters
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	limit := 10 // default
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	// Validate limit bounds
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	// Parse optional filters
	filter := appproduct.NewFilter()
	if category := r.URL.Query().Get("category"); category != "" {
		filter.SetCategory(category)
	}
	if maxPrice := r.URL.Query().Get("maxPrice"); maxPrice != "" {
		// Attempt to parse maxPrice as decimal
		if price, err := appproduct.ParseDecimal(maxPrice); err == nil {
			filter.SetMaxPrice(price)
		}
	}

	// Execute use case with filters
	res, err := h.listCatalogUC.Execute(r.Context(), offset, limit, filter.ToFilterOptions()...)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return JSON response
	httputil.WriteSuccess(w, res)
}

// HandleGetProductDetail handles GET /catalog/:code requests.
// Returns product details including variants and category.
func (h *ProductHandler) HandleGetProductDetail(w http.ResponseWriter, r *http.Request) {
	// Extract product code from URL path parameter
	// Go 1.22+ supports path patterns like "GET /catalog/{code}"
	code := r.PathValue("code")
	if code == "" {
		httputil.WriteError(w, http.StatusBadRequest, "product code is required")
		return
	}

	// Execute use case
	res, err := h.getProductDetailUC.Execute(r.Context(), code)
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Return JSON response
	httputil.WriteSuccess(w, res)
}
