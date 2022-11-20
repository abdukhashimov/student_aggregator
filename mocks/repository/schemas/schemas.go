package users

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/mocks/utils"
)

const (
	ValidSchemaID1   = "507f1f77bcf86cd000000001"
	ValidSchemaID2   = "507f1f77bcf86cd000000002"
	NotFoundSchemaID = "507f1f77bcf86cd999999999"
)

var InternalError = errors.New("internal error")

var EtalonSchema1 = domain.Schema{
	ID:         ValidSchemaID1,
	Name:       "RSS",
	Slug:       "rss",
	Version:    "1.0.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []domain.FieldSchema{
		{Name: "first_name", Col: "A"},
		{Name: "last_name", Col: "B"},
		{Name: "email", Col: "C"},
	},
}

var EtalonSchema2 = domain.Schema{
	ID:         ValidSchemaID2,
	Name:       "WAC",
	Slug:       "wac",
	Version:    "1.0.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []domain.FieldSchema{
		{Name: "name", Col: "A"},
		{Name: "surname", Col: "B"},
		{Name: "email", Col: "C"},
	},
}

var _ ports.SchemaStore = (*mockSchemasRepository)(nil)

type mockSchemasRepository struct {
	schemasStorage map[string]*domain.Schema
	lastSchemaId   string
	mutex          *sync.RWMutex
}

func NewMockSchemasRepository() *mockSchemasRepository {
	return &mockSchemasRepository{
		schemasStorage: map[string]*domain.Schema{
			ValidSchemaID1: utils.CopySchema(&EtalonSchema1),
			ValidSchemaID2: utils.CopySchema(&EtalonSchema2),
		},
		lastSchemaId: EtalonSchema2.ID,
		mutex:        &sync.RWMutex{},
	}
}

func (m *mockSchemasRepository) Create(ctx context.Context, input domain.Schema) (string, error) {
	if utils.WithError(ctx) {
		return "", InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	slug := utils.GetSlug(input.Name)
	if m.isDuplicate(slug, "") {
		return "", domain.DuplicationError
	}

	newId := utils.IncrementMongoId(m.lastSchemaId)
	m.lastSchemaId = newId

	newSchema := &domain.Schema{
		ID:         newId,
		Name:       input.Name,
		Slug:       slug,
		Version:    input.Version,
		SchemaType: input.SchemaType,
		Headers:    input.Headers,
		Fields:     input.Fields,
	}
	m.schemasStorage[newId] = newSchema

	return newId, nil
}

func (m *mockSchemasRepository) FindAll(ctx context.Context) ([]domain.Schema, error) {
	if utils.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var result []domain.Schema
	for _, v := range m.schemasStorage {
		schemaCopy := utils.CopySchema(v)
		result = append(result, *schemaCopy)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (m *mockSchemasRepository) GetById(ctx context.Context, id string) (*domain.Schema, error) {
	if utils.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	schema, ok := m.schemasStorage[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	schemaCopy := utils.CopySchema(schema)

	return schemaCopy, nil
}

func (m *mockSchemasRepository) Update(ctx context.Context, id string, input domain.UpdateSchemaInput) error {
	if utils.WithError(ctx) {
		return InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	schema, ok := m.schemasStorage[id]
	if !ok {
		return domain.ErrNotFound
	}

	if input.Name != nil {
		slug := utils.GetSlug(*input.Name)
		if m.isDuplicate(slug, id) {
			return domain.DuplicationError
		}
		schema.Name = *input.Name
		schema.Slug = slug
	}

	if input.Version != nil {
		schema.Version = *input.Version
	}

	if input.SchemaType != nil {
		schema.SchemaType = *input.SchemaType
	}

	if input.Headers != nil {
		schema.Headers = *input.Headers
	}

	if input.Fields != nil {
		schema.Fields = *input.Fields
	}

	return nil
}

func (m *mockSchemasRepository) Delete(ctx context.Context, id string) error {
	if utils.WithError(ctx) {
		return InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.schemasStorage[id]; !ok {
		return domain.ErrNotFound
	}

	delete(m.schemasStorage, id)

	return nil
}

func (m *mockSchemasRepository) isDuplicate(slug string, id string) bool {
	for _, s := range m.schemasStorage {
		if s.Slug == slug && s.ID != id {
			return true
		}
	}

	return false
}
