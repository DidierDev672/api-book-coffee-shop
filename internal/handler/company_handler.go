package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type CompanyHandler struct {
	uc usecase.CompanyUseCase
}

func NewCompanyHandler(uc usecase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{uc: uc}
}

func (h *CompanyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/companies")
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

type companyRequest struct {
	NIT              string `json:"nit"`
	SocialReason     string `json:"social_reason"`
	BusinessName     string `json:"business_name"`
	TypePerson       string `json:"type_person"`
	CompanyType      string `json:"company_type"`
	Status           string `json:"status"`
	ConstitutionDate string `json:"constitution_date"`
}

func (h *CompanyHandler) create(w http.ResponseWriter, r *http.Request) {
	var req companyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	c, err := h.uc.Create(
		req.NIT, req.SocialReason, req.BusinessName,
		req.TypePerson, req.CompanyType, req.Status, req.ConstitutionDate,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, c, http.StatusCreated)
}

func (h *CompanyHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	c, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *CompanyHandler) getAll(w http.ResponseWriter, r *http.Request) {
	companies, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, companies, http.StatusOK)
}

func (h *CompanyHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req companyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	c, err := h.uc.Update(
		id, req.NIT, req.SocialReason, req.BusinessName,
		req.TypePerson, req.CompanyType, req.Status, req.ConstitutionDate,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *CompanyHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
