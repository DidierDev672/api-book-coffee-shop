package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type AuthorHandler struct {
	uc usecase.AuthorUseCase
}

func NewAuthorHandler(uc usecase.AuthorUseCase) *AuthorHandler {
	return &AuthorHandler{uc: uc}
}

func (h *AuthorHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/authors")
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
			http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
			return
		}
		h.update(w, r, id)
	case http.MethodDelete:
		if id == "" {
			http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
			return
		}
		h.delete(w, r, id)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *AuthorHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string   `json:"name"`
		Country  string   `json:"country"`
		Genres   []string `json:"genres"`
		BirthDay string   `json:"birthDay"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	author, err := h.uc.Create(req.Name, req.Country, req.Genres, req.BirthDay)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, author, http.StatusCreated)
}

func (h *AuthorHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	author, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, author, http.StatusOK)
}

func (h *AuthorHandler) getAll(w http.ResponseWriter, r *http.Request) {
	authors, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, authors, http.StatusOK)
}

func (h *AuthorHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name     string   `json:"name"`
		Country  string   `json:"country"`
		Genres   []string `json:"genres"`
		BirthDay string   `json:"birthDay"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	author, err := h.uc.Update(id, req.Name, req.Country, req.Genres, req.BirthDay)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, author, http.StatusOK)
}

func (h *AuthorHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
