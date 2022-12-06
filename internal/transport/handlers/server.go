package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abdukhashimov/student_aggregator/docs"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/core/repository/mongodb"
	"github.com/abdukhashimov/student_aggregator/internal/core/services"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	server            *http.Server
	router            *mux.Router
	userService       ports.UsersService
	schemasService    ports.SchemaService
	storageService    ports.StorageService
	aggregatorService ports.AggregatorService
	config            *config.Config
}

func NewServer(db *mongo.Database, storageClient *minio.Client, cfg *config.Config) *Server {
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
	s.schemasService = servs.Schemas
	s.aggregatorService = servs.Aggregator

	s.storageService = servs.Storage
	s.storageService.SetClient(storageClient)

	logger.Log.Info("services successfully initialized")

	s.routes()
	s.server.Handler = s.router
	logger.Log.Info("handlers successfully initialized")

	return &s
}

// @title Student Aggregator API
// @description This API contains the source for the Student Aggregator app

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// @BasePath /api/v1

// Run initializes http server
func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	portInt, err := strconv.Atoi(port[1:])
	if err != nil {
		panic(err)
	}

	if s.config.Project.SwaggerEnabled {
		docs.SwaggerInfo.Version = s.config.Project.Version

		s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", portInt)),
		)).Methods(http.MethodGet)

		logger.Log.Infof("swagger is enabled")
	}

	s.server.Addr = port
	logger.Log.Infof("server starting on %s", port)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	_ = s.server.Shutdown(ctx)
}
