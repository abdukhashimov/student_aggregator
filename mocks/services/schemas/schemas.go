package schemas

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/mocks/services"
)

const (
	ValidSchemaID1   = "1"
	ValidSchemaID2   = "2"
	NotFoundSchemaID = "999999999"
)

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

var InternalError = errors.New("internal error")

var _ ports.SchemaService = (*mockSchemasService)(nil)

type mockSchemasService struct {
	schemasStorage map[string]*domain.Schema
	lastSchemaId   string
	mutex          *sync.RWMutex
}

func NewMockSchemasService() *mockSchemasService {
	return &mockSchemasService{
		schemasStorage: map[string]*domain.Schema{
			ValidSchemaID1: copySchema(&EtalonSchema1),
			ValidSchemaID2: copySchema(&EtalonSchema2),
		},
		lastSchemaId: ValidSchemaID2,
		mutex:        &sync.RWMutex{},
	}
}

func (m *mockSchemasService) NewSchema(ctx context.Context, input domain.NewSchemaInput) (*domain.Schema, error) {
	if services.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	slug := getSlug(input.Name)
	if m.isDuplicate(slug, "") {
		return nil, domain.DuplicationError
	}

	m.incrementId()
	m.schemasStorage[m.lastSchemaId] = &domain.Schema{
		ID:         m.lastSchemaId,
		Name:       input.Name,
		Slug:       slug,
		Version:    input.Version,
		SchemaType: input.SchemaType,
		Headers:    input.Headers,
		Fields:     input.Fields,
	}

	schemaCopy := copySchema(m.schemasStorage[m.lastSchemaId])

	return schemaCopy, nil
}

func (m *mockSchemasService) ListSchemas(ctx context.Context) ([]domain.Schema, error) {
	if services.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var result []domain.Schema
	for _, v := range m.schemasStorage {
		schemaCopy := copySchema(v)
		result = append(result, *schemaCopy)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

func (m *mockSchemasService) GetSchemaById(ctx context.Context, id string) (*domain.Schema, error) {
	if services.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	schema, ok := m.schemasStorage[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	schemaCopy := copySchema(schema)

	return schemaCopy, nil
}

func (m *mockSchemasService) UpdateSchema(ctx context.Context, id string, input domain.UpdateSchemaInput) (*domain.Schema, error) {
	if services.WithError(ctx) {
		return nil, InternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	schema, ok := m.schemasStorage[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	if input.Name != nil {
		slug := getSlug(*input.Name)
		if m.isDuplicate(slug, id) {
			return nil, domain.DuplicationError
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

	schemaCopy := copySchema(schema)

	return schemaCopy, nil
}

func (m *mockSchemasService) DeleteSchema(ctx context.Context, id string) error {
	if services.WithError(ctx) {
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

func (m *mockSchemasService) incrementId() {
	id, _ := strconv.Atoi(m.lastSchemaId)
	id++
	m.lastSchemaId = strconv.Itoa(id)
}

func (m *mockSchemasService) isDuplicate(slug string, id string) bool {
	for _, s := range m.schemasStorage {
		if s.Slug == slug && s.ID != id {
			return true
		}
	}

	return false
}

func getSlug(in string) string {
	space := regexp.MustCompile(`\s+`)
	result := space.ReplaceAllString(in, " ")
	result = strings.TrimSpace(result)
	result = space.ReplaceAllString(result, "-")
	result = strings.ToLower(result)

	return result
}

func copySchema(schema *domain.Schema) *domain.Schema {
	schemaCopy := *schema
	schemaCopy.Fields = make([]domain.FieldSchema, len(schema.Fields))
	copy(schemaCopy.Fields, schema.Fields)

	return &schemaCopy
}
