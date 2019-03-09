package middleware

import "context"

// Naming convention for middleware is noun with -er, -or or -ar.

type contextKey string

func (c contextKey) WithValue(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, c, v)
}

func (c contextKey) Value(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(c).(string)
	return v, ok
}
