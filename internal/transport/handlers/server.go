package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/core/repository/mongodb"
	"github.com/abdukhashimov/student_aggregator/internal/core/services"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	server      *http.Server
	router      *mux.Router
	userService ports.UsersService
	config      *config.Config
}

func NewServer(db *mongo.Database, cfg *config.Config) *Server {
	s := Server{
		server: &http.Server{
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
		router: mux.NewRouter().StrictSlash(true),
		config: cfg,
	}

	repos := mongodb.NewRepositories(db)
	logger.Log.Info("repositories successfully initialized")

	servs := services.NewServices(repos, cfg)
	s.userService = servs.Users
	logger.Log.Info("services successfully initialized")

	s.routes()
	s.server.Handler = s.router
	logger.Log.Info("handlers successfully initialized")

	return &s
}

func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	s.server.Addr = port
	logger.Log.Infof("server starting on %s", port)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	_ = s.server.Shutdown(ctx)
}
