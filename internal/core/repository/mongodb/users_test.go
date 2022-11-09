package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const (
	testDbName = "testDbName"
	validId    = "507f1f77bcf86cd799439011"
)

var CreateTestCases = []struct {
	name        string
	mongoRes    bson.D
	expectedErr bool
}{
	{
		name:        "success",
		mongoRes:    mtest.CreateSuccessResponse(),
		expectedErr: false,
	},
	{
		name: "failure",
		mongoRes: mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    123,
			Message: "some error",
		}),
		expectedErr: true,
	},
}

var GetUserResponses = []struct {
	name        string
	getMongoRes func(*domain.User) ([]bson.D, error)
	expectedErr error
}{
	{
		name: "success",
		getMongoRes: func(expectedUser *domain.User) ([]bson.D, error) {
			bsonD, err := toBson(expectedUser)
			if err != nil {
				return nil, err
			}

			res := mtest.CreateCursorResponse(
				1,
				fmt.Sprintf("%s.%s", testDbName, usersCollection),
				mtest.FirstBatch,
				bsonD)

			end := mtest.CreateCursorResponse(
				0,
				fmt.Sprintf("%s.%s", testDbName, usersCollection),
				mtest.NextBatch)

			return []bson.D{res, end}, nil
		},
		expectedErr: nil,
	},
	{
		name: "not-found",
		getMongoRes: func(expectedUser *domain.User) ([]bson.D, error) {
			res := mtest.CreateCursorResponse(
				1,
				fmt.Sprintf("%s.%s", testDbName, usersCollection),
				mtest.FirstBatch)
			end := mtest.CreateCursorResponse(
				0,
				fmt.Sprintf("%s.%s", testDbName, usersCollection),
				mtest.NextBatch)

			return []bson.D{res, end}, nil
		},
		expectedErr: domain.ErrUserNotFound,
	},
}

var GetUserTestCases = []struct {
	name         string
	testFunction func(store ports.UsersStore) (*domain.User, error)
}{
	{
		name: "GetByCredentials",
		testFunction: func(store ports.UsersStore) (*domain.User, error) {
			return store.GetByCredentials(context.Background(), "test@ts.ts", "123456")
		},
	},
	{
		name: "GetByRefreshToken",
		testFunction: func(store ports.UsersStore) (*domain.User, error) {
			return store.GetByRefreshToken(context.Background(), "token")
		},
	},
	{
		name: "GetById",
		testFunction: func(store ports.UsersStore) (*domain.User, error) {
			return store.GetById(context.Background(), validId)
		},
	},
}

var StoreRefreshTokenTestCases = []struct {
	name        string
	getMongoRes func(*domain.User) bson.D
	expectedErr bool
}{
	{
		name: "success",
		getMongoRes: func(expectedUser *domain.User) bson.D {
			return bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "_id", Value: expectedUser.ID},
				}},
			}
		},
	},
	{
		name: "failure",
		getMongoRes: func(user *domain.User) bson.D {
			return mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "some error",
			})
		},
		expectedErr: true,
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

func TestCreateUser(t *testing.T) {
	mt := getMockTest(t)
	defer mt.Close()

	for _, tc := range CreateTestCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.mongoRes)
			usersRepo := NewUsersRepo(mt.DB)

			err := usersRepo.Create(context.Background(), domain.User{
				Username:     "Test",
				Email:        "test@ts.ts",
				Password:     "123456",
				RefreshToken: domain.RefreshToken{},
			})

			if tc.expectedErr && err == nil {
				t.Error("expected an error")
				return
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	mt := getMockTest(t)
	defer mt.Close()

	for _, tc := range GetUserTestCases {
		for _, ur := range GetUserResponses {
			mt.Run(fmt.Sprintf("%s.%s", tc.name, ur.name), func(mt *mtest.T) {
				testUser := domain.User{
					ID:       validId,
					Username: "Max",
					Email:    "test@ts.ts",
					Password: "123456",
					RefreshToken: domain.RefreshToken{
						Token:     "token",
						ExpiresAt: time.Now(),
					},
				}
				moks, err := ur.getMongoRes(&testUser)
				if err != nil {
					t.Errorf("unexpecting error: %s", err.Error())
				}
				mt.AddMockResponses(moks...)
				usersRepo := NewUsersRepo(mt.DB)

				gotUser, err := tc.testFunction(usersRepo)
				if ur.expectedErr != nil {
					if err == nil {
						t.Error("an error expected")
						return
					}
					if ur.expectedErr != err {
						t.Errorf("unexpecting error: %s", err.Error())
						return
					}
					return
				}

				if err != nil {
					t.Errorf("unexpecting error: %s", err.Error())
					return
				}

				err = compareUsers(&testUser, gotUser)
				if err != nil {
					t.Error(err.Error())
					return
				}
			})
		}
	}
}

func TestStoreRefreshToken(t *testing.T) {
	mt := getMockTest(t)
	defer mt.Close()

	for _, tc := range StoreRefreshTokenTestCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			testUser := domain.User{
				ID:       validId,
				Username: "Max",
				Email:    "test@ts.ts",
				Password: "123456",
				RefreshToken: domain.RefreshToken{
					Token:     "token",
					ExpiresAt: time.Now(),
				},
			}
			moks := tc.getMongoRes(&testUser)
			mt.AddMockResponses(moks)
			usersRepo := NewUsersRepo(mt.DB)

			err := usersRepo.StoreRefreshToken(context.Background(), testUser.ID, domain.RefreshToken{
				Token:     "123456",
				ExpiresAt: time.Now(),
			})

			if tc.expectedErr && err == nil {
				t.Error("expected an error")
				return
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}
		})
	}
}

func getMockTest(t *testing.T) *mtest.T {
	opts := mtest.NewOptions().DatabaseName(testDbName).ClientType(mtest.Mock)

	return mtest.New(t, opts)
}

func compareUsers(expected, got *domain.User) error {
	if expected.ID != got.ID {
		return getDoNotMatchError("ID")
	}
	if expected.Email != got.Email {
		return getDoNotMatchError("email")
	}
	if expected.Username != got.Username {
		return getDoNotMatchError("username")
	}
	if expected.Password != got.Password {
		return getDoNotMatchError("password")
	}
	if expected.RefreshToken.Token != got.RefreshToken.Token {
		return getDoNotMatchError("token")
	}
	if expected.RefreshToken.ExpiresAt.UnixMilli() != got.RefreshToken.ExpiresAt.UnixMilli() {
		return getDoNotMatchError("expiresAt")
	}

	return nil
}

func getDoNotMatchError(field string) error {
	return fmt.Errorf("%s field does not match expected", field)
}

func toBson(v interface{}) (bsonD bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &bsonD)
	return
}
