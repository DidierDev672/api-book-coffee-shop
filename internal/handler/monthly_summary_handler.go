package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type MonthlySummaryHandler struct {
	uc usecase.MonthlySummaryUseCase
}

func NewMonthlySummaryHandler(uc usecase.MonthlySummaryUseCase) *MonthlySummaryHandler {
	return &MonthlySummaryHandler{uc: uc}
}

func (h *MonthlySummaryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/monthly-summaries")
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

func (h *MonthlySummaryHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Product        string  `json:"product"`
		BeginningStock float64 `json:"beginning_stock"`
		IncomingOrders float64 `json:"incoming_orders"`
		OutgoingOrders float64 `json:"outgoing_orders"`
		EndingStock    float64 `json:"ending_stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ms, err := h.uc.Create(req.Product, req.BeginningStock, req.IncomingOrders, req.OutgoingOrders, req.EndingStock)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, ms, http.StatusCreated)
}

func (h *MonthlySummaryHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	ms, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, ms, http.StatusOK)
}

func (h *MonthlySummaryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, summaries, http.StatusOK)
}

func (h *MonthlySummaryHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Product        string  `json:"product"`
		BeginningStock float64 `json:"beginning_stock"`
		IncomingOrders float64 `json:"incoming_orders"`
		OutgoingOrders float64 `json:"outgoing_orders"`
		EndingStock    float64 `json:"ending_stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ms, err := h.uc.Update(id, req.Product, req.BeginningStock, req.IncomingOrders, req.OutgoingOrders, req.EndingStock)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, ms, http.StatusOK)
}

func (h *MonthlySummaryHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
