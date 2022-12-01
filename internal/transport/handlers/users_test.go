package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks/services/users"
)

var LoginUserTestCases = []struct {
	name         string
	email        string
	password     string
	expectedBody string
	expectedCode int
}{
	{
		name:         "success",
		email:        users.ValidUserEmail,
		password:     users.ValidUserPassword,
		expectedBody: fmt.Sprintf(`{"tokens":{"access_token":"%s","refresh_token":"%s"}}`, users.ValidAccessToken, users.ValidRefreshToken),
		expectedCode: http.StatusOK,
	},
	{
		name:         "notFound",
		email:        users.NotFoundUserEmail,
		password:     "123456",
		expectedBody: `{"errors":"invalid authentication credentials"}`,
		expectedCode: http.StatusUnprocessableEntity,
	},
	{
		name:         "internalServerError",
		email:        users.InternalServerErrorEmail,
		password:     "123456",
		expectedBody: `{"errors":"internal error"}`,
		expectedCode: http.StatusInternalServerError,
	},
}

var CreateUserTestCases = []struct {
	name         string
	username     string
	email        string
	password     string
	expectedBody string
	expectedCode int
}{
	{
		name:         "success",
		username:     "user name",
		email:        "newemail@ts.ts",
		password:     "123456",
		expectedBody: "",
		expectedCode: http.StatusCreated,
	},
	{
		name:         "failure",
		username:     "user name",
		email:        users.InternalServerErrorEmail,
		password:     "123456",
		expectedBody: `{"message":"internal error"}`,
		expectedCode: http.StatusInternalServerError,
	},
}

var GetCurrentUserTestCases = []struct {
	name           string
	prepareRequest func(r *http.Request) *http.Request
	expectedBody   string
	expectedCode   int
}{
	{
		name: "success",
		prepareRequest: func(r *http.Request) *http.Request {
			ctx := context.WithValue(r.Context(), userKey, &users.EtalonUser)

			return r.WithContext(ctx)
		},
		expectedBody: string(users.EtalonProfileJson),
		expectedCode: http.StatusOK,
	},
	{
		name: "failure",
		prepareRequest: func(r *http.Request) *http.Request {
			ctx := context.WithValue(r.Context(), userKey, nil)

			return r.WithContext(ctx)
		},
		expectedBody: `{"errors":"internal error"}`,
		expectedCode: http.StatusInternalServerError,
	},
}

var RefreshTokenTestCases = []struct {
	name         string
	token        string
	expectedBody string
	expectedCode int
}{
	{
		name:         "success",
		token:        users.ValidRefreshToken,
		expectedBody: fmt.Sprintf(`{"tokens":{"access_token":"%s","refresh_token":"%s"}}`, users.ValidAccessToken, users.ValidRefreshToken),
		expectedCode: http.StatusOK,
	},
	{
		name:         "expired",
		token:        users.ExpiredRefreshToken,
		expectedBody: `{"errors":"invalid or missing refresh token"}`,
		expectedCode: http.StatusUnauthorized,
	},
	{
		name:         "notfound",
		token:        "just random string",
		expectedBody: `{"errors":"invalid or missing refresh token"}`,
		expectedCode: http.StatusUnauthorized,
	},
	{
		name:         "failure",
		token:        users.InternalServerErrorRefreshToken,
		expectedBody: `{"errors":"internal error"}`,
		expectedCode: http.StatusInternalServerError,
	},
}

func TestLoginUser(t *testing.T) {
	for _, tc := range LoginUserTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := users.NewMockUsersService()
			server := &Server{
				userService: mockUsersService,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			input := domain.SignInUserInput{
				Email:    tc.email,
				Password: tc.password,
			}
			r = inputToContext(r, &input)

			// Test loginUser method
			server.loginUser(w, r)

			result := w.Result()
			body, _ := io.ReadAll(result.Body)
			bodyString := string(body)

			if bodyString != tc.expectedBody {
				t.Error("unexpected response")
				return
			}

			if result.StatusCode != tc.expectedCode {
				t.Error("unexpected status code")
				return
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	for _, tc := range CreateUserTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := users.NewMockUsersService()
			server := &Server{
				userService: mockUsersService,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			input := domain.SignUpUserInput{
				Username: tc.username,
				Email:    tc.email,
				Password: tc.password,
			}
			r = inputToContext(r, &input)

			// Test createUser method
			server.createUser(w, r)

			result := w.Result()
			body, _ := io.ReadAll(result.Body)
			bodyString := string(body)

			if bodyString != tc.expectedBody {
				t.Error("unexpected response")
				return
			}

			if result.StatusCode != tc.expectedCode {
				t.Error("unexpected status code")
				return
			}

			if tc.expectedCode == http.StatusCreated {
				_, err := mockUsersService.SignIn(context.Background(), domain.SignInUserInput{
					Email:    tc.email,
					Password: tc.password,
				})
				if err != nil {
					t.Errorf("unexpecting error: %s", err.Error())
					return
				}
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	for _, tc := range GetCurrentUserTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := users.NewMockUsersService()
			server := &Server{
				userService: mockUsersService,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = tc.prepareRequest(r)

			// Test getCurrentUser method
			server.getCurrentUser(w, r)

			result := w.Result()
			body, _ := io.ReadAll(result.Body)
			bodyString := string(body)

			if bodyString != tc.expectedBody {
				t.Error("unexpected response")
				return
			}

			if result.StatusCode != tc.expectedCode {
				t.Error("unexpected status code")
				return
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	for _, tc := range RefreshTokenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := users.NewMockUsersService()
			server := &Server{
				userService: mockUsersService,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			input := domain.TokenInput{
				Token: tc.token,
			}
			r = inputToContext(r, &input)

			// Test refreshToken method
			server.refreshToken(w, r)

			result := w.Result()
			body, _ := io.ReadAll(result.Body)
			bodyString := string(body)

			if bodyString != tc.expectedBody {
				t.Error("unexpected response")
				return
			}

			if result.StatusCode != tc.expectedCode {
				t.Error("unexpected status code")
				return
			}
		})
	}
}
