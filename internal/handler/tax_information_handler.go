package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type TaxInformationHandler struct {
	uc usecase.TaxInformationUseCase
}

func NewTaxInformationHandler(uc usecase.TaxInformationUseCase) *TaxInformationHandler {
	return &TaxInformationHandler{uc: uc}
}

func (h *TaxInformationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/tax-information")
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

type taxInformationRequest struct {
	UserID              string `json:"user_id"`
	BusinessID          string `json:"business_id"`
	TaxRegime           string `json:"tax_regime"`
	VATResponsible      bool   `json:"vat_responsible"`
	WithholdingTaxpayer bool   `json:"withholding_taxpayer"`
	LargeTaxpayer       bool   `json:"large_taxpayer"`
}

func (h *TaxInformationHandler) create(w http.ResponseWriter, r *http.Request) {
	var req taxInformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	t, err := h.uc.Create(
		req.UserID, req.BusinessID, req.TaxRegime,
		req.VATResponsible, req.WithholdingTaxpayer, req.LargeTaxpayer,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, t, http.StatusCreated)
}

func (h *TaxInformationHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	t, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, t, http.StatusOK)
}

func (h *TaxInformationHandler) getAll(w http.ResponseWriter, r *http.Request) {
	records, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, records, http.StatusOK)
}

func (h *TaxInformationHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req taxInformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	t, err := h.uc.Update(
		id, req.UserID, req.BusinessID, req.TaxRegime,
		req.VATResponsible, req.WithholdingTaxpayer, req.LargeTaxpayer,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, t, http.StatusOK)
}

func (h *TaxInformationHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
