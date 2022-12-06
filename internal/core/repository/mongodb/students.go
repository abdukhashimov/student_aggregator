package mongodb

import (
	"context"
	"errors"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const studentsCollection = "students"

type StudentCollection interface {
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type StudentsRepo struct {
	col StudentCollection
}

func NewStudentsRepo(db *mongo.Database) *StudentsRepo {
	return newRepo(db.Collection(studentsCollection))
}

func newRepo(col StudentCollection) *StudentsRepo {
	return &StudentsRepo{col}
}

func (sr *StudentsRepo) SaveRSS(ctx context.Context, email string, student domain.StudentRSS) (string, error) {
	s := domain.StudentRecord{
		Source:     domain.RSS,
		Email:      email,
		StudentRSS: student,
	}
	return sr.save(ctx, s)
}

func (sr *StudentsRepo) SaveWAC(ctx context.Context, email string, student domain.StudentWAC) (string, error) {
	s := domain.StudentRecord{
		Source:     domain.WAC,
		Email:      email,
		StudentWAC: student,
	}
	return sr.save(ctx, s)
}

func (sr *StudentsRepo) GetById(ctx context.Context, id string) (*domain.StudentRecord, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var student domain.StudentRecord
	if err := sr.col.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &student, nil
}

func (sr *StudentsRepo) save(ctx context.Context, student domain.StudentRecord) (string, error) {
	res, err := sr.col.InsertOne(ctx, student)
	if err != nil {
		return "", err
	}

	oid := getIdFromObjectID(res.InsertedID)

	return oid, nil
}

func (sr *StudentsRepo) GetAll(ctx context.Context, options domain.ListStudentsOptions) ([]domain.StudentRecord, error) {
	var students []domain.StudentRecord
	opts := getPaginationOpts(options.Limit, options.Skip)
	opts.SetSort(options.Sort)

	filter := bson.M{}
	if options.Email != "" {
		filter["email"] = options.Email
	}
	if options.Source != "" {
		filter["source"] = options.Source
	}

	cur, err := sr.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &students)

	return students, err
}

func (sr *StudentsRepo) Update(ctx context.Context, id string, input domain.StudentRecord) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = sr.col.UpdateOne(ctx,
		bson.M{"_id": objectId}, bson.M{"$set": input})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}
