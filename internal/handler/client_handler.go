package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type ClientHandler struct {
	uc usecase.ClientUseCase
}

func NewClientHandler(uc usecase.ClientUseCase) *ClientHandler {
	return &ClientHandler{uc: uc}
}

func (h *ClientHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/clients")
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

func (h *ClientHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NameFull string `json:"name_full"`
		Phone    string `json:"phone"`
		Correo   string `json:"correo"`
		Address  string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	c, err := h.uc.Create(req.NameFull, req.Phone, req.Correo, req.Address)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, c, http.StatusCreated)
}

func (h *ClientHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	c, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *ClientHandler) getAll(w http.ResponseWriter, r *http.Request) {
	clients, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, clients, http.StatusOK)
}

func (h *ClientHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		NameFull string `json:"name_full"`
		Phone    string `json:"phone"`
		Correo   string `json:"correo"`
		Address  string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	c, err := h.uc.Update(id, req.NameFull, req.Phone, req.Correo, req.Address)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *ClientHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
