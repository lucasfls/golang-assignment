package handler

import (
	"fmt"
	"net/http"
	"strconv"

	appproduct "github.com/mytheresa/go-hiring-challenge/internal/application/product"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/mytheresa/go-hiring-challenge/internal/ports/http/response"
	"github.com/shopspring/decimal"
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
		if price, err := decimal.NewFromString(maxPriceStr); err == nil {
			filter.MaxPrice = price
		}
	}

	res, err := h.listCatalog.List(r.Context(), offset, limit, filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, res)
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
		response.Error(w, http.StatusBadRequest, fmt.Errorf("product code is required"))
		return
	}

	res, err := h.getDetail.Get(r.Context(), code)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	response.JSON(w, http.StatusOK, res)
}
