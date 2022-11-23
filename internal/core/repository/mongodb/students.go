package mongodb

import (
	"context"
	"fmt"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//const studentsCollection = "students"

type StudentCollection interface {
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type StudentsRepo struct {
	col StudentCollection
}

//func NewStudentsRepo(db *mongo.Database) *StudentsRepo {
//	return newRepo(db.Collection(studentsCollection))
//}

func newRepo(col StudentCollection) *StudentsRepo {
	return &StudentsRepo{col}
}

func (sr *StudentsRepo) SaveRSS(ctx context.Context, student domain.StudentRSS) (string, error) {
	s := domain.StudentRecord{
		Source:     domain.RSS,
		StudentRSS: student,
	}
	return sr.save(ctx, s)
}

func (sr *StudentsRepo) SaveWAC(ctx context.Context, student domain.StudentWAC) (string, error) {
	s := domain.StudentRecord{
		Source:     domain.WAC,
		StudentWAC: student,
	}
	return sr.save(ctx, s)
}

func (sr *StudentsRepo) save(ctx context.Context, student domain.StudentRecord) (string, error) {
	res, err := sr.col.InsertOne(ctx, student)
	if err != nil {
		return "", err
	}

	oid := getIdFromObjectID(res.InsertedID)

	return oid, nil
}
