package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
)

type WineryHandler struct {
	uc usecase.WineryUseCase
}

func NewWineryHandler(uc usecase.WineryUseCase) *WineryHandler {
	return &WineryHandler{uc: uc}
}

func (h *WineryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/wineries")
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

func (h *WineryHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegisteredDate string `json:"registered_date"`
		UserID         string `json:"user_id"`
		CompanyID      string `json:"company_id"`
		Area           string `json:"area"`
		Units          string `json:"units"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	winery, err := h.uc.Create(req.RegisteredDate, req.UserID, req.CompanyID, req.Area, req.Units)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, winery, http.StatusCreated)
}

func (h *WineryHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	winery, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, winery, http.StatusOK)
}

func (h *WineryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("company_id")
	var wineries []*domain.Winery
	var err error
	if companyID != "" {
		wineries, err = h.uc.GetByCompanyID(companyID)
	} else {
		wineries, err = h.uc.GetAll()
	}
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, wineries, http.StatusOK)
}

func (h *WineryHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		RegisteredDate string `json:"registered_date"`
		UserID         string `json:"user_id"`
		CompanyID      string `json:"company_id"`
		Area           string `json:"area"`
		Units          string `json:"units"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	winery, err := h.uc.Update(id, req.RegisteredDate, req.UserID, req.CompanyID, req.Area, req.Units)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, winery, http.StatusOK)
}

func (h *WineryHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
