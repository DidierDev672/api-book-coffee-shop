package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"book-coffee-shop/internal/usecase"
)

type ExerciseHandler struct {
	uc usecase.ExerciseUseCase
}

func NewExerciseHandler(uc usecase.ExerciseUseCase) *ExerciseHandler {
	return &ExerciseHandler{uc: uc}
}

func (h *ExerciseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/exercises")
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
			writeError(w, "id is required", http.StatusBadRequest, "")
			return
		}
		h.update(w, r, id)
	case http.MethodDelete:
		if id == "" {
			writeError(w, "id is required", http.StatusBadRequest, "")
			return
		}
		h.delete(w, r, id)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed, "")
	}
}

func (h *ExerciseHandler) HandleByEquipmentID(w http.ResponseWriter, r *http.Request, equipmentID string) {
	w.Header().Set("Content-Type", "application/json")

	if h == nil || h.uc == nil {
		writeError(w, "exercise handler not initialized", http.StatusInternalServerError, "")
		return
	}

	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	exercises, err := h.uc.GetByEquipmentID(equipmentID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "equipment not found" {
			status = http.StatusNotFound
		}
		writeError(w, err.Error(), status, "")
		return
	}
	writeJSON(w, exercises, http.StatusOK)
}

func (h *ExerciseHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EquipmentID string `json:"equipment_id"`
		Name        string `json:"name"`
		MuscleGroup string `json:"muscle_group"`
		Difficulty  string `json:"difficulty"`
		VideoURL    string `json:"video_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest, "")
		return
	}

	exercise, err := h.uc.Create(req.EquipmentID, req.Name, req.MuscleGroup, req.Difficulty, req.VideoURL)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	writeJSON(w, exercise, http.StatusCreated)
}

func (h *ExerciseHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	exercise, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound, "")
		return
	}
	writeJSON(w, exercise, http.StatusOK)
}

func (h *ExerciseHandler) getAll(w http.ResponseWriter, r *http.Request) {
	muscleGroup := r.URL.Query().Get("muscle_group")
	difficulty := r.URL.Query().Get("difficulty")

	var exercises interface{}
	var err error

	if muscleGroup != "" || difficulty != "" {
		exercises, err = h.uc.GetAllFiltered(muscleGroup, difficulty)
	} else {
		exercises, err = h.uc.GetAll()
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError, "")
		return
	}
	writeJSON(w, exercises, http.StatusOK)
}

func (h *ExerciseHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		EquipmentID string `json:"equipment_id"`
		Name        string `json:"name"`
		MuscleGroup string `json:"muscle_group"`
		Difficulty  string `json:"difficulty"`
		VideoURL    string `json:"video_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest, "")
		return
	}

	exercise, err := h.uc.Update(id, req.EquipmentID, req.Name, req.MuscleGroup, req.Difficulty, req.VideoURL)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}
	writeJSON(w, exercise, http.StatusOK)
}

func (h *ExerciseHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleUseCaseError(w http.ResponseWriter, err error) {
	var fieldErr *usecase.FieldError
	if errors.As(err, &fieldErr) {
		writeError(w, fieldErr.Message, http.StatusBadRequest, fieldErr.Field)
		return
	}
	if strings.Contains(err.Error(), "not found") {
		writeError(w, err.Error(), http.StatusNotFound, "")
		return
	}
	writeError(w, err.Error(), http.StatusBadRequest, "")
}
