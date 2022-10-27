package ports

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, input domain.SignUpUserInput) error
	SignIn(ctx context.Context, input domain.SignInUserInput) (string, error)
	SetRefreshToken(ctx context.Context, id string, token string) error
	UserById(ctx context.Context, id string) (*domain.User, error)
	GenerateUserTokens(ctx context.Context, id string) (*domain.Tokens, error)
	UserByAccessToken(ctx context.Context, token string) (*domain.User, error)
	UserByRefreshToken(ctx context.Context, token string) (*domain.User, error)
}

type UsersStore interface {
	Create(ctx context.Context, user domain.User) error
	Update(ctx context.Context, inp domain.UpdateUserInput) error
	GetByCredentials(ctx context.Context, email, password string) (*domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	StoreRefreshToken(ctx context.Context, id string, token domain.RefreshToken) (int, error)
}
