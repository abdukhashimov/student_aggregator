package handlers

import (
	"github.com/rs/cors"
	"net/http"
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		noAuth.Handle("/health", healthCheck()).Methods(MethodGet)
		noAuth.Handle("/users/login", s.loginUser()).Methods(MethodPost)
	}

}

func healthCheck() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := M{
			"status":  "available",
			"message": "healthy",
			"data":    M{"hello": "world"},
		}
		writeJSON(rw, http.StatusOK, resp)
	})
}
