package mongodb

import (
	"context"
	"errors"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	schemasCollection = "schemas"
)

var _ ports.SchemaStore = (*SchemaRepo)(nil)

type SchemaRepo struct {
	db *mongo.Collection
}

func NewSchemaRepo(db *mongo.Database) *SchemaRepo {
	return &SchemaRepo{
		db: db.Collection(schemasCollection),
	}
}

func (sr *SchemaRepo) Create(ctx context.Context, schema domain.Schema) (string, error) {

	res, err := sr.db.InsertOne(ctx, schema)

	if err != nil {
		if IsDuplicate(err) {
			return "", domain.DuplicationError
		}
		return "", err
	}

	stringId := getIdFromObjectID(res.InsertedID)

	logger.Log.Debugf("new schema created - %s", stringId)

	return stringId, nil
}

func (sr *SchemaRepo) GetById(ctx context.Context, id string) (*domain.Schema, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var schema domain.Schema
	if err := sr.db.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&schema); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &schema, nil
}

func (sr *SchemaRepo) Update(ctx context.Context, id string, input domain.UpdateSchemaInput) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = sr.db.UpdateOne(ctx,
		bson.M{"_id": objectId}, bson.M{"$set": input})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotFound
		}

		if IsDuplicate(err) {
			return domain.DuplicationError
		}
		return err
	}

	return nil
}

func (sr *SchemaRepo) FindAll(ctx context.Context) ([]domain.Schema, error) {
	var schemas []domain.Schema

	cur, err := sr.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &schemas)

	return schemas, err
}

func (sr *SchemaRepo) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := sr.db.DeleteOne(ctx, bson.M{"_id": objectId})
	if res != nil && res.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return err
}
