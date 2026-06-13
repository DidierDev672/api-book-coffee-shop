package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
	"book-coffee-shop/internal/utils"
)

type AuthHandler struct {
	uc usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NameFull    string `json:"name_full"`
		Phone       string `json:"phone"`
		IDNumber    string `json:"id_number"`
		DateOfBirth string `json:"date_of_birth"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token := extractBearerToken(r)

	var u *domain.User
	err := utils.TryExecute(r.Context(), func() error {
		if err := utils.ValidateRegisterFields(req.NameFull, req.Phone, req.IDNumber, req.DateOfBirth, req.Email, req.Password); err != nil {
			return err
		}
		var ucErr error
		u, ucErr = h.uc.Register(token, req.NameFull, req.Phone, req.IDNumber, req.DateOfBirth, req.Email, req.Password)
		return ucErr
	})

	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "authorization") || strings.Contains(err.Error(), "context cancelled") {
			status = http.StatusUnauthorized
		}
		writeError(w, err.Error(), status)
		return
	}

	writeJSON(w, u, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token := extractBearerToken(r)
	u, authToken, err := h.uc.Login(token, req.Email, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "invalid email or password") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "authorization") {
			status = http.StatusUnauthorized
		}
		writeError(w, err.Error(), status)
		return
	}

	writeJSON(w, map[string]any{
		"token": authToken,
		"user":  u,
	}, http.StatusOK)
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, users, http.StatusOK)
}

func extractBearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
