package handlers

import (
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

type UserProfileResponse struct {
	User domain.UserProfile `json:"user"`
}

// @Summary User SignIn
// @Description user sign in process
// @Tags user-auth
// @Param request body domain.SignInUserInput true "query params"
// @Success 200 {object} domain.Tokens
// @Failure 422
// @Failure 500
// @Accept json
// @Produce json
// @Router /users/login [post]
func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	input, err := inputFromContext[domain.SignInUserInput](r.Context())
	if err != nil {
		sendServerError(w, err)
		return
	}

	userId, err := s.userService.SignIn(r.Context(), *input)

	if err != nil {
		if err == domain.ErrNotFound {
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

// @Summary User SignUp
// @Description user sign up process
// @Tags user-auth
// @Param request body domain.SignUpUserInput true "query params"
// @Success 201
// @Failure 422
// @Failure 500
// @Accept json
// @Produce json
// @Router /users [post]
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	input, err := inputFromContext[domain.SignUpUserInput](r.Context())
	if err != nil {
		sendServerError(w, err)
		return
	}

	_, err = s.userService.SignUp(r.Context(), *input)
	if err != nil {
		// toDo: check other error types
		writeJSON(w, http.StatusInternalServerError, M{"message": "internal error"})
		return
	}

	sendCode(w, http.StatusCreated)
}

// @Summary User Profile
// @Description retrieves user profile
// @Security UsersAuth
// @Tags user
// @Success 200 {object} UserProfileResponse
// @Failure 401
// @Failure 500
// @Accept json
// @Produce json
// @Router /user [get]
func (s *Server) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := userFromContext(ctx)
	if err != nil {
		sendServerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, UserProfileResponse{
		User: *user.GetProfile(),
	})
}

// @Summary User Refresh Tokens
// @Description user refresh tokens process
// @Tags user-auth
// @Param request body domain.TokenInput true "query params"
// @Success 201
// @Failure 422
// @Failure 500
// @Accept json
// @Produce json
// @Router /auth/refresh [post]
func (s *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	input, err := inputFromContext[domain.TokenInput](r.Context())
	if err != nil {
		sendServerError(w, err)
		return
	}

	user, err := s.userService.UserByRefreshToken(r.Context(), input.Token)
	if err != nil {
		if err == domain.ErrNotFound {
			sendInvalidRefreshTokenError(w)
			return
		}
		sendServerError(w, err)
		return
	}

	tokens, err := s.userService.GenerateUserTokens(r.Context(), user.ID)
	if err != nil {
		sendServerError(w, err)
		return
	}

	logger.Log.Debugf("tokens successful refreshed. userId: %s", user.ID)

	writeJSON(w, http.StatusOK, M{"tokens": domain.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}})
}
