package services

import (
	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/core/repository/mongodb"
)

type Services struct {
	Users ports.UsersService
}

func NewServices(repos *mongodb.Repositories, cfg *config.Config) *Services {
	usersService := NewUsersService(repos.Users, cfg)

	return &Services{
		Users: usersService,
	}
}
