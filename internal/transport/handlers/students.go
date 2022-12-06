package handlers

import (
	"errors"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/gorilla/mux"
)

type StudentResponse struct {
	Student domain.StudentRecord `json:"student"`
}
type StudentsResponse struct {
	Students []domain.StudentRecord `json:"students"`
}

// @Summary Get Student By ID
// @Description get student by id
// @Security UsersAuth
// @Tags student
// @Success 200 {object} StudentResponse
// @Param id path string true "student id"
// @Failure 404
// @Failure 500
// @Accept  json
// @Produce  json
// @Router /students/{id} [get]
func (s *Server) getStudentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		sendUnprocessableEntityError(w, errors.New("id should not be empty"))
		return
	}

	student, err := s.studentsService.GetStudentById(r.Context(), id)

	if err != nil {
		if err == domain.ErrNotFound {
			sendNotFoundError(w)
			return
		}
		sendServerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, StudentResponse{
		Student: *student,
	})
}

// @Summary List Students
// @Description retrieves all students
// @Security UsersAuth
// @Tags student
// @Success 200 {object} StudentsResponse
// @Param limit query int false "limit"
// @Param skip query int false "skip"
// @Param email query string false "email"
// @Param source query string false "source"
// @Param sort query string false "sort"
// @Failure 401
// @Failure 500
// @Accept json
// @Produce json
// @Router /students [get]
func (s *Server) listStudents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	limit, skip := getLimitSkip(params)
	sort := getSort(params)

	students, err := s.studentsService.ListStudents(r.Context(), domain.ListStudentsOptions{
		Email:  params.Get("email"),
		Source: params.Get("source"),
		Sort:   sort,
		Limit:  limit,
		Skip:   skip,
	})
	if err != nil {
		sendServerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, StudentsResponse{
		Students: students,
	})
}
