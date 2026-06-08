package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type ProductHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/products")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getByID(w, r, id)
		} else {
			h.getAll(w, r)
		}
	case http.MethodPost:
		h.create(w, r)
	case http.MethodPut:
		if id == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}
		h.update(w, r, id)
	case http.MethodDelete:
		if id == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}
		h.delete(w, r, id)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductCode  string   `json:"product_code"`
		Categories   []string `json:"categories"`
		Unit         string   `json:"unit"`
		MinimumStock float64  `json:"minimum_stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.uc.Create(req.ProductCode, req.Categories, req.Unit, req.MinimumStock)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, p, http.StatusCreated)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	p, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, products, http.StatusOK)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		ProductCode  string   `json:"product_code"`
		Categories   []string `json:"categories"`
		Unit         string   `json:"unit"`
		MinimumStock float64  `json:"minimum_stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.uc.Update(id, req.ProductCode, req.Categories, req.Unit, req.MinimumStock)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
