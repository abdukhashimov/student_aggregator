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
