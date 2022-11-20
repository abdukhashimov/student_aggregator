package users

import (
	"context"
	"errors"
	"github.com/abdukhashimov/student_aggregator/mocks/utils"
	"sync"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

const (
	ValidMongoId           = "507f1f77bcf86cd000000001"
	ExpiredMongoId         = "507f1f77bcf86cd000000002"
	InvalidMongoId         = "123456"
	NotFoundMongoId        = "654321"
	TakenUserEmail         = "test@ts.ts"
	ErrorToCreateUserEmail = "test@ts.ts"
	ValidRefreshToken      = "123456"
	MissedRefreshToken     = "654321"
	ExpiredRefreshToken    = "112233"
)

var _ ports.UsersStore = (*mockUsersRepository)(nil)

type mockUsersRepository struct {
	usersStorage map[string]*domain.User
	lastId       string
	mutex        *sync.RWMutex
}

func NewMockUsersRepository() *mockUsersRepository {
	return &mockUsersRepository{
		usersStorage: map[string]*domain.User{
			ValidMongoId: {
				ID:       ValidMongoId,
				Username: "Name 1",
				Email:    TakenUserEmail,
				Password: "123456",
				RefreshToken: domain.RefreshToken{
					Token:     ValidRefreshToken,
					ExpiresAt: time.Now().Add(1 * time.Hour),
				},
			},
			ExpiredMongoId: {
				ID:       ExpiredMongoId,
				Username: "Name 2",
				Email:    "test2@ts.ts",
				Password: "123456",
				RefreshToken: domain.RefreshToken{
					Token:     ExpiredRefreshToken,
					ExpiresAt: time.Now(),
				},
			},
		},
		lastId: ExpiredMongoId,
		mutex:  &sync.RWMutex{},
	}
}

func (m *mockUsersRepository) Create(ctx context.Context, user domain.User) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if user.Email == ErrorToCreateUserEmail {
		return "", errors.New("internal error")
	}
	if user.Email == "" {
		return "", errors.New("email is missed")
	}
	for _, u := range m.usersStorage {
		if user.Email == u.Email {
			return "", errors.New("email is taken")
		}
	}

	newId := utils.IncrementMongoId(m.lastId)
	m.lastId = newId
	m.usersStorage[newId] = &domain.User{
		ID:           newId,
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: domain.RefreshToken{},
	}

	return newId, nil
}

func (m *mockUsersRepository) Update(ctx context.Context, id string, inp domain.UpdateUserInput) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, ok := m.usersStorage[id]
	if !ok {
		return domain.ErrNotFound
	}

	if inp.Email != nil {
		user.Email = *inp.Email
	}
	if inp.Username != nil {
		user.Username = *inp.Username
	}
	if inp.Username != nil {
		user.Username = *inp.Username
	}
	if inp.RefreshToken != nil {
		user.RefreshToken = *inp.RefreshToken
	}

	return nil
}

func (m *mockUsersRepository) GetByCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, u := range m.usersStorage {
		if u.Email == email && u.Password == password {
			copyUser := *u
			return &copyUser, nil
		}
	}

	return nil, errors.New("user not found")
}

func (m *mockUsersRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, u := range m.usersStorage {
		if u.RefreshToken.Token == refreshToken && time.Now().Before(u.RefreshToken.ExpiresAt) {
			copyUser := *u
			return &copyUser, nil
		}
	}

	return nil, errors.New("user not found")
}

func (m *mockUsersRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	if id == InvalidMongoId {
		return nil, errors.New("internal error")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, ok := m.usersStorage[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copyUser := *user

	return &copyUser, nil
}

func (m *mockUsersRepository) StoreRefreshToken(ctx context.Context, id string, token domain.RefreshToken) error {
	if id == InvalidMongoId {
		return errors.New("internal error")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, ok := m.usersStorage[id]
	if !ok {
		return domain.ErrNotFound
	}
	user.RefreshToken = token

	return nil
}
