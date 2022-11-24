package mongodb

import (
	"context"
	"errors"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const nameCollectionWAC = "wac"

var _ ports.SchemaStore = (*SchemaRepo)(nil)

type WACRepo struct {
	db *mongo.Collection
}

func NewWACRepo(db *mongo.Database) *WACRepo {
	return &WACRepo{
		db: db.Collection(nameCollectionWAC),
	}
}

func (wr *WACRepo) AddStudent(ctx context.Context, student domain.StudentWAC) (id string, err error) {
	res, err := wr.db.InsertOne(ctx, student)

	if err != nil {
		return "", err
	}

	idObject := res.InsertedID.(primitive.ObjectID)

	return idObject.String(), nil
}

func (wr *WACRepo) GetById(ctx context.Context, id string) (*domain.StudentWAC, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var student domain.StudentWAC
	result := wr.db.FindOne(ctx, bson.M{
		"_id": objectId,
	})

	if err := result.Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &student, nil
}

func (wr *WACRepo) Update(ctx context.Context, id string, input interface{}) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = wr.db.UpdateByID(ctx, bson.M{"_id": objectId}, bson.M{"$set": input})

	return err
}

func (wr *WACRepo) Delete(ctx context.Context, id string) error{
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return  err
	}

	result, err := wr.db.DeleteOne(ctx, bson.M{
		"_id": objectId,
	})

	if result != nil && result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return err
}
