package ports

import (
	"context"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type StudentsService interface {
	GetStudentById(ctx context.Context, id string) (*domain.StudentRecord, error)
}

type StudentsStore interface {
	SaveRSS(ctx context.Context, email string, student domain.StudentRSS) (string, error)
	SaveWAC(ctx context.Context, email string, student domain.StudentWAC) (string, error)
	GetById(ctx context.Context, id string) (*domain.StudentRecord, error)
}
