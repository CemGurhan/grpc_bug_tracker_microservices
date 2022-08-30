package middleware

import (
	"context"

	user "github.com/cemgurhan/auth-microservice/structs"
)

type contextKey string

var UserContextKey contextKey = "user"

func UserFromContext(ctx context.Context) *user.GoogleUser {
	if user, ok := ctx.Value(UserContextKey).(*user.GoogleUser); ok {
		return user
	}

	return nil
}

func WithUser(ctx context.Context, user *user.GoogleUser) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}
