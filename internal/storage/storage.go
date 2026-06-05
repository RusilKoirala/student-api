package storage

import "github.com/rusilkoirala/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudent(id int64) (types.Student, error)
	GetAllStudent() ([]types.Student, error)
	Delete(id int64) error
}
