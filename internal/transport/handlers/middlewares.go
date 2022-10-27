package handlers

import (
	"net/http"
	"strings"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			sendInvalidAuthTokenError(w)
			return
		}

		splitAuthHeader := strings.Split(authHeader, " ")
		if len(splitAuthHeader) < 2 {
			sendInvalidAuthTokenError(w)
			return
		}

		user, err := s.userService.UserByAccessToken(r.Context(), splitAuthHeader[1])
		if err != nil {
			sendInvalidAuthTokenError(w)
			return
		}

		r = setContextUser(r, &domain.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})

		next.ServeHTTP(w, r)
	})
}
