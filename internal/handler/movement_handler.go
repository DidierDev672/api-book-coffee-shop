package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type MovementHandler struct {
	uc usecase.MovementUseCase
}

func NewMovementHandler(uc usecase.MovementUseCase) *MovementHandler {
	return &MovementHandler{uc: uc}
}

func (h *MovementHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/movements")
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

func (h *MovementHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Date           string  `json:"date"`
		Code           string  `json:"code"`
		Product        string  `json:"product"`
		Unit           string  `json:"unit"`
		Entrance       float64 `json:"entrance"`
		Output         float64 `json:"output"`
		Balance        float64 `json:"balance"`
		UnitCost       float64 `json:"unit_cost"`
		ValorValue     float64 `json:"valor_value"`
		MovementTypeID string  `json:"movement_type_id"`
		Observations   string  `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	m, err := h.uc.Create(req.Date, req.Code, req.Product, req.Unit, req.Entrance, req.Output, req.Balance, req.UnitCost, req.ValorValue, req.MovementTypeID, req.Observations, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, m, http.StatusCreated)
}

func (h *MovementHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	m, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, m, http.StatusOK)
}

func (h *MovementHandler) getAll(w http.ResponseWriter, r *http.Request) {
	movements, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, movements, http.StatusOK)
}

func (h *MovementHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Date           string  `json:"date"`
		Code           string  `json:"code"`
		Product        string  `json:"product"`
		Unit           string  `json:"unit"`
		Entrance       float64 `json:"entrance"`
		Output         float64 `json:"output"`
		Balance        float64 `json:"balance"`
		UnitCost       float64 `json:"unit_cost"`
		ValorValue     float64 `json:"valor_value"`
		MovementTypeID string  `json:"movement_type_id"`
		Observations   string  `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	m, err := h.uc.Update(id, req.Date, req.Code, req.Product, req.Unit, req.Entrance, req.Output, req.Balance, req.UnitCost, req.ValorValue, req.MovementTypeID, req.Observations, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, m, http.StatusOK)
}

func (h *MovementHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
