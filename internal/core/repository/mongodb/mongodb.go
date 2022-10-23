package mongodb

import (
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type Repositories struct {
	Users ports.UsersStore
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
