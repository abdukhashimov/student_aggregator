package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/mocks"
	"github.com/abdukhashimov/student_aggregator/mocks/repository"
)

const (
	testPassword = "123456"
)

var testConfig = &config.Config{
	Project: config.ProjectConfig{
		Salt:      "123456",
		JwtSecret: "654321",
	},
	Http: config.HttpConfig{
		AccessTokenTTLMinutes: 5,
		RefreshTokenTTLHours:  1,
	},
}

var SignUpTestCases = []struct {
	name          string
	testEmail     string
	testPassword  string
	expectedError bool
	postCheck     func(us *UsersService, id string) error
}{
	{name: "success", testEmail: "newvalidemail@ts.ts", testPassword: testPassword, expectedError: false, postCheck: func(us *UsersService, id string) error {
		user, err := us.UserById(context.Background(), id)
		if err != nil {
			return err
		}
		hashedPassword, _ := us.hasher.Hash(testPassword)
		if user.Password != hashedPassword {
			return errors.New("password is not hashed")
		}

		return nil
	}},
	{name: "duplicate", testEmail: repository.TakenUserEmail, testPassword: testPassword, expectedError: true},
	{name: "failure", testEmail: repository.ErrorToCreateUserEmail, testPassword: testPassword, expectedError: true},
}

var SignInTestCases = []struct {
	name           string
	usernameSignUp string
	emailSignUp    string
	passwordSignUp string
	emailSignIn    string
	passwordSignIn string
	expectedError  bool
}{
	{
		name:           "success",
		usernameSignUp: "test 1",
		emailSignUp:    "validemail@ts.ts",
		passwordSignUp: "123456",
		emailSignIn:    "validemail@ts.ts",
		passwordSignIn: "123456",
		expectedError:  false,
	},
	{
		name:           "wrong_password",
		usernameSignUp: "test 1",
		emailSignUp:    "validemail@ts.ts",
		passwordSignUp: "123456",
		emailSignIn:    "validemail@ts.ts",
		passwordSignIn: "654321",
		expectedError:  true,
	},
	{
		name:           "wrong_email",
		usernameSignUp: "test 1",
		emailSignUp:    "validemail@ts.ts",
		passwordSignUp: "123456",
		emailSignIn:    "wrongemail@ts.ts",
		passwordSignIn: "123456",
		expectedError:  true,
	},
	{
		name:           "empty",
		usernameSignUp: "test 1",
		emailSignUp:    "validemail@ts.ts",
		passwordSignUp: "123456",
		emailSignIn:    "",
		passwordSignIn: "",
		expectedError:  true,
	},
}

var UserByAccessTokenTestCases = []struct {
	name          string
	userId        string
	ttl           time.Duration
	expectedError bool
	corrupted     bool
}{
	{name: "success", userId: repository.ValidMongoId, ttl: 1 * time.Minute, expectedError: false},
	{name: "corrupted", userId: repository.ValidMongoId, ttl: 1 * time.Minute, expectedError: true, corrupted: true},
	{name: "expired", userId: repository.ValidMongoId, ttl: -1 * time.Minute, expectedError: true},
	{name: "failure", userId: "1", ttl: 1 * time.Minute, expectedError: true},
}

var UserByRefreshTokenTestCases = []struct {
	name          string
	token         string
	expectedId    string
	expectedError bool
}{
	{name: "success", token: repository.ValidRefreshToken, expectedId: repository.ValidMongoId, expectedError: false},
	{name: "expired", token: repository.ExpiredRefreshToken, expectedError: true},
	{name: "notfound", token: repository.MissedRefreshToken, expectedError: true},
}

var GeneralTestCases = []struct {
	name          string
	userId        string
	expectedError bool
}{
	{name: "success", userId: repository.ValidMongoId, expectedError: false},
	{name: "failure", userId: repository.InvalidMongoId, expectedError: true},
	{name: "notfound", userId: repository.NotFoundMongoId, expectedError: true},
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

func TestSignUp(t *testing.T) {
	for _, tc := range SignUpTestCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			id, err := us.SignUp(context.Background(), domain.SignUpUserInput{
				Username: "test user",
				Email:    tc.testEmail,
				Password: tc.testPassword,
			})

			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if tc.postCheck != nil {
				err := tc.postCheck(us, id)
				if err != nil {
					t.Errorf("unexpecting error: %s", err.Error())
					return
				}
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	for _, tc := range SignInTestCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			newUserId, err := us.SignUp(context.Background(), domain.SignUpUserInput{
				Username: tc.usernameSignUp,
				Email:    tc.emailSignUp,
				Password: tc.passwordSignUp,
			})
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			signInId, err := us.SignIn(context.Background(), domain.SignInUserInput{
				Email:    tc.emailSignIn,
				Password: tc.passwordSignIn,
			})
			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if newUserId != signInId {
				t.Error("user is wrong")
				return
			}
		})
	}
}

func TestUserByAccessToken(t *testing.T) {
	for _, tc := range UserByAccessTokenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)
			accessToken, _ := us.tokeManager.NewJWT(tc.userId, tc.ttl)

			if tc.corrupted {
				accessToken += "$"
			}

			user, err := us.UserByAccessToken(context.Background(), accessToken)
			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if user.ID != tc.userId {
				t.Error("user is wrong")
				return
			}
		})
	}
}

func TestUserByRefreshToken(t *testing.T) {
	for _, tc := range UserByRefreshTokenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			user, err := us.UserByRefreshToken(context.Background(), tc.token)
			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if user.ID != tc.expectedId {
				t.Error("user is wrong")
				return
			}
		})
	}
}

func TestUserById(t *testing.T) {
	for _, tc := range GeneralTestCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			user, err := us.UserById(context.Background(), tc.userId)
			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if user.ID != tc.userId {
				t.Error("user is wrong")
				return
			}
		})
	}
}

func TestGenerateUserTokens(t *testing.T) {
	for _, tc := range GeneralTestCases {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Now()
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			tokens, err := us.GenerateUserTokens(context.Background(), tc.userId)
			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			user, err := us.UserById(context.Background(), tc.userId)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if user.RefreshToken.Token != tokens.RefreshToken {
				t.Error("token is not stored")
				return
			}

			delta := 1 * time.Minute
			expectedTokenExpiresAt := now.Add(time.Hour * time.Duration(us.cfg.Http.RefreshTokenTTLHours))
			if user.RefreshToken.ExpiresAt.Add(delta).Before(expectedTokenExpiresAt) || user.RefreshToken.ExpiresAt.Add(-delta).After(expectedTokenExpiresAt) {
				t.Error("invalid expiration date")
				return
			}

			idFromToken, err := us.tokeManager.Parse(tokens.AccessToken)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if tc.userId != idFromToken {
				t.Error("id from access token is wrong")
				return
			}
		})
	}
}

func TestSetRefreshToken(t *testing.T) {
	for _, tc := range GeneralTestCases {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Now()
			newRefreshToken := "newRandomToken"
			userRepository := repository.NewMockUsersRepository()
			us := NewUsersService(userRepository, testConfig)

			err := us.SetRefreshToken(context.Background(), tc.userId, newRefreshToken)

			if tc.expectedError {
				if err == nil {
					t.Error("expected an error")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			user, err := us.UserById(context.Background(), tc.userId)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if user.RefreshToken.Token != newRefreshToken {
				t.Error("token is not stored")
				return
			}

			delta := 1 * time.Minute
			expectedTokenExpiresAt := now.Add(time.Hour * time.Duration(us.cfg.Http.RefreshTokenTTLHours))
			if user.RefreshToken.ExpiresAt.Add(delta).Before(expectedTokenExpiresAt) || user.RefreshToken.ExpiresAt.Add(-delta).After(expectedTokenExpiresAt) {
				t.Error("invalid expiration date")
				return
			}
		})
	}
}
