package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/abdukhashimov/student_aggregator/mocks"
)

const etalonUserJSON = `{"id":1,"name":"One","a_bool":false,"a_float":25.6,"nested":{"id":2,"name":"Two","a_bool":true,"a_float":0},"reference":{"id":3,"name":"Three","a_bool":false,"a_float":-69.333,"nested":{"id":0,"name":"","a_bool":false,"a_float":0},"reference":null}}`

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ABool     bool    `json:"a_bool"`
	AFloat    float64 `json:"a_float"`
	Nested    Nested  `json:"nested"`
	Reference *User   `json:"reference"`
}

type Nested struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	ABool  bool    `json:"a_bool"`
	AFloat float64 `json:"a_float"`
}

var etalonUser = User{
	ID:     1,
	Name:   "One",
	ABool:  false,
	AFloat: 25.6,
	Nested: Nested{
		ID:     2,
		Name:   "Two",
		ABool:  true,
		AFloat: 0,
	},
	Reference: &User{
		ID:        3,
		Name:      "Three",
		ABool:     false,
		AFloat:    -69.333,
		Nested:    Nested{},
		Reference: nil,
	},
}

type ErrorResponseTest struct {
	Name         string
	ExpectedCode int
	ExpectedBody string
	testFunction func(w http.ResponseWriter)
}

var ErrorResponseTests = []ErrorResponseTest{
	{
		"writeJSON",
		http.StatusOK,
		etalonUserJSON,
		func(w http.ResponseWriter) {
			writeJSON(w, http.StatusOK, etalonUser)
		},
	},
	{
		"sendServerError",
		http.StatusInternalServerError,
		`{"errors":"internal error"}`,
		func(w http.ResponseWriter) {
			sendServerError(w, errors.New("error 1"))
		},
	},
	{
		"sendUnprocessableEntityError",
		http.StatusUnprocessableEntity,
		`{"errors":"unprocessable entity"}`,
		func(w http.ResponseWriter) {
			sendUnprocessableEntityError(w, errors.New("error 1"))
		},
	},
	{
		"sendUnauthorizedError",
		http.StatusUnprocessableEntity,
		`{"errors":"invalid authentication credentials"}`,
		func(w http.ResponseWriter) {
			sendUnauthorizedError(w, "test@ts.ts")
		},
	},
	{
		"sendInvalidAuthTokenError",
		http.StatusUnauthorized,
		`{"errors":"invalid or missing authentication token"}`,
		sendInvalidAuthTokenError,
	},
	{
		"sendInvalidRefreshTokenError",
		http.StatusUnauthorized,
		`{"errors":"invalid or missing refresh token"}`,
		sendInvalidRefreshTokenError,
	},
	{
		"sendInvalidRefreshTokenError",
		http.StatusInternalServerError,
		`{"errors":"error 1"}`,
		func(w http.ResponseWriter) {
			writeErrorResponse(w, http.StatusInternalServerError, "error 1")
		},
	},
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// mock global logger
	mocks.MockAppLogger()

	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	// Do something here.

	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestReadJSON(t *testing.T) {
	t.Run("correct JSON", func(t *testing.T) {
		reader := bytes.NewReader([]byte(etalonUserJSON))
		user := User{}

		err := readJSON(reader, &user)
		if err != nil {
			t.Errorf("unexpecting error: %s", err.Error())
			return
		}

		if !reflect.DeepEqual(user, etalonUser) {
			t.Error("data is lost")
			return
		}
	})

	t.Run("corrupted JSON", func(t *testing.T) {
		reader := bytes.NewReader([]byte("{" + etalonUserJSON))
		user := &User{}

		err := readJSON(reader, user)
		if err == nil {
			t.Error("corrupted JSON must be not parsed")
			return
		}
	})
}

func TestWriteJSON(t *testing.T) {
	for _, test := range ErrorResponseTests {
		t.Run(test.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			test.testFunction(w)

			result := w.Result()
			if err := checkResult(result, test.ExpectedCode, test.ExpectedBody); err != nil {
				t.Error(err.Error())
				return
			}
		})
	}
}

func checkResult(result *http.Response, expectedCode int, expectedBody string) error {
	if result.StatusCode != expectedCode {
		return errors.New("status code is not correct")
	}

	body, err := io.ReadAll(result.Body)
	bodyString := string(body)
	if err != nil {
		errMsg := fmt.Sprintf("unexpecting error: %s", err.Error())
		return errors.New(errMsg)
	}
	if bodyString != expectedBody {
		return errors.New("data is lost")
	}

	contentType, ok := result.Header["Content-Type"]
	if !ok {
		return errors.New("Content-Type header is missing")
	}

	if len(contentType) == 0 || contentType[0] != "application/json" {
		return errors.New("Content-Type header must be application/json")
	}

	return nil
}
