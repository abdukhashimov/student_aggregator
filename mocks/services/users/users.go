package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

var _ ports.UsersService = (*mockUsersService)(nil)

type mockUsersService struct {
	usersStorage map[string]*domain.User
	lastUserId   string
	mutex        *sync.RWMutex
}

type UserProfileResponse struct {
	User domain.UserProfile `json:"user"`
}

const (
	ValidAccessToken   = "validAccessToken"
	InvalidAccessToken = "invalidAccessToken"

	ValidUserId       = "1"
	ValidUserEmail    = "test1@ts.ts"
	ValidUserPassword = "123456"
	ValidRefreshToken = "validRefreshToken"

	ExpiredAccessToken  = "expiredAccessToken"
	ExpiredUserId       = "2"
	ExpiredUserEmail    = "expired@ts.ts"
	ExpiredUserPassword = "123456"
	ExpiredRefreshToken = "expiredRefreshToken"

	NotFoundUserEmail = "notFoundEmail@ts.ts"

	InternalServerErrorUserId       = "500_UserId"
	InternalServerErrorEmail        = "500@ts.ts"
	InternalServerErrorRefreshToken = "500_RefreshToken"
	InternalServerErrorAccessToken  = "500_AccessToken"
)

var EtalonUser = domain.User{
	ID:       ValidUserId,
	Username: "Name 1",
	Email:    ValidUserEmail,
	Password: ValidUserPassword,
	RefreshToken: domain.RefreshToken{
		Token:     ValidRefreshToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	},
}

var ExpiredUser = domain.User{
	ID:       ExpiredUserId,
	Username: "Name 1",
	Email:    ExpiredUserEmail,
	Password: ExpiredUserPassword,
	RefreshToken: domain.RefreshToken{
		Token:     ExpiredRefreshToken,
		ExpiresAt: time.Now(),
	},
}

var EtalonProfileJson, _ = json.Marshal(UserProfileResponse{
	User: *EtalonUser.GetProfile(),
})

func NewMockUsersService() *mockUsersService {
	return &mockUsersService{
		usersStorage: map[string]*domain.User{
			ValidAccessToken:   &EtalonUser,
			ExpiredAccessToken: &ExpiredUser,
		},
		lastUserId: ExpiredUserId,
		mutex:      &sync.RWMutex{},
	}
}

func (m *mockUsersService) SignUp(ctx context.Context, input domain.SignUpUserInput) (string, error) {
	if input.Email == InternalServerErrorEmail {
		return "", domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, u := range m.usersStorage {
		if u.Email == input.Email {
			return "", errors.New("email is taken")
		}
	}

	m.incrementId()
	accessToken := fmt.Sprintf("access_token_%d", time.Now().UnixNano())
	m.usersStorage[accessToken] = &domain.User{
		ID:       m.lastUserId,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RefreshToken: domain.RefreshToken{
			Token:     fmt.Sprintf("refresh_token_%d", time.Now().UnixNano()),
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}
	return m.lastUserId, nil
}

func (m *mockUsersService) SignIn(ctx context.Context, input domain.SignInUserInput) (string, error) {
	if input.Email == InternalServerErrorEmail {
		return "", domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, u := range m.usersStorage {
		if u.Email == input.Email && u.Password == input.Password {
			return u.ID, nil
		}
	}

	return "", domain.ErrNotFound
}

func (m *mockUsersService) SetRefreshToken(ctx context.Context, id string, token string) error {
	if id == InternalServerErrorUserId {
		return domain.ErrInternalError
	}

	return nil
}

func (m *mockUsersService) UserById(ctx context.Context, id string) (*domain.User, error) {
	if id == InternalServerErrorUserId {
		return nil, domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, us := range m.usersStorage {
		if us.ID == id {
			usCopy := *us
			return &usCopy, nil
		}
	}

	return nil, domain.ErrNotFound
}

func (m *mockUsersService) GenerateUserTokens(ctx context.Context, id string) (*domain.Tokens, error) {
	if id == InternalServerErrorUserId {
		return nil, domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for accessToken, u := range m.usersStorage {
		if u.ID == id {
			tokens := domain.Tokens{
				AccessToken:  accessToken,
				RefreshToken: u.RefreshToken.Token,
			}

			return &tokens, nil
		}
	}

	return nil, domain.ErrNotFound
}

func (m *mockUsersService) UserByAccessToken(ctx context.Context, token string) (*domain.User, error) {
	if token == InternalServerErrorAccessToken {
		return nil, domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, ok := m.usersStorage[token]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return user, nil
}

func (m *mockUsersService) UserByRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	if token == InternalServerErrorRefreshToken {
		return nil, domain.ErrInternalError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, u := range m.usersStorage {
		if u.RefreshToken.Token == token && u.RefreshToken.ExpiresAt.After(time.Now()) {
			userCopy := *u
			return &userCopy, nil
		}
	}

	return nil, domain.ErrNotFound
}

func (m *mockUsersService) incrementId() {
	id, _ := strconv.Atoi(m.lastUserId)
	id++
	m.lastUserId = strconv.Itoa(id)
}
