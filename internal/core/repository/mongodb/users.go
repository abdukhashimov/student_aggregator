package mongodb

import (
	"context"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ ports.UsersStore = (*UsersRepo)(nil)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (u *UsersRepo) Create(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UsersRepo) Update(ctx context.Context, inp domain.UpdateUserInput) error {
	//TODO implement me
	panic("implement me")
}

func (u *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}
