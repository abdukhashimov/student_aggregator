package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks/repository/schemas"
	"github.com/abdukhashimov/student_aggregator/mocks/utils"
)

var newSchemaInput = domain.NewSchemaInput{
	Name:       "NewSchemaName",
	Version:    "1.0.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []domain.FieldSchema{
		{Name: "name", Col: "A"},
		{Name: "surname", Col: "B"},
		{Name: "email", Col: "C"},
	},
}

var newSchemaName = "NewSchemaName"

type SchemasTestCaseGroup struct {
	name          string
	executeMethod func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error
	testCases     []SchemasTestCase
}

type SchemasTestCase struct {
	name          string
	inputID       string
	input         interface{}
	getContext    func(ctx context.Context) context.Context
	expectedError error
	postCheck     func(s *SchemaService) error
}

var schemasTestCaseGroup = []SchemasTestCaseGroup{
	{
		name: "NewSchema",
		executeMethod: func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error {
			in, ok := input.(domain.NewSchemaInput)
			if !ok {
				return errors.New("invalid input type")
			}
			newSchema, err := s.NewSchema(ctx, in)

			if expectedError != nil {
				if expectedError != err {
					return errors.New("expected an error")
				}

				return nil
			}

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if newSchema.ID != schemas.NextSchemaID {
				return errors.New("invalid schema ID")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:  "success",
				input: newSchemaInput,
				postCheck: func(s *SchemaService) error {
					schema, err := s.repo.GetById(context.Background(), schemas.NextSchemaID)
					if err != nil {
						return fmt.Errorf("unexpecting error: %s", err.Error())
					}
					// TODO: check all fields
					if schema.Name != newSchemaInput.Name {
						return errors.New("data is not stored")
					}

					return nil
				},
			},
			{
				name: "duplicate",
				input: domain.NewSchemaInput{
					Name: schemas.EtalonSchema1.Name,
				},
				expectedError: domain.DuplicationError,
				postCheck: func(s *SchemaService) error {
					_, err := s.repo.GetById(context.Background(), schemas.NextSchemaID)
					if err != domain.ErrNotFound {
						return errors.New("schema should not be created")
					}

					return nil
				},
			},
			{
				name:          "internalError",
				input:         newSchemaInput,
				expectedError: schemas.InternalError,
				getContext: func(ctx context.Context) context.Context {
					return utils.SetWithErrorToContext(ctx, true)
				},
				postCheck: func(s *SchemaService) error {
					_, err := s.repo.GetById(context.Background(), schemas.NextSchemaID)
					if err != domain.ErrNotFound {
						return errors.New("schema should not be created")
					}

					return nil
				},
			},
		},
	},
	{
		name: "ListSchemas",
		executeMethod: func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error {
			schemasList, err := s.ListSchemas(ctx)

			if expectedError != nil {
				if expectedError != err {
					return errors.New("expected an error")
				}

				return nil
			}

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if !reflect.DeepEqual(schemasList, []domain.Schema{schemas.EtalonSchema1, schemas.EtalonSchema2}) {
				return errors.New("unexpected schema list")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name: "success",
			},
			{
				name:          "internalError",
				expectedError: schemas.InternalError,
				getContext: func(ctx context.Context) context.Context {
					return utils.SetWithErrorToContext(ctx, true)
				},
			},
		},
	},
	{
		name: "GetSchemaById",
		executeMethod: func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error {
			schema, err := s.GetSchemaById(ctx, inputID)

			if expectedError != nil {
				if expectedError != err {
					return errors.New("expected an error")
				}

				return nil
			}

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if schema.ID != inputID {
				return errors.New("invalid schema ID")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: schemas.ValidSchemaID1,
			},
			{
				name:          "notFound",
				inputID:       schemas.NotFoundSchemaID,
				expectedError: domain.ErrNotFound,
			},
			{
				name:          "internalError",
				inputID:       schemas.ValidSchemaID1,
				expectedError: schemas.InternalError,
				getContext: func(ctx context.Context) context.Context {
					return utils.SetWithErrorToContext(ctx, true)
				},
			},
		},
	},
	{
		name: "UpdateSchema",
		executeMethod: func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error {
			in, ok := input.(domain.UpdateSchemaInput)
			if !ok {
				return errors.New("invalid input type")
			}
			updatedSchema, err := s.UpdateSchema(ctx, inputID, in)

			if expectedError != nil {
				if expectedError != err {
					return errors.New("expected an error")
				}

				return nil
			}

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			// TODO: check all fields
			if updatedSchema.Name != *in.Name {
				return errors.New("schema is not updated")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: schemas.ValidSchemaID1,
				input: domain.UpdateSchemaInput{
					Name: &newSchemaName,
				},
				postCheck: func(s *SchemaService) error {
					schema, err := s.repo.GetById(context.Background(), schemas.ValidSchemaID1)
					if err != nil {
						return fmt.Errorf("unexpecting error: %s", err.Error())
					}
					if schema.Name != newSchemaInput.Name {
						return errors.New("data is not stored")
					}

					return nil
				},
			},
			{
				name:    "duplicate",
				inputID: schemas.ValidSchemaID1,
				input: domain.UpdateSchemaInput{
					Name: &schemas.EtalonSchema2.Name,
				},
				expectedError: domain.DuplicationError,
				postCheck: func(s *SchemaService) error {
					schema, err := s.repo.GetById(context.Background(), schemas.ValidSchemaID1)
					if err != nil {
						return fmt.Errorf("unexpecting error: %s", err.Error())
					}
					if schema.Name != schemas.EtalonSchema1.Name {
						return errors.New("schema should not be updated")
					}

					return nil
				},
			},
			{
				name:    "notFound",
				inputID: schemas.NotFoundSchemaID,
				input: domain.UpdateSchemaInput{
					Name: &newSchemaName,
				},
				expectedError: domain.ErrNotFound,
			},
			{
				name:    "internalError",
				inputID: schemas.ValidSchemaID1,
				input: domain.UpdateSchemaInput{
					Name: &newSchemaName,
				},
				expectedError: schemas.InternalError,
				getContext: func(ctx context.Context) context.Context {
					return utils.SetWithErrorToContext(ctx, true)
				},
			},
		},
	},
	{
		name: "DeleteSchema",
		executeMethod: func(ctx context.Context, s *SchemaService, inputID string, input interface{}, expectedError error) error {
			err := s.DeleteSchema(ctx, inputID)

			if expectedError != nil {
				if expectedError != err {
					return errors.New("expected an error")
				}

				return nil
			}

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: schemas.ValidSchemaID1,
				postCheck: func(s *SchemaService) error {
					_, err := s.repo.GetById(context.Background(), schemas.ValidSchemaID1)
					if err != domain.ErrNotFound {
						return errors.New("schema is not deleted")
					}

					return nil
				},
			},
			{
				name:          "notFound",
				inputID:       schemas.NotFoundSchemaID,
				expectedError: domain.ErrNotFound,
			},
			{
				name:          "internalError",
				inputID:       schemas.ValidSchemaID1,
				expectedError: schemas.InternalError,
				getContext: func(ctx context.Context) context.Context {
					return utils.SetWithErrorToContext(ctx, true)
				},
			},
		},
	},
}

func TestSchemasService(t *testing.T) {
	for _, tcGroup := range schemasTestCaseGroup {
		for _, tc := range tcGroup.testCases {
			t.Run(fmt.Sprintf("%s_%s", tcGroup.name, tc.name), func(t *testing.T) {
				schemasRepository := schemas.NewMockSchemasRepository()
				ss := NewSchemaService(schemasRepository, testConfig)
				ctx := context.Background()
				if tc.getContext != nil {
					ctx = tc.getContext(ctx)
				}
				err := tcGroup.executeMethod(ctx, ss, tc.inputID, tc.input, tc.expectedError)

				if err != nil {
					t.Error(err.Error())
					return
				}

				if tc.postCheck != nil {
					err := tc.postCheck(ss)

					if err != nil {
						t.Error(err.Error())
						return
					}
				}
			})
		}
	}
}
