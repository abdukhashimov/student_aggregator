package services

import (
	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/core/repository"
)

type Services struct {
	Users   ports.UsersService
	Schemas ports.SchemaService
	Storage ports.StorageService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	usersService := NewUsersService(repos.Users, cfg)
	schemasService := NewSchemaService(repos.Schemas, cfg)
	storageService := NewStorageService(cfg)

	return &Services{
		Users:   usersService,
		Schemas: schemasService,
		Storage: storageService,
	}
}
