package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/abdukhashimov/student_aggregator/mocks/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks/services/schemas"
	"github.com/gorilla/mux"
)

var (
	updateSchemaName       = "updateSchemaName"
	updateSchemaVersion    = "updateSchemaVersion"
	updateSchemaSchemaType = "updateSchemaSchemaType"
	updateSchemaHeaders    = false
	updateSchemaFields     = []domain.FieldSchema{
		{Name: "name", Col: "D"},
		{Name: "surname", Col: "E"},
		{Name: "em", Col: "F"},
	}
)

type SchemaTestCaseGroup struct {
	name          string
	requestMethod string
	getHandler    func(s *Server) http.HandlerFunc
	testCases     []SchemaTestCase
}

type SchemaTestCase struct {
	name           string
	requestInput   any
	prepareRequest func(r *http.Request) *http.Request
	expectedBody   string
	expectedCode   int
	postCheck      func(s *Server) error
}

var SchemasTestCases = []SchemaTestCaseGroup{
	{
		name:          "listSchemas",
		requestMethod: http.MethodGet,
		getHandler: func(s *Server) http.HandlerFunc {
			return s.listSchemas
		},
		testCases: []SchemaTestCase{
			{
				name:         "success",
				expectedBody: `{"schemas":[{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]},{"id":"2","name":"WAC","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"surname","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}]}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "internalError",
				prepareRequest: func(r *http.Request) *http.Request {
					return utils.SetWithErrorToRequest(r, true)
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
		},
	},
	{
		name:          "createSchema",
		requestMethod: http.MethodPost,
		getHandler: func(s *Server) http.HandlerFunc {
			return s.createSchema
		},
		testCases: []SchemaTestCase{
			{
				name: "success",
				requestInput: &domain.NewSchemaInput{
					Name:       "New Schema",
					Version:    "1.0.0",
					SchemaType: "coords",
					Headers:    true,
					Fields: []domain.FieldSchema{
						{Name: "name", Col: "A"},
						{Name: "surname", Col: "B"},
						{Name: "email", Col: "C"},
					},
				},
				expectedBody: `{"schema":{"id":"3","name":"New Schema","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"surname","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusCreated,
			},
			{
				name: "duplicate",
				requestInput: &domain.NewSchemaInput{
					Name:       "RSS",
					Version:    "1.0.0",
					SchemaType: "coords",
					Headers:    true,
					Fields: []domain.FieldSchema{
						{Name: "name", Col: "A"},
						{Name: "surname", Col: "B"},
						{Name: "email", Col: "C"},
					},
				},
				expectedBody: `{"errors":"the field [name] is taken"}`,
				expectedCode: http.StatusUnprocessableEntity,
			},
			{
				name: "internalError",
				requestInput: &domain.NewSchemaInput{
					Name:       "New Schema",
					Version:    "1.0.0",
					SchemaType: "coords",
					Headers:    true,
					Fields: []domain.FieldSchema{
						{Name: "name", Col: "A"},
						{Name: "surname", Col: "B"},
						{Name: "email", Col: "C"},
					},
				},
				prepareRequest: func(r *http.Request) *http.Request {
					return utils.SetWithErrorToRequest(r, true)
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
		},
	},
	{
		name:          "getSchemaById",
		requestMethod: http.MethodGet,
		getHandler: func(s *Server) http.HandlerFunc {
			return s.getSchemaById
		},
		testCases: []SchemaTestCase{
			{
				name: "success",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "notFound",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.NotFoundSchemaID,
					})
					return r
				},
				expectedBody: `{"errors":"resource not found"}`,
				expectedCode: http.StatusNotFound,
			},
			{
				name: "internalError",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.NotFoundSchemaID,
					})
					return utils.SetWithErrorToRequest(r, true)
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
		},
	},
	{
		name:          "updateSchema",
		requestMethod: http.MethodPost,
		getHandler: func(s *Server) http.HandlerFunc {
			return s.updateSchema
		},
		testCases: []SchemaTestCase{
			{
				name: "successAllFields",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Name:       &updateSchemaName,
					Version:    &updateSchemaVersion,
					SchemaType: &updateSchemaSchemaType,
					Headers:    &updateSchemaHeaders,
					Fields:     &updateSchemaFields,
				},
				expectedBody: `{"schema":{"id":"1","name":"updateSchemaName","version":"updateSchemaVersion","schema_type":"updateSchemaSchemaType","headers":false,"fields":[{"col":"D","name":"name","is_multiple":false,"is_map":false,"map_start":false},{"col":"E","name":"surname","is_multiple":false,"is_map":false,"map_start":false},{"col":"F","name":"em","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "successName",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Name: &updateSchemaName,
				},
				expectedBody: `{"schema":{"id":"1","name":"updateSchemaName","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "successVersion",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Version: &updateSchemaVersion,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"updateSchemaVersion","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "successSchemaType",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					SchemaType: &updateSchemaSchemaType,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"updateSchemaSchemaType","headers":true,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "successHeaders",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Headers: &updateSchemaHeaders,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":false,"fields":[{"col":"A","name":"first_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"B","name":"last_name","is_multiple":false,"is_map":false,"map_start":false},{"col":"C","name":"email","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "successFields",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Fields: &updateSchemaFields,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"D","name":"name","is_multiple":false,"is_map":false,"map_start":false},{"col":"E","name":"surname","is_multiple":false,"is_map":false,"map_start":false},{"col":"F","name":"em","is_multiple":false,"is_map":false,"map_start":false}]}}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "emptyAllFields",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
			{
				name: "notFound",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.NotFoundSchemaID,
					})
					return r
				},
				requestInput: &domain.UpdateSchemaInput{
					Name:       &updateSchemaName,
					Version:    &updateSchemaVersion,
					SchemaType: &updateSchemaSchemaType,
					Headers:    &updateSchemaHeaders,
					Fields:     &updateSchemaFields,
				},
				expectedBody: `{"errors":"resource not found"}`,
				expectedCode: http.StatusNotFound,
			},
			{
				name: "internalError",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return utils.SetWithErrorToRequest(r, true)
				},
				requestInput: &domain.UpdateSchemaInput{
					Name:       &updateSchemaName,
					Version:    &updateSchemaVersion,
					SchemaType: &updateSchemaSchemaType,
					Headers:    &updateSchemaHeaders,
					Fields:     &updateSchemaFields,
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
		},
	},
	{
		name:          "deleteSchema",
		requestMethod: http.MethodDelete,
		getHandler: func(s *Server) http.HandlerFunc {
			return s.deleteSchema
		},
		testCases: []SchemaTestCase{
			{
				name: "success",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return r
				},
				expectedBody: "",
				expectedCode: http.StatusOK,
				postCheck: func(s *Server) error {
					_, err := s.schemasService.GetSchemaById(context.Background(), schemas.ValidSchemaID1)
					if err != domain.ErrNotFound {
						return errors.New("schema is not deleted")
					}
					return nil
				},
			},
			{
				name: "notFound",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.NotFoundSchemaID,
					})
					return r
				},
				expectedBody: `{"errors":"resource not found"}`,
				expectedCode: http.StatusNotFound,
			},
			{
				name: "internalError",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.ValidSchemaID1,
					})
					return utils.SetWithErrorToRequest(r, true)
				},
				expectedBody: `{"errors":"internal error"}`,
				expectedCode: http.StatusInternalServerError,
			},
		},
	},
}

func TestSchemas(t *testing.T) {
	for _, tcGroup := range SchemasTestCases {
		for _, tc := range tcGroup.testCases {
			t.Run(fmt.Sprintf("%s_%s", tcGroup.name, tc.name), func(t *testing.T) {
				mockSchemasService := schemas.NewMockSchemasService()
				server := &Server{
					schemasService: mockSchemasService,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(tcGroup.requestMethod, "/", nil)
				if tc.prepareRequest != nil {
					r = tc.prepareRequest(r)
				}
				if tc.requestInput != nil {
					r = inputToContext(r, tc.requestInput)
				}

				tcGroup.getHandler(server).ServeHTTP(w, r)

				result := w.Result()
				responseBody, _ := io.ReadAll(result.Body)
				responseBodyString := string(responseBody)

				if responseBodyString != tc.expectedBody {
					t.Error("unexpected response")
					return
				}

				if result.StatusCode != tc.expectedCode {
					t.Error("unexpected status code")
					return
				}
				if tc.postCheck != nil {
					err := tc.postCheck(server)
					if err != nil {
						t.Error(err.Error())
						return
					}
				}
			})
		}
	}
}
