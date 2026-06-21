package middleware

import (
	"log"
	"net/http"
	"strings"

	"book-coffee-shop/internal/repository"
	"book-coffee-shop/internal/utils"
)

func NewAuthMiddleware(tokenService repository.TokenService, userRepo repository.UserRepository, publicPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			path := strings.TrimRight(r.URL.Path, "/")
			for _, p := range publicPaths {
				if path == strings.TrimRight(p, "/") {
					next.ServeHTTP(w, r)
					return
				}
			}

			token := ExtractBearerToken(r)
			if token == "" {
				utils.WriteError(w, "authorization token is required", http.StatusUnauthorized)
				return
			}

			userID, err := tokenService.Validate(token)
			if err != nil {
				log.Printf("[auth] token validation failed for %s: %v", r.URL.Path, err)
				utils.WriteError(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				utils.WriteError(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := utils.SetUserContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
