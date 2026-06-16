package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type ProviderHandler struct {
	uc usecase.ProviderUseCase
}

func NewProviderHandler(uc usecase.ProviderUseCase) *ProviderHandler {
	return &ProviderHandler{uc: uc}
}

func (h *ProviderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/providers")
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

type providerRequest struct {
	Code              string `json:"code"`
	TypePerson        string `json:"type_person"`
	DocumentType      string `json:"document_type"`
	DocumentNumber    string `json:"document_number"`
	VerificationDigit string `json:"verification_digit"`
	BusinessName      string `json:"business_name"`
	BusinessActivity  string `json:"business_activity"`
	Status            bool   `json:"status"`
}

func (h *ProviderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req providerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.uc.Create(
		req.Code, req.TypePerson, req.DocumentType,
		req.DocumentNumber, req.VerificationDigit,
		req.BusinessName, req.BusinessActivity, req.Status,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, p, http.StatusCreated)
}

func (h *ProviderHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	p, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *ProviderHandler) getAll(w http.ResponseWriter, r *http.Request) {
	providers, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, providers, http.StatusOK)
}

func (h *ProviderHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req providerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.uc.Update(
		id, req.Code, req.TypePerson, req.DocumentType,
		req.DocumentNumber, req.VerificationDigit,
		req.BusinessName, req.BusinessActivity, req.Status,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *ProviderHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
