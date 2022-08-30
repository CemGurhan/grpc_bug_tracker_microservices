package middleware

import (
	"context"

	user "github.com/cemgurhan/auth-microservice/structs"
)

type contextKey string

var userContextKey contextKey = "user"

func WithUser(ctx context.Context, user *user.GoogleUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
