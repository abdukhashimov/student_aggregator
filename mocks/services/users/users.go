package users

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

var _ ports.UsersService = (*mockUsersService)(nil)

type mockUsersService struct {
	tokenStorage map[string]*domain.User
}

const (
	ValidAccessToken   = "first_token"
	InvalidAccessToken = "second_token"

	ValidUserId = "1"
)

func NewMockUsersService() *mockUsersService {
	return &mockUsersService{
		tokenStorage: map[string]*domain.User{
			ValidAccessToken: {
				ID:           ValidUserId,
				Username:     "Name 1",
				Email:        "test1@ts.ts",
				Password:     "123456",
				RefreshToken: domain.RefreshToken{},
			},
		},
	}
}

func (m *mockUsersService) SignUp(ctx context.Context, input domain.SignUpUserInput) error {
	return nil
}

func (m *mockUsersService) SignIn(ctx context.Context, input domain.SignInUserInput) (string, error) {
	return "", nil
}

func (m *mockUsersService) SetRefreshToken(ctx context.Context, id string, token string) error {
	return nil
}

func (m *mockUsersService) UserById(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User
	for _, us := range m.tokenStorage {
		if us.ID == id {
			user = us
		}
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (m *mockUsersService) GenerateUserTokens(ctx context.Context, id string) (*domain.Tokens, error) {
	return nil, nil
}

func (m *mockUsersService) UserByAccessToken(ctx context.Context, token string) (*domain.User, error) {
	user, ok := m.tokenStorage[token]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (m *mockUsersService) UserByRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	return nil, nil
}
