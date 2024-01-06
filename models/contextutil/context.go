package contextutil

import (
	"context"

	"github.com/twsm000/lenslocked/models/entities"
)

type ctxKey string

const (
	userKey ctxKey = "user"
)

// WithUser return a new context with user stored into it
func WithUser(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser extract the user from the context
func GetUser(ctx context.Context) (user *entities.User, ok bool) {
	user, ok = WithValueAs[*entities.User](ctx, userKey)
	return
}

// WithValueAs extract the value from the context with typesafe cast
func WithValueAs[T any](ctx context.Context, key any) (t T, ok bool) {
	t, ok = ctx.Value(key).(T)
	return
}
