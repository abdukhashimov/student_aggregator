package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

type M map[string]interface{}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		sendServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)

	if err != nil {
		logger.Log.Error(err)
	}
}

func sendCode(w http.ResponseWriter, code int) {
	writeJSON(w, code, M{})
}

func sendServerError(w http.ResponseWriter, err error) {
	logger.Log.Error(err)
	writeErrorResponse(w, http.StatusInternalServerError, "internal error")
}

func sendUnprocessableEntityError(w http.ResponseWriter, err error) {
	logger.Log.Error(err)
	writeErrorResponse(w, http.StatusUnprocessableEntity, "unprocessable entity")
}

func sendUnauthorizedError(w http.ResponseWriter, email string) {
	logger.Log.Debugf("invalid authentication credentials. email: %s", email)
	writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid authentication credentials")
}

func writeErrorResponse(w http.ResponseWriter, code int, errs interface{}) {
	writeJSON(w, code, M{"errors": errs})
}
