package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/abdukhashimov/student_aggregator/internal/core/repository"
)

const (
	usersCollection = "users"
)

func NewRepositories(db *mongo.Database) *repository.Repositories {
	return &repository.Repositories{
		Users: NewUsersRepo(db),
	}
}
