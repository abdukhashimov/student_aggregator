package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

type M map[string]interface{}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	var jsonBytes []byte
	if data != nil {
		jB, err := json.Marshal(data)
		if err != nil {
			sendServerError(w, err)
			return
		}
		jsonBytes = jB
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if len(jsonBytes) > 0 {
		_, err := w.Write(jsonBytes)

		if err != nil {
			logger.Log.Error(err)
		}
	}
}

func sendCode(w http.ResponseWriter, code int) {
	writeJSON(w, code, nil)
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

func sendValidationError(w http.ResponseWriter, errs []string) {
	writeErrorResponse(w, http.StatusUnprocessableEntity, errs)
}

func sendInvalidAuthTokenError(w http.ResponseWriter) {
	logger.Log.Info("invalid or missing authentication token")
	w.Header().Set("WWW-Authenticate", "Token")
	msg := "invalid or missing authentication token"
	writeErrorResponse(w, http.StatusUnauthorized, msg)
}

func sendInvalidRefreshTokenError(w http.ResponseWriter) {
	msg := "invalid or missing refresh token"
	logger.Log.Info(msg)
	writeErrorResponse(w, http.StatusUnauthorized, msg)
}

func sendNotFoundError(w http.ResponseWriter) {
	writeErrorResponse(w, http.StatusNotFound, "resource not found")
}

func sendDuplicatedError(w http.ResponseWriter, field string) {
	writeErrorResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("the field [%s] is taken", field))
}

func writeErrorResponse(w http.ResponseWriter, code int, errs interface{}) {
	writeJSON(w, code, M{"errors": errs})
}
