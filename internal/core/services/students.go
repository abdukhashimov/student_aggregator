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

func (s *StudentsService) ListStudents(ctx context.Context, options domain.ListStudentsOptions) ([]domain.StudentRecord, error) {
	schemas, err := s.repo.GetAll(ctx, options)

	return schemas, err
}

func (s *StudentsService) UpdateStudent(ctx context.Context, id string, input domain.StudentRecord) (*domain.StudentRecord, error) {
	err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	student, err := s.repo.GetById(ctx, id)

	return student, err
}

func (s *StudentsService) DeleteStudent(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)

	return err
}

func (s *StudentsService) DeleteStudentByFileName(ctx context.Context, fileName string) error {
	err := s.repo.DeleteByFileName(ctx, fileName)

	return err
}
