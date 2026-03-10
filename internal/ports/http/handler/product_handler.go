package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	appproduct "github.com/mytheresa/go-hiring-challenge/internal/application/product"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

type ProductHandler struct {
	listCatalog *appproduct.ListCatalog
	getDetail   *appproduct.GetProductDetail
}

func NewProductHandler(listCatalog *appproduct.ListCatalog, getDetail *appproduct.GetProductDetail) *ProductHandler {
	return &ProductHandler{
		listCatalog: listCatalog,
		getDetail:   getDetail,
	}
}

func (h *ProductHandler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	offset := parseIntParam(r, "offset", 0)
	limit := parseIntParam(r, "limit", 10)

	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	filter := product.FindAllFilter{
		CategoryCode: r.URL.Query().Get("category"),
	}

	if maxPriceStr := r.URL.Query().Get("maxPrice"); maxPriceStr != "" {
		if price, err := appproduct.ParsePrice(maxPriceStr); err == nil {
			filter.MaxPrice = price
		}
	}

	res, err := h.listCatalog.Execute(r.Context(), offset, limit, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func parseIntParam(r *http.Request, key string, defaultVal int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed >= 0 {
			return parsed
		}
	}
	return defaultVal
}

func (h *ProductHandler) HandleGetProductDetail(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "product code is required", http.StatusBadRequest)
		return
	}

	res, err := h.getDetail.Execute(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
