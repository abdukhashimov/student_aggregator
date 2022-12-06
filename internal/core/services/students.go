package services

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

var _ ports.StudentsService = (*StudentsService)(nil)

type StudentsService struct {
	repo ports.StudentsStore
	cfg  *config.Config
}

func NewStudentsService(repo ports.StudentsStore, cfg *config.Config) *StudentsService {
	return &StudentsService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *StudentsService) GetStudentById(ctx context.Context, id string) (*domain.StudentRecord, error) {
	student, err := s.repo.GetById(ctx, id)

	return student, err
}
