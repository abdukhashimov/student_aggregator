package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
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

func (u *UsersRepo) GetById(ctx context.Context, id string) (*domain.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := u.db.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
