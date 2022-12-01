package handlers

import (
	"net/http"

	"github.com/rs/cors"
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		noAuth.Handle("/health", http.HandlerFunc(healthCheck)).Methods(http.MethodGet)
		noAuth.Handle("/users/login", http.HandlerFunc(s.loginUser)).Methods(http.MethodPost)
		noAuth.Handle("/users", http.HandlerFunc(s.createUser)).Methods(http.MethodPost)
		noAuth.Handle("/auth/refresh", http.HandlerFunc(s.refreshToken)).Methods(http.MethodPost)
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate)
	{
		authApiRoutes.Handle("/user", http.HandlerFunc(s.getCurrentUser)).Methods(http.MethodGet)
		authApiRoutes.Handle("/schemas", http.HandlerFunc(s.createSchema)).Methods(http.MethodPost)
		authApiRoutes.Handle("/schemas", http.HandlerFunc(s.listSchemas)).Methods(http.MethodGet)
		authApiRoutes.Handle("/schemas/{id}", http.HandlerFunc(s.getSchemaById)).Methods(http.MethodGet)
		authApiRoutes.Handle("/schemas/{id}", http.HandlerFunc(s.updateSchema)).Methods(http.MethodPatch)
		authApiRoutes.Handle("/schemas/{id}", http.HandlerFunc(s.deleteSchema)).Methods(http.MethodDelete)
		authApiRoutes.Handle("/storage/upload", s.blobUpload()).Methods(http.MethodPost)
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
func healthCheck(rw http.ResponseWriter, r *http.Request) {
	resp := M{
		"status":  "available",
		"message": "healthy",
		"data":    M{"hello": "world"},
	}
	writeJSON(rw, http.StatusOK, resp)
}
