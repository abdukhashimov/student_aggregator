package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
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

func (ur *UsersRepo) Create(ctx context.Context, user domain.User) error {

	res, err := ur.db.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	logger.Log.Debugf("Created user with %s", res.InsertedID)

	return nil
}

func (ur *UsersRepo) Update(ctx context.Context, inp domain.UpdateUserInput) error {
	//TODO implement me
	panic("implement me")
}

func (ur *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	user := &domain.User{}
	if err := ur.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (ur *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	user := &domain.User{}
	if err := ur.db.FindOne(ctx, bson.M{
		"refresh_token.token":      refreshToken,
		"refresh_token.expires_at": bson.M{"$gt": time.Now()},
	}).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (ur *UsersRepo) GetById(ctx context.Context, id string) (*domain.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := ur.db.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UsersRepo) StoreRefreshToken(ctx context.Context, id string, token domain.RefreshToken) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = ur.db.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{"refresh_token": token}})
	if err != nil {
		return err
	}

	return nil
}
