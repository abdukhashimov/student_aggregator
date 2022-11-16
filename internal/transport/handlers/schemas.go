package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type SchemaResponse struct {
	Schema domain.Schema `json:"schema"`
}
type SchemasResponse struct {
	Schemas []domain.Schema `json:"schemas"`
}

// @Summary List Schemas
// @Description retrieves all schemas
// @Security UsersAuth
// @Tags schema
// @Success 200 {object} SchemasResponse
// @Failure 401
// @Failure 500
// @Accept json
// @Produce json
// @Router /schemas [get]
func (s *Server) listSchemas() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schemas, err := s.schemasService.ListSchemas(r.Context())
		if err != nil {
			sendServerError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, SchemasResponse{
			Schemas: schemas,
		})
	}
}

// @Summary Create Schema
// @Description schema creation process
// @Tags schema
// @Param request body domain.NewSchemaInput true "query params"
// @Success 201
// @Failure 422
// @Failure 500
// @Accept json
// @Produce json
// @Router /schemas [post]
func (s *Server) createSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := domain.NewSchemaInput{}

		err := readJSON(r.Body, &input)
		if err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		schema, err := s.schemasService.NewSchema(r.Context(), input)
		if err != nil {
			if err == domain.DuplicationError {
				sendDuplicatedError(w, "name")
				return
			}
			sendServerError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, SchemaResponse{
			Schema: *schema,
		})
	}
}

// @Summary Get Schema By ID
// @Description get course by id
// @Tags schema
// @Success 200 {object} SchemaResponse
// @Param id path string true "schema id"
// @Failure 404
// @Failure 500
// @Accept  json
// @Produce  json
// @Router /schemas/{id} [get]
func (s *Server) getSchemaById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			sendUnprocessableEntityError(w, errors.New("id should not be empty"))
			return
		}

		schema, err := s.schemasService.GetSchemaById(r.Context(), id)

		if err != nil {
			if err == domain.ErrNotFound {
				sendNotFoundError(w)
				return
			}
			sendServerError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, SchemaResponse{
			Schema: *schema,
		})
	}
}

// @Summary Update Schema By ID
// @Description update schema by id
// @Tags schema
// @Param input body domain.UpdateSchemaInput true "update info"
// @Success 200 {object} SchemaResponse
// @Failure 404
// @Failure 422
// @Failure 500
// @Accept  json
// @Produce  json
// @Router /schemas/{id} [patch]
func (s *Server) updateSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if id == "" {
			sendUnprocessableEntityError(w, errors.New("id should not be empty"))
			return
		}

		input := domain.UpdateSchemaInput{}
		err := readJSON(r.Body, &input)
		if err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		schema, err := s.schemasService.UpdateSchema(r.Context(), id, input)
		if err != nil {
			if err == domain.DuplicationError {
				sendDuplicatedError(w, "name")
				return
			}
			if err == domain.ErrNotFound {
				sendNotFoundError(w)
				return
			}
			sendServerError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, SchemaResponse{
			Schema: *schema,
		})
	}
}

// @Summary Delete Schema
// @Description delete schema
// @Tags schema
// @Param id path string true "schema id"
// @Success 200
// @Failure 404
// @Failure 500
// @Accept  json
// @Produce  json
// @Router /schemas/{id} [delete]
func (s *Server) deleteSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if id == "" {
			sendUnprocessableEntityError(w, errors.New("id should not be empty"))
			return
		}

		err := s.schemasService.DeleteSchema(r.Context(), id)
		if err != nil {
			if err == domain.ErrNotFound {
				sendNotFoundError(w)
				return
			}
			sendServerError(w, err)
			return
		}

		sendCode(w, http.StatusOK)
	}
}
