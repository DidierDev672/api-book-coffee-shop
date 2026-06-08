package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type BookHandler struct {
	uc usecase.BookUseCase
}

func NewBookHandler(uc usecase.BookUseCase) *BookHandler {
	return &BookHandler{uc: uc}
}

func (h *BookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/books")
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

func (h *BookHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Author          string   `json:"author"`
		Genres          []string `json:"genres"`
		Photos          []string `json:"photos"`
		PublicationDate string   `json:"publicationDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.uc.Create(req.Title, req.Description, req.Author, req.Genres, req.Photos, req.PublicationDate)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, book, http.StatusCreated)
}

func (h *BookHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	book, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, book, http.StatusOK)
}

func (h *BookHandler) getAll(w http.ResponseWriter, r *http.Request) {
	books, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, books, http.StatusOK)
}

func (h *BookHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Author          string   `json:"author"`
		Genres          []string `json:"genres"`
		Photos          []string `json:"photos"`
		PublicationDate string   `json:"publicationDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.uc.Update(id, req.Title, req.Description, req.Author, req.Genres, req.Photos, req.PublicationDate)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, book, http.StatusOK)
}

func (h *BookHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
