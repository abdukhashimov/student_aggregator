package handlers

import (
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
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
