package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type contextKey string

const (
	userKey contextKey = "user"
)

func setContextUser(r *http.Request, u *domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, u)
	return r.WithContext(ctx)
}

func userFromContext(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(userKey).(*domain.User)

	if !ok {
		return nil, errors.New("can not get user from context")
	}

	return user, nil
}
