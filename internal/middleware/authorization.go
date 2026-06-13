package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/utils"
)

// RequireRoles es un wrapper de orden superior para validar permisos
func RequireRoles(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Obtener usuario del contexto (inyectado por AuthMiddleware)
			user, ok := utils.GetUserFromContext(r.Context())
			if !ok || user == nil {
				utils.WriteError(w, "authentication required", http.StatusUnauthorized)
				return
			}

			// Verificar si el usuario tiene algún rol permitido
			hasRole := false
			for _, userRole := range user.Roles {
				if slices.Contains(allowedRoles, userRole) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.WriteError(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			// Usuario autorizado - continuar
			next(w, r)
		}
	}
}

// RequirePermission wrapper más granular por acción específica
func RequirePermission(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := utils.GetUserFromContext(r.Context())
			if !ok || user == nil {
				utils.WriteError(w, "authentication required", http.StatusUnauthorized)
				return
			}

			if !user.HasPermission(permission) {
				utils.WriteError(w, "missing required permission: "+permission, http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}

func ExtractBearerToken(r *http.Request) string {
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

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	return utils.GetUserFromContext(ctx)
}
