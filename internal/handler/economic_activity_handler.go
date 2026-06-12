package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type EconomicActivityHandler struct {
	uc usecase.EconomicActivityUseCase
}

func NewEconomicActivityHandler(uc usecase.EconomicActivityUseCase) *EconomicActivityHandler {
	return &EconomicActivityHandler{uc: uc}
}

func (h *EconomicActivityHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/economic-activities")
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

type economicActivityRequest struct {
	UserID      string `json:"user_id"`
	CompanyID   string `json:"company_id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (h *EconomicActivityHandler) create(w http.ResponseWriter, r *http.Request) {
	var req economicActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	a, err := h.uc.Create(req.UserID, req.CompanyID, req.Code, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, a, http.StatusCreated)
}

func (h *EconomicActivityHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	a, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *EconomicActivityHandler) getAll(w http.ResponseWriter, r *http.Request) {
	activities, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, activities, http.StatusOK)
}

func (h *EconomicActivityHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req economicActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	a, err := h.uc.Update(id, req.UserID, req.CompanyID, req.Code, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *EconomicActivityHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
