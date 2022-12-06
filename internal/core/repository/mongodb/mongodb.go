package mongodb

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abdukhashimov/student_aggregator/internal/core/repository"
)

func NewRepositories(db *mongo.Database) *repository.Repositories {
	return &repository.Repositories{
		Users:    NewUsersRepo(db),
		Schemas:  NewSchemaRepo(db),
		Students: NewStudentsRepo(db),
	}
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}

func getIdFromObjectID(in interface{}) string {
	if p, ok := in.(primitive.ObjectID); ok {
		return p.Hex()
	}

	return ""
}

func getPaginationOpts(l int, s int) *options.FindOptions {
	limit := int64(l)
	skip := int64(s)
	return &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}
}
