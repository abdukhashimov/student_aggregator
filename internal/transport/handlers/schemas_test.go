package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks/services"
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
	requestBody    any
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
			return s.listSchemas()
		},
		testCases: []SchemaTestCase{
			{
				name:         "success",
				expectedBody: `{"schemas":[{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]},{"id":"2","name":"WAC","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"name"},{"col":"B","name":"surname"},{"col":"C","name":"email"}]}]}`,
				expectedCode: http.StatusOK,
			},
			{
				name: "internalError",
				prepareRequest: func(r *http.Request) *http.Request {
					return services.SetWithError(r, true)
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
			return s.createSchema()
		},
		testCases: []SchemaTestCase{
			{
				name: "success",
				requestBody: domain.NewSchemaInput{
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
				expectedBody: `{"schema":{"id":"3","name":"New Schema","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"name"},{"col":"B","name":"surname"},{"col":"C","name":"email"}]}}`,
				expectedCode: http.StatusCreated,
			},
			{
				name: "duplicate",
				requestBody: domain.NewSchemaInput{
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
				requestBody: domain.NewSchemaInput{
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
					return services.SetWithError(r, true)
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
			return s.getSchemaById()
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
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]}}`,
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
					return services.SetWithError(r, true)
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
			return s.updateSchema()
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
				requestBody: domain.UpdateSchemaInput{
					Name:       &updateSchemaName,
					Version:    &updateSchemaVersion,
					SchemaType: &updateSchemaSchemaType,
					Headers:    &updateSchemaHeaders,
					Fields:     &updateSchemaFields,
				},
				expectedBody: `{"schema":{"id":"1","name":"updateSchemaName","version":"updateSchemaVersion","schema_type":"updateSchemaSchemaType","headers":false,"fields":[{"col":"D","name":"name"},{"col":"E","name":"surname"},{"col":"F","name":"em"}]}}`,
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
				requestBody: domain.UpdateSchemaInput{
					Name: &updateSchemaName,
				},
				expectedBody: `{"schema":{"id":"1","name":"updateSchemaName","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]}}`,
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
				requestBody: domain.UpdateSchemaInput{
					Version: &updateSchemaVersion,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"updateSchemaVersion","schema_type":"coords","headers":true,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]}}`,
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
				requestBody: domain.UpdateSchemaInput{
					SchemaType: &updateSchemaSchemaType,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"updateSchemaSchemaType","headers":true,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]}}`,
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
				requestBody: domain.UpdateSchemaInput{
					Headers: &updateSchemaHeaders,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":false,"fields":[{"col":"A","name":"first_name"},{"col":"B","name":"last_name"},{"col":"C","name":"email"}]}}`,
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
				requestBody: domain.UpdateSchemaInput{
					Fields: &updateSchemaFields,
				},
				expectedBody: `{"schema":{"id":"1","name":"RSS","version":"1.0.0","schema_type":"coords","headers":true,"fields":[{"col":"D","name":"name"},{"col":"E","name":"surname"},{"col":"F","name":"em"}]}}`,
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
				expectedBody: `{"errors":"unprocessable entity"}`,
				expectedCode: http.StatusUnprocessableEntity,
			},
			{
				name: "notFound",
				prepareRequest: func(r *http.Request) *http.Request {
					r = mux.SetURLVars(r, map[string]string{
						"id": schemas.NotFoundSchemaID,
					})
					return r
				},
				requestBody: domain.UpdateSchemaInput{
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
					return services.SetWithError(r, true)
				},
				requestBody: domain.UpdateSchemaInput{
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
			return s.deleteSchema()
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
					return services.SetWithError(r, true)
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

				var requestDody io.Reader
				if tc.requestBody != nil {
					data, err := json.Marshal(tc.requestBody)
					if err != nil {
						t.Error("unexpected response")
						return
					}

					requestDody = bytes.NewReader(data)
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(tcGroup.requestMethod, "/", requestDody)
				if tc.prepareRequest != nil {
					r = tc.prepareRequest(r)
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
