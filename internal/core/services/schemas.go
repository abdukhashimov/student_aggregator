package services

import (
	"context"
	"regexp"
	"strings"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

var _ ports.SchemaService = (*SchemaService)(nil)

type SchemaService struct {
	repo ports.SchemaStore
	cfg  *config.Config
}

func NewSchemaService(repo ports.SchemaStore, cfg *config.Config) *SchemaService {
	return &SchemaService{
		repo: repo,
		cfg:  cfg,
	}
}

func (ss *SchemaService) NewSchema(ctx context.Context, input domain.NewSchemaInput) (*domain.Schema, error) {
	schemaId, err := ss.repo.Create(ctx, domain.Schema{
		Name:       input.Name,
		Slug:       getSlug(input.Name),
		Version:    input.Version,
		SchemaType: input.SchemaType,
		Headers:    input.Headers,
		Fields:     input.Fields,
	})
	if err != nil {
		return nil, err
	}

	schema, err := ss.repo.GetById(ctx, schemaId)

	return schema, err
}

func (ss *SchemaService) ListSchemas(ctx context.Context) ([]domain.Schema, error) {
	schemas, err := ss.repo.FindAll(ctx)

	return schemas, err
}

func (ss *SchemaService) GetSchemaById(ctx context.Context, id string) (*domain.Schema, error) {
	schema, err := ss.repo.GetById(ctx, id)

	return schema, err
}

func (ss *SchemaService) UpdateSchema(ctx context.Context, id string, input domain.UpdateSchemaInput) (*domain.Schema, error) {
	input.Slug = nil
	if input.Name != nil {
		slug := getSlug(*input.Name)
		input.Slug = &slug
	}
	err := ss.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	schema, err := ss.repo.GetById(ctx, id)

	return schema, err
}

func (ss *SchemaService) DeleteSchema(ctx context.Context, id string) error {
	err := ss.repo.Delete(ctx, id)

	return err
}

func getSlug(in string) string {
	space := regexp.MustCompile(`\s+`)
	result := space.ReplaceAllString(in, " ")
	result = strings.TrimSpace(result)
	result = space.ReplaceAllString(result, "-")
	result = strings.ToLower(result)

	return result
}
