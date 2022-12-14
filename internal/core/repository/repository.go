package repository

import (
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
)

type Repositories struct {
	Users    ports.UsersStore
	Schemas  ports.SchemaStore
	Students ports.StudentsStore
}
