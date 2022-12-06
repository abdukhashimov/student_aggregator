package services

import (
	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/core/repository"
)

type Services struct {
	Users      ports.UsersService
	Schemas    ports.SchemaService
	Students   ports.StudentsService
	Storage    ports.StorageService
	Aggregator ports.AggregatorService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	usersService := NewUsersService(repos.Users, cfg)
	schemasService := NewSchemaService(repos.Schemas, cfg)
	studentsService := NewStudentsService(repos.Students, cfg)
	storageService := NewStorageService(cfg)
	parserService := NewAggregatorService(repos.Students, repos.Schemas, storageService)

	return &Services{
		Users:      usersService,
		Schemas:    schemasService,
		Students:   studentsService,
		Storage:    storageService,
		Aggregator: parserService,
	}
}
