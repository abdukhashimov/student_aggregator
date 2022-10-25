package services

import (
	"context"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/pkg/hash"
)

var _ ports.UsersService = (*UsersService)(nil)

type UsersService struct {
	repo   ports.UsersStore
	hasher hash.PasswordHasher
}

func NewUsersService(repo ports.UsersStore, cfg *config.Config) *UsersService {
	hasher := hash.NewSHA1Hasher(cfg.Project.Salt)

	return &UsersService{
		repo:   repo,
		hasher: hasher,
	}
}

func (us *UsersService) SignUp(ctx context.Context, input domain.SignUpUserInput) error {
	hashedPassword, err := us.hasher.Hash(input.Password)
	if err != nil {
		logger.Log.Errorf("error to create user: %w", err)
		return domain.ErrInternalError
	}

	err = us.repo.Create(ctx, domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

func (us *UsersService) SignIn(ctx context.Context, input domain.SignInUserInput) (*domain.Tokens, error) {
	return nil, nil
}

func (us *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.Tokens, error) {
	return nil, nil
}

func (us *UsersService) UserById(ctx context.Context, id string) (*domain.User, error) {
	user, err := us.repo.GetById(ctx, id)

	return user, err
}
