package handlers

import (
	"github.com/rs/cors"
	"net/http"
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		noAuth.Handle("/health", healthCheck()).Methods(http.MethodGet)
		noAuth.Handle("/users/login", s.loginUser()).Methods(http.MethodPost)
		noAuth.Handle("/users", s.createUser()).Methods(http.MethodPost)
		noAuth.Handle("/auth/refresh", s.refreshToken()).Methods(http.MethodPost)
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate)
	{
		authApiRoutes.Handle("/user", s.getCurrentUser()).Methods(http.MethodGet)
	}
}

// @Summary Health Check
// @Description Health Check
// @Tags health
// @Success 200
// @Failure 500
// @Accept json
// @Produce json
// @Router /health [get]
func healthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		resp := M{
			"status":  "available",
			"message": "healthy",
			"data":    M{"hello": "world"},
		}
		writeJSON(rw, http.StatusOK, resp)
	}
}
