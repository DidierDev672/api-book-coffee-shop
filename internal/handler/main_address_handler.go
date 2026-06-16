package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type MainAddressHandler struct {
	uc usecase.MainAddressUseCase
}

func NewMainAddressHandler(uc usecase.MainAddressUseCase) *MainAddressHandler {
	return &MainAddressHandler{uc: uc}
}

func (h *MainAddressHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/main-addresses")
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

type mainAddressRequest struct {
	UserID     string `json:"user_id"`
	CompanyID  string `json:"company_id"`
	Country    string `json:"country"`
	Department string `json:"department"`
	Municipio  string `json:"municipio"`
	Address    string `json:"address"`
	Postcode   string `json:"postcode"`
}

func (h *MainAddressHandler) create(w http.ResponseWriter, r *http.Request) {
	var req mainAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	a, err := h.uc.Create(
		req.UserID, req.CompanyID, req.Country, req.Department, req.Municipio, req.Address, req.Postcode,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, a, http.StatusCreated)
}

func (h *MainAddressHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	a, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *MainAddressHandler) getAll(w http.ResponseWriter, r *http.Request) {
	addresses, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, addresses, http.StatusOK)
}

func (h *MainAddressHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req mainAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	a, err := h.uc.Update(
		id, req.UserID, req.CompanyID, req.Country, req.Department, req.Municipio, req.Address, req.Postcode,
	)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *MainAddressHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
