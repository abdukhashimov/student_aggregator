package handlers

import (
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

func (s *Server) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := domain.SignInUserInput{}

		if err := readJSON(r.Body, &input); err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		userId, err := s.userService.SignIn(r.Context(), input)

		if err != nil {
			if err == domain.ErrUserNotFound {
				sendUnauthorizedError(w, input.Email)
				return
			}
			sendServerError(w, err)
			return
		}

		tokens, err := s.userService.GenerateUserTokens(r.Context(), userId)
		if err != nil {
			sendServerError(w, err)
			return
		}

		logger.Log.Debugf("user successful logged in. userId: %s", userId)

		writeJSON(w, http.StatusOK, M{"tokens": domain.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}})
	}
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := domain.SignUpUserInput{}

		err := readJSON(r.Body, &input)
		if err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		//ToDo: add input validation

		err = s.userService.SignUp(r.Context(), input)
		if err != nil {
			// toDo: check other error types
			writeJSON(w, http.StatusInternalServerError, M{"message": "internal error"})
			return
		}

		sendCode(w, http.StatusCreated)
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
