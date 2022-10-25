package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

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

		tokenStr := splitAuthHeader[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}

			return s.hmacSampleSecret, nil
		})

		if err != nil {
			sendInvalidAuthTokenError(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendInvalidAuthTokenError(w)
			return
		}

		userIdClaim, ok := claims["id"]
		if !ok {
			sendInvalidAuthTokenError(w)
			return
		}

		userIdString, ok := userIdClaim.(string)
		if !ok {
			sendInvalidAuthTokenError(w)
			return
		}

		user, err := s.userService.UserById(r.Context(), userIdString)
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
