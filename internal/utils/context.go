package utils

import (
	"context"

	"book-coffee-shop/internal/domain"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}

func SetUserContext(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}
