package handlers

import (
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"net/http"

	"github.com/rs/cors"
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		// health
		noAuth.Handle("/health", http.HandlerFunc(healthCheck)).Methods(http.MethodGet)
		// user
		noAuth.Handle("/users/login", validatorWrapper[domain.SignInUserInput](s.loginUser)).Methods(http.MethodPost)
		noAuth.Handle("/users", validatorWrapper[domain.SignUpUserInput](s.createUser)).Methods(http.MethodPost)
		noAuth.Handle("/auth/refresh", validatorWrapper[domain.TokenInput](s.refreshToken)).Methods(http.MethodPost)
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate)
	{
		// user
		authApiRoutes.Handle("/user", http.HandlerFunc(s.getCurrentUser)).Methods(http.MethodGet)
		// schema
		authApiRoutes.Handle("/schemas", validatorWrapper[domain.NewSchemaInput](s.createSchema)).Methods(http.MethodPost)
		authApiRoutes.Handle("/schemas", http.HandlerFunc(s.listSchemas)).Methods(http.MethodGet)
		authApiRoutes.Handle("/schemas/{id}", http.HandlerFunc(s.getSchemaById)).Methods(http.MethodGet)
		authApiRoutes.Handle("/schemas/{id}", validatorWrapper[domain.UpdateSchemaInput](s.updateSchema)).Methods(http.MethodPatch)
		authApiRoutes.Handle("/schemas/{id}", http.HandlerFunc(s.deleteSchema)).Methods(http.MethodDelete)
		// storage
		authApiRoutes.Handle("/storage/upload", s.blobUpload()).Methods(http.MethodPost)
		// aggregator
		authApiRoutes.Handle("/aggregator/parse", validatorWrapper[domain.ParseFileInput](s.parseFile)).Methods(http.MethodPost)
		// student
		authApiRoutes.Handle("/students/{id}", http.HandlerFunc(s.getStudentById)).Methods(http.MethodGet)

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
