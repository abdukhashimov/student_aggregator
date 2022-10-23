package ports

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, input domain.SignUpUserInput) error
	SignIn(ctx context.Context, input domain.SignInUserInput) (*domain.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.Tokens, error)
}

type UsersStore interface {
	Create(ctx context.Context, user domain.User) error
	Update(ctx context.Context, inp domain.UpdateUserInput) error
	GetByCredentials(ctx context.Context, email, password string) (*domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
}
