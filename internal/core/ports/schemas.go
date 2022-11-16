package ports

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type SchemaService interface {
	NewSchema(ctx context.Context, input domain.NewSchemaInput) (*domain.Schema, error)
	ListSchemas(ctx context.Context) ([]domain.Schema, error)
	GetSchemaById(ctx context.Context, id string) (*domain.Schema, error)
	UpdateSchema(ctx context.Context, id string, input domain.UpdateSchemaInput) (*domain.Schema, error)
	DeleteSchema(ctx context.Context, id string) error
}

type SchemaStore interface {
	Create(ctx context.Context, input domain.Schema) (string, error)
	FindAll(ctx context.Context) ([]domain.Schema, error)
	GetById(ctx context.Context, id string) (*domain.Schema, error)
	Update(ctx context.Context, id string, input domain.UpdateSchemaInput) error
	Delete(ctx context.Context, id string) error
}
