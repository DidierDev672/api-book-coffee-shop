package handler

import (
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type InventoryHistoryHandler struct {
	svc *usecase.HistoryService
}

func NewInventoryHistoryHandler(svc *usecase.HistoryService) *InventoryHistoryHandler {
	return &InventoryHistoryHandler{svc: svc}
}

func (h *InventoryHistoryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/history")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			parts := strings.SplitN(id, "/", 2)
			if len(parts) == 2 {
				h.getByDocument(w, r, parts[0], parts[1])
				return
			}
		}
		h.getAll(w, r)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHistoryHandler) getByDocument(w http.ResponseWriter, r *http.Request, documentType, documentID string) {
	events, err := h.svc.GetByDocument(documentType, documentID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, events, http.StatusOK)
}

func (h *InventoryHistoryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.svc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, events, http.StatusOK)
}
