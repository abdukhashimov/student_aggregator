package services

import (
	"context"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/pkg/hash"
	"github.com/abdukhashimov/student_aggregator/pkg/jwt"
)

var _ ports.UsersService = (*UsersService)(nil)

type UsersService struct {
	repo        ports.UsersStore
	cfg         *config.Config
	hasher      hash.PasswordHasher
	tokeManager jwt.TokenManager
}

func NewUsersService(repo ports.UsersStore, cfg *config.Config) *UsersService {
	hasher := hash.NewSHA1Hasher(cfg.Project.Salt)
	tokeManager := jwt.NewManager(cfg.Project.JwtSecret)

	return &UsersService{
		repo:        repo,
		cfg:         cfg,
		hasher:      hasher,
		tokeManager: tokeManager,
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

func (us *UsersService) SignIn(ctx context.Context, input domain.SignInUserInput) (string, error) {
	hashedPassword, err := us.hasher.Hash(input.Password)
	if err != nil {
		return "", domain.ErrInternalError
	}

	user, err := us.repo.GetByCredentials(ctx, input.Email, hashedPassword)
	if err != nil {
		return "", err
	}

	return user.ID, err
}

func (us *UsersService) UserByAccessToken(ctx context.Context, token string) (*domain.User, error) {
	userId, err := us.tokeManager.Parse(token)
	if err != nil {
		return nil, err
	}

	user, err := us.UserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (us *UsersService) UserByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	user, err := us.repo.GetByRefreshToken(ctx, refreshToken)

	return user, err
}

func (us *UsersService) UserById(ctx context.Context, id string) (*domain.User, error) {
	user, err := us.repo.GetById(ctx, id)

	return user, err
}

func (us *UsersService) GenerateUserTokens(ctx context.Context, id string) (*domain.Tokens, error) {
	accessToken, err := us.tokeManager.NewJWT(id, time.Minute*time.Duration(us.cfg.Http.AccessTokenTTLMinutes))
	if err != nil {
		return nil, err
	}

	refreshToken, err := us.tokeManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = us.SetRefreshToken(ctx, id, refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (us *UsersService) SetRefreshToken(ctx context.Context, id string, token string) error {
	tokenExpiresAt := time.Now().Add(time.Hour * time.Duration(us.cfg.Http.RefreshTokenTTLHours))

	err := us.repo.StoreRefreshToken(ctx, id, domain.RefreshToken{
		Token:     token,
		ExpiresAt: tokenExpiresAt,
	})

	if err != nil {
		return err
	}

	return nil
}
