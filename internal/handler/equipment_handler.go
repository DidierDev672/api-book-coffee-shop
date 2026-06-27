package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"book-coffee-shop/internal/usecase"
)

type EquipmentHandler struct {
	uc usecase.EquipmentUseCase
}

func NewEquipmentHandler(uc usecase.EquipmentUseCase) *EquipmentHandler {
	return &EquipmentHandler{uc: uc}
}

func (h *EquipmentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/equipment")
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

func (h *EquipmentHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		Status          string `json:"status"`
		LastMaintenance string `json:"lastMaintenance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var lastMaint time.Time
	if req.LastMaintenance != "" {
		var err error
		lastMaint, err = time.Parse("2006-01-02", req.LastMaintenance)
		if err != nil {
			lastMaint, err = time.Parse(time.RFC3339, req.LastMaintenance)
			if err != nil {
				writeError(w, "invalid lastMaintenance format, use YYYY-MM-DD or RFC3339", http.StatusBadRequest)
				return
			}
		}
	}

	equipment, err := h.uc.Create(req.Name, req.Type, req.Status, lastMaint)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, equipment, http.StatusCreated)
}

func (h *EquipmentHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	equipment, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, equipment, http.StatusOK)
}

func (h *EquipmentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	equipments, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, equipments, http.StatusOK)
}

func (h *EquipmentHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		Status          string `json:"status"`
		LastMaintenance string `json:"lastMaintenance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var lastMaint time.Time
	if req.LastMaintenance != "" {
		var err error
		lastMaint, err = time.Parse("2006-01-02", req.LastMaintenance)
		if err != nil {
			lastMaint, err = time.Parse(time.RFC3339, req.LastMaintenance)
			if err != nil {
				writeError(w, "invalid lastMaintenance format, use YYYY-MM-DD or RFC3339", http.StatusBadRequest)
				return
			}
		}
	}

	equipment, err := h.uc.Update(id, req.Name, req.Type, req.Status, lastMaint)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, equipment, http.StatusOK)
}

func (h *EquipmentHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
