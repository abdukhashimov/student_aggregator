package services

import (
	"context"
	"net/http"
)

type contextKey string

const (
	withErrorKey contextKey = "withError"
)

func WithError(ctx context.Context) bool {
	withError, ok := ctx.Value(withErrorKey).(bool)

	if !ok {
		return false
	}

	return withError
}

func SetWithError(r *http.Request, withError bool) *http.Request {
	ctx := context.WithValue(r.Context(), withErrorKey, withError)

	return r.WithContext(ctx)
}
