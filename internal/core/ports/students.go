package ports

import (
	"context"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type StudentsStore interface {
	SaveRSS(ctx context.Context, student domain.StudentRSS)
	SaveWAC(ctx context.Context, student domain.StudentWAC)
}
