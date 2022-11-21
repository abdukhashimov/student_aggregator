package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const (
	ValidMongoId  = "507f1f77bcf86cd000000001"
	ValidMongoId2 = "507f1f77bcf86cd000000002"
)

var generalError = errors.New("general")

var EtalonSchema = domain.Schema{
	ID:         ValidMongoId,
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
	ID:         ValidMongoId2,
	Name:       "WAA",
	Slug:       "waa",
	Version:    "1.0.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []domain.FieldSchema{
		{Name: "name", Col: "A"},
		{Name: "surname", Col: "B"},
		{Name: "email", Col: "C"},
	},
}

var newSchemaInput = domain.Schema{
	Name:       "New Schema Name",
	Slug:       "new-schema-name",
	Version:    "1.0.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []domain.FieldSchema{
		{Name: "name", Col: "A"},
		{Name: "surname", Col: "B"},
		{Name: "email", Col: "C"},
	},
}

var (
	updateSchemaName       = "updateSchemaName"
	updateSchemaSlug       = "updateschemaname"
	updateSchemaVersion    = "updateSchemaVersion"
	updateSchemaSchemaType = "updateSchemaSchemaType"
	updateSchemaHeaders    = false
	updateSchemaFields     = []domain.FieldSchema{
		{Name: "name", Col: "D"},
		{Name: "surname", Col: "E"},
		{Name: "em", Col: "F"},
	}
)

type SchemasTestCaseGroup struct {
	name          string
	executeMethod func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error
	testCases     []SchemasTestCase
}

type SchemasTestCase struct {
	name          string
	inputID       string
	input         interface{}
	getMongoRes   func() ([]bson.D, error)
	expectedError error
}

var schemasTestCaseGroup = []SchemasTestCaseGroup{
	{
		name: "Create",
		executeMethod: func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error {
			in, ok := tc.input.(domain.Schema)
			if !ok {
				return errors.New("invalid input type")
			}
			newSchemaId, err := repo.Create(ctx, in)

			err, skip := checkError(tc, err)

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if skip {
				return nil
			}

			if newSchemaId == "" {
				return errors.New("invalid schema ID")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:  "success",
				input: newSchemaInput,
				getMongoRes: func() ([]bson.D, error) {
					success := mtest.CreateSuccessResponse()
					return []bson.D{success}, nil
				},
			},
			{
				name:  "failure",
				input: newSchemaInput,
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					})
					return []bson.D{err}, nil
				},
				expectedError: generalError,
			},
			{
				name:  "duplicate",
				input: newSchemaInput,
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    11000,
						Message: "duplicate error",
					})
					return []bson.D{err}, nil
				},
				expectedError: domain.DuplicationError,
			},
		},
	},
	{
		name: "GetById",
		executeMethod: func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error {
			schema, err := repo.GetById(ctx, tc.inputID)

			err, skip := checkError(tc, err)

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if skip {
				return nil
			}

			if !reflect.DeepEqual(*schema, EtalonSchema) {
				return errors.New("invalid result")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
			},
			{
				name:    "notFound",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					return getNotFoundSchemaMongoRes()
				},
				expectedError: domain.ErrNotFound,
			},
			{
				name:    "invalidId",
				inputID: "1",
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
				expectedError: generalError,
			},
			{
				name:    "failure",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					})
					return []bson.D{err}, nil
				},
				expectedError: generalError,
			},
		},
	},
	{
		name: "Update",
		executeMethod: func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error {
			in, ok := tc.input.(domain.UpdateSchemaInput)
			if !ok {
				return errors.New("invalid input type")
			}
			err := repo.Update(ctx, tc.inputID, in)

			err, _ = checkError(tc, err)

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: ValidMongoId,
				input: domain.UpdateSchemaInput{
					Name:       &updateSchemaName,
					Slug:       &updateSchemaSlug,
					Version:    &updateSchemaVersion,
					SchemaType: &updateSchemaSchemaType,
					Headers:    &updateSchemaHeaders,
					Fields:     &updateSchemaFields,
				},
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
			},
			{
				name:    "successPartial",
				inputID: ValidMongoId,
				input: domain.UpdateSchemaInput{
					Fields: &updateSchemaFields,
				},
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
			},
			{
				name:    "invalidId",
				inputID: "1",
				input: domain.UpdateSchemaInput{
					Fields: &updateSchemaFields,
				},
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
				expectedError: generalError,
			},
			{
				name:    "failure",
				inputID: ValidMongoId,
				input: domain.UpdateSchemaInput{
					Fields: &updateSchemaFields,
				},
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					})
					return []bson.D{err}, nil
				},
				expectedError: generalError,
			},
		},
	},
	{
		name: "FindAll",
		executeMethod: func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error {
			schemasList, err := repo.FindAll(ctx)

			err, skip := checkError(tc, err)
			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			if skip {
				return nil
			}

			if !reflect.DeepEqual(schemasList, []domain.Schema{EtalonSchema, EtalonSchema2}) {
				return errors.New("invalid result")
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name: "success",
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema, EtalonSchema2})
				},
			},
			{
				name: "failure",
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					})
					return []bson.D{err}, nil
				},
				expectedError: generalError,
			},
		},
	},
	{
		name: "Delete",
		executeMethod: func(ctx context.Context, repo *SchemaRepo, tc *SchemasTestCase) error {
			err := repo.Delete(ctx, tc.inputID)

			err, _ = checkError(tc, err)

			if err != nil {
				return fmt.Errorf("unexpecting error: %s", err.Error())
			}

			return nil
		},
		testCases: []SchemasTestCase{
			{
				name:    "success",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					return []bson.D{{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}}}, nil
				},
			},
			{
				name:    "notFound",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					return []bson.D{{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}}}, nil
				},
				expectedError: domain.ErrNotFound,
			},
			{
				name:    "invalidId",
				inputID: "1",
				getMongoRes: func() ([]bson.D, error) {
					return getSuccessSchemaMongoRes([]domain.Schema{EtalonSchema})
				},
				expectedError: generalError,
			},
			{
				name:    "failure",
				inputID: ValidMongoId,
				getMongoRes: func() ([]bson.D, error) {
					err := mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					})

					return []bson.D{err}, nil
				},
				expectedError: generalError,
			},
		},
	},
}

func TestSchemasRepo(t *testing.T) {
	mt := getMockTest(t)
	defer mt.Close()

	for _, tcGroup := range schemasTestCaseGroup {
		for _, tc := range tcGroup.testCases {
			mt.Run(fmt.Sprintf("%s_%s", tcGroup.name, tc.name), func(mt *mtest.T) {
				mocks, err := tc.getMongoRes()
				if err != nil {
					t.Errorf("unexpecting error: %s", err.Error())
				}
				mt.AddMockResponses(mocks...)
				schemaRepo := NewSchemaRepo(mt.DB)
				err = tcGroup.executeMethod(context.Background(), schemaRepo, &tc)

				if err != nil {
					t.Error(err.Error())
					return
				}
			})
		}
	}
}

func checkError(tc *SchemasTestCase, err error) (error, bool) {
	skip := true
	if tc.expectedError != nil {
		if tc.expectedError != generalError {
			if tc.expectedError != err {
				return errors.New("errors do not match"), skip
			}
			return nil, skip
		}
		if err == nil {
			return errors.New("expected an error"), skip
		}
		return nil, skip
	}
	skip = false

	return err, skip
}

func getSuccessSchemaMongoRes(expectedSchemas []domain.Schema) ([]bson.D, error) {
	var result []bson.D
	for index, schema := range expectedSchemas {
		batch := mtest.NextBatch
		if index == 0 {
			batch = mtest.FirstBatch
		}
		bsonD, err := toBson(schema)
		if err != nil {
			return nil, err
		}
		result = append(result, mtest.CreateCursorResponse(
			1,
			fmt.Sprintf("%s.%s", testDbName, schemasCollection),
			batch,
			bsonD),
		)
	}

	result = append(result, mtest.CreateCursorResponse(
		0,
		fmt.Sprintf("%s.%s", testDbName, schemasCollection),
		mtest.NextBatch),
	)

	return result, nil
}

func getNotFoundSchemaMongoRes() ([]bson.D, error) {
	res := mtest.CreateCursorResponse(
		1,
		fmt.Sprintf("%s.%s", testDbName, schemasCollection),
		mtest.FirstBatch)
	end := mtest.CreateCursorResponse(
		0,
		fmt.Sprintf("%s.%s", testDbName, schemasCollection),
		mtest.NextBatch)

	return []bson.D{res, end}, nil
}
