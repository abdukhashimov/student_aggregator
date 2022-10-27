package handlers

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

func TestUserContext(t *testing.T) {
	t.Run("user context", func(t *testing.T) {
		request := &http.Request{}
		originalUser := domain.User{
			ID:       "Test ID",
			Username: "Test Username",
			Email:    "test@ts.ts",
		}
		user1 := originalUser

		requestWithUser := setContextUser(request, &user1)

		user2, err := userFromContext(requestWithUser.Context())

		if err != nil {
			t.Errorf("unexpecting error: %s", err.Error())
			return
		}

		if !reflect.DeepEqual(originalUser, *user2) {
			t.Error("user data is lost")
		}

		if &user1 != user2 {
			t.Error("data from context is not the same")
		}
	})
}
