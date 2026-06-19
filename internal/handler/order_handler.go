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

	if id != "" && strings.HasSuffix(id, "/approve") {
		if r.Method == http.MethodPatch {
			h.approve(w, r, strings.TrimSuffix(id, "/approve"))
			return
		}
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		OrderNumeric     string                  `json:"order_numeric"`
		OrderType        string                  `json:"order_type"`
		Date             string                  `json:"date"`
		CompanyID        string                  `json:"company_id"`
		UserID           string                  `json:"user_id"`
		RequestedBy      string                  `json:"requested_by"`
		Details          []domain.OrderDetail    `json:"details"`
		FinancialSummary domain.FinancialSummary `json:"financial_summary"`
		Status           string                  `json:"status"`
		ReasonForOrder   string                  `json:"reason_for_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.Create(req.OrderNumeric, req.OrderType, req.Date, req.CompanyID, req.UserID, req.RequestedBy, req.Details, req.FinancialSummary, req.Status, req.ReasonForOrder, extractIP(r))
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
		OrderNumeric     string                  `json:"order_numeric"`
		OrderType        string                  `json:"order_type"`
		Date             string                  `json:"date"`
		CompanyID        string                  `json:"company_id"`
		UserID           string                  `json:"user_id"`
		RequestedBy      string                  `json:"requested_by"`
		Details          []domain.OrderDetail    `json:"details"`
		FinancialSummary domain.FinancialSummary `json:"financial_summary"`
		Status           string                  `json:"status"`
		ReasonForOrder   string                  `json:"reason_for_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.Update(id, req.OrderNumeric, req.OrderType, req.Date, req.CompanyID, req.UserID, req.RequestedBy, req.Details, req.FinancialSummary, req.Status, req.ReasonForOrder, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id, extractIP(r)); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) approve(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.uc.Approve(id, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, order, http.StatusOK)
}
