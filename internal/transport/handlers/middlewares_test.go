package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks/services/users"
)

type TestCase struct {
	Name           string
	Token          string
	ExpectedCode   int
	ExpectedBody   string
	ExpectedUserId string
}

var TestCases = []TestCase{
	{
		"Valid token",
		users.ValidAccessToken,
		http.StatusOK,
		"requestHandler is called",
		users.ValidUserId,
	},
	{
		"Invalid token",
		users.InvalidAccessToken,
		http.StatusUnauthorized,
		`{"errors":"invalid or missing authentication token"}`,
		"",
	},
	{
		"Empty token",
		"",
		http.StatusUnauthorized,
		`{"errors":"invalid or missing authentication token"}`,
		"",
	},
}

func TestAuthenticateMiddleware(t *testing.T) {
	mockUsersService := users.NewMockUsersService()
	server := &Server{
		userService: mockUsersService,
	}

	for _, test := range TestCases {
		t.Run(test.Name, func(t *testing.T) {
			var (
				userFromCtx *domain.User
				handlerErr  error
			)
			var requestHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
				_, handlerErr = w.Write([]byte("requestHandler is called"))
				if handlerErr == nil {
					userFromCtx, handlerErr = userFromContext(r.Context())
				}
			}

			// Test authenticate middleware function
			httpHandler := server.authenticate(requestHandler)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.Token != "" {
				r.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", test.Token)}
			}

			httpHandler.ServeHTTP(w, r)

			result := w.Result()

			if handlerErr != nil {
				t.Errorf("unexpecting error: %s", handlerErr.Error())
				return
			}

			if result.StatusCode != test.ExpectedCode {
				t.Error("status code is not correct")
				return
			}

			body, err := io.ReadAll(result.Body)
			bodyString := string(body)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}
			if bodyString != test.ExpectedBody {
				t.Error("data is lost")
				return
			}

			if test.ExpectedUserId != "" && userFromCtx.ID != test.ExpectedUserId {
				t.Error("user is not stored to context")
				return
			}
		})
	}
}
