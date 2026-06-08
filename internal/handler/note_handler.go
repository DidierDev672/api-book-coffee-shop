package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type NoteHandler struct {
	uc usecase.NoteUseCase
}

func NewNoteHandler(uc usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{uc: uc}
}

func (h *NoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/notes")
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

func (h *NoteHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Color   string `json:"color"`
		IDTopic string `json:"id_topic"`
		IDBook  string `json:"id_book"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	note, err := h.uc.Create(req.Name, req.Content, req.Type, req.Color, req.IDTopic, req.IDBook)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, note, http.StatusCreated)
}

func (h *NoteHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	note, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, note, http.StatusOK)
}

func (h *NoteHandler) getAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, notes, http.StatusOK)
}

func (h *NoteHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Color   string `json:"color"`
		IDTopic string `json:"id_topic"`
		IDBook  string `json:"id_book"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	note, err := h.uc.Update(id, req.Name, req.Content, req.Type, req.Color, req.IDTopic, req.IDBook)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, note, http.StatusOK)
}

func (h *NoteHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
