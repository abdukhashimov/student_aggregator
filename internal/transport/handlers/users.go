package handlers

import (
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"net/http"
)

func (s *Server) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := domain.SignInUserInput{}

		if err := readJSON(r.Body, &input); err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		tokens, err := s.userService.SignIn(r.Context(), input)

		if err != nil || tokens == nil {
			sendUnauthorizedError(w, input.Email)
			return
		}

		writeJSON(w, http.StatusOK, M{"tokens": tokens})
	}
}

func (s *Server) getCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := userFromContext(ctx)
		if err != nil {
			sendServerError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"user": user})
	}
}
