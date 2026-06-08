package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
)

type OrderHandler struct {
	uc usecase.OrderUseCase
}

func NewOrderHandler(uc usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/orders")
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

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderNumeric  string              `json:"order_numeric"`
		Date          string              `json:"date"`
		Hour          string              `json:"hour"`
		AttendedBy    string              `json:"attended_by"`
		ClientID      string              `json:"client_id"`
		Details       []domain.OrderDetail `json:"details"`
		PaymentMethod string              `json:"payment_method"`
		Status        string              `json:"status"`
		Observations  string              `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.Create(req.OrderNumeric, req.Date, req.Hour, req.AttendedBy, req.ClientID, req.Details, req.PaymentMethod, req.Status, req.Observations)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, order, http.StatusCreated)
}

func (h *OrderHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) getAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, orders, http.StatusOK)
}

func (h *OrderHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		OrderNumeric  string              `json:"order_numeric"`
		Date          string              `json:"date"`
		Hour          string              `json:"hour"`
		AttendedBy    string              `json:"attended_by"`
		ClientID      string              `json:"client_id"`
		Details       []domain.OrderDetail `json:"details"`
		PaymentMethod string              `json:"payment_method"`
		Status        string              `json:"status"`
		Observations  string              `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.Update(id, req.OrderNumeric, req.Date, req.Hour, req.AttendedBy, req.ClientID, req.Details, req.PaymentMethod, req.Status, req.Observations)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
