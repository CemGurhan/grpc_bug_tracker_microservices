package middleware

import (
	"context"

	user "github.com/cemgurhan/auth-microservice/structs"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

// UserContextKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
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
