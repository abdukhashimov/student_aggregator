package handlers

import (
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

// @Summary Parse File
// @Description File Parsing process
// @Security UsersAuth
// @Tags parser
// @Param request body domain.ParseFileInput true "query params"
// @Success 200
// @Failure 422
// @Failure 500
// @Accept json
// @Produce json
// @Router /aggregator/parse [post]
func (s *Server) parseFile(w http.ResponseWriter, r *http.Request) {
	input, err := inputFromContext[domain.ParseFileInput](r.Context())
	if err != nil {
		sendServerError(w, err)
		return
	}

	err = s.aggregatorService.ParseFile(r.Context(), input.FileName, input.SchemaID)
	if err != nil {
		sendServerError(w, err)
		return
	}

	// TODO: add appropriate response
	writeJSON(w, http.StatusOK, M{"message": "ok"})
}
