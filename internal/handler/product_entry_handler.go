package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
)

type ProductEntryHandler struct {
	uc usecase.ProductEntryUseCase
}

func NewProductEntryHandler(uc usecase.ProductEntryUseCase) *ProductEntryHandler {
	return &ProductEntryHandler{uc: uc}
}

func (h *ProductEntryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/product-entries")
	id := strings.TrimPrefix(path, "/")

	if r.Method == http.MethodGet && id == "by-product-codes" {
		h.getByProductCodes(w, r)
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

func (h *ProductEntryHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EntryNumber      string                     `json:"entry_number"`
		RegisteredDate   string                     `json:"registered_date"`
		MovementType     string                     `json:"movement_type"`
		Warehouse        string                     `json:"warehouse"`
		ResponsibleParty string                     `json:"responsible_party"`
		CompanyID        string                     `json:"company_id"`
		Details          []domain.ProductEntryDetail `json:"details"`
		FinancialSummary domain.FinancialSummary      `json:"financial_summary"`
		Observations     string                     `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pe, err := h.uc.Create(
		req.EntryNumber, req.RegisteredDate, req.MovementType,
		req.Warehouse, req.ResponsibleParty, req.CompanyID,
		req.Details, req.FinancialSummary, req.Observations,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, pe, http.StatusCreated)
}

func (h *ProductEntryHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	pe, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, pe, http.StatusOK)
}

func (h *ProductEntryHandler) getByProductCodes(w http.ResponseWriter, r *http.Request) {
	codesParam := r.URL.Query().Get("codes")
	companyID := r.URL.Query().Get("company_id")
	if codesParam == "" {
		writeError(w, "codes query parameter is required", http.StatusBadRequest)
		return
	}
	if companyID == "" {
		writeError(w, "company_id query parameter is required", http.StatusBadRequest)
		return
	}
	codes := strings.Split(codesParam, ",")

	entries, err := h.uc.GetByProductCodes(codes, companyID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, entries, http.StatusOK)
}

func (h *ProductEntryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	entries, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, entries, http.StatusOK)
}

func (h *ProductEntryHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		EntryNumber      string                     `json:"entry_number"`
		RegisteredDate   string                     `json:"registered_date"`
		MovementType     string                     `json:"movement_type"`
		Warehouse        string                     `json:"warehouse"`
		ResponsibleParty string                     `json:"responsible_party"`
		CompanyID        string                     `json:"company_id"`
		Details          []domain.ProductEntryDetail `json:"details"`
		FinancialSummary domain.FinancialSummary      `json:"financial_summary"`
		Observations     string                     `json:"observations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pe, err := h.uc.Update(
		id, req.EntryNumber, req.RegisteredDate, req.MovementType,
		req.Warehouse, req.ResponsibleParty, req.CompanyID,
		req.Details, req.FinancialSummary, req.Observations,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, pe, http.StatusOK)
}

func (h *ProductEntryHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
