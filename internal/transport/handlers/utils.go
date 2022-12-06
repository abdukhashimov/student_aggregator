package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

const (
	defaultLimit       = 20
	defaultSortField   = "_id,DESC"
	defaultSortType    = -1
	DescSortType       = -1
	AscSortType        = 1
	DescSortTypeString = "DESC"
	AscSortTypeString  = "ASC"
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

func getLimitSkip(params url.Values) (int, int) {
	limit := defaultLimit
	if params.Has("limit") {
		var err error
		limit, err = strconv.Atoi(params.Get("limit"))
		if err != nil {
			limit = defaultLimit
		}
	}
	skip := 0
	if params.Has("skip") {
		var err error
		skip, err = strconv.Atoi(params.Get("skip"))
		if err != nil {
			skip = 0
		}
	}
	return limit, skip
}

func getSort(params url.Values) map[string]int {
	sort := map[string]int{
		defaultSortField: defaultSortType,
	}
	if params.Has("sort") {
		parts := strings.Split(params.Get("sort"), ",")
		if len(parts) == 1 {
			delete(sort, defaultSortField)
			sort[parts[0]] = defaultSortType
		} else if len(parts) == 2 {
			if parts[1] == DescSortTypeString {
				delete(sort, defaultSortField)
				sort[parts[0]] = DescSortType
			}
			if parts[1] == AscSortTypeString {
				delete(sort, defaultSortField)
				sort[parts[0]] = AscSortType
			}
		}

	}
	return sort
}
