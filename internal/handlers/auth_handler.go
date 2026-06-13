package handlers

import (
	"net/http"
	"strings"

	"book-coffee-shop/internal/middleware"
	"book-coffee-shop/internal/models"
	"book-coffee-shop/internal/usecase"
	"book-coffee-shop/internal/utils"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(useCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	validationMiddleware := middleware.ValidatePayload(&req)
	validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := middleware.ExtractBearerToken(r)

		u, authToken, err := h.useCase.Login(token, req.Email, req.Password)
		if err != nil {
			status := h.determineErrorStatus(err)
			utils.WriteError(w, err.Error(), status)
			return
		}

		utils.WriteJSON(w, map[string]any{
			"token": authToken,
			"user":  u,
		}, http.StatusOK)
	})).ServeHTTP(w, r)
}

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	authWrapper := middleware.RequireRoles("admin", "user")
	authWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok || user == nil {
			utils.WriteError(w, "authentication required", http.StatusUnauthorized)
			return
		}

		profile, err := h.useCase.GetProfile(user.ID)
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusNotFound)
			return
		}

		utils.WriteJSON(w, profile, http.StatusOK)
	})).ServeHTTP(w, r)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	wrapper := middleware.RequireRoles("admin")
	wrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, map[string]string{"status": "deleted"}, http.StatusOK)
	})).ServeHTTP(w, r)
}

func (h *AuthHandler) determineErrorStatus(err error) int {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "invalid email or password"):
		return http.StatusUnauthorized
	case strings.Contains(errMsg, "token") || strings.Contains(errMsg, "authorization"):
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}
