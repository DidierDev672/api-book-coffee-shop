package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type EstablishmentHandler struct {
	uc usecase.EstablishmentUseCase
}

func NewEstablishmentHandler(uc usecase.EstablishmentUseCase) *EstablishmentHandler {
	return &EstablishmentHandler{uc: uc}
}

func (h *EstablishmentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/establishments")
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

func (h *EstablishmentHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EstablishmentName    string `json:"establishment_name"`
		InventoryManager     string `json:"inventory_manager"`
		WarehousePointOfSale string `json:"warehouse_point_of_sale"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	e, err := h.uc.Create(req.EstablishmentName, req.InventoryManager, req.WarehousePointOfSale)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, e, http.StatusCreated)
}

func (h *EstablishmentHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	e, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, e, http.StatusOK)
}

func (h *EstablishmentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	establishments, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, establishments, http.StatusOK)
}

func (h *EstablishmentHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		EstablishmentName    string `json:"establishment_name"`
		InventoryManager     string `json:"inventory_manager"`
		WarehousePointOfSale string `json:"warehouse_point_of_sale"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	e, err := h.uc.Update(id, req.EstablishmentName, req.InventoryManager, req.WarehousePointOfSale)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, e, http.StatusOK)
}

func (h *EstablishmentHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
