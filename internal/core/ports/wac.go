package ports

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type WACService interface {
	AddStudent(ctx context.Context, student domain.StudentWAC) (id string, err error)
	GetStudentById(ctx context.Context, id string) (*domain.StudentWAC, error)
	UpdateStudent(ctx context.Context, id string, input domain.StudentWAC) (*domain.StudentWAC, error)
	DeleteStudent(ctx context.Context, id string) error
}
