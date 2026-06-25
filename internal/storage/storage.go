package storage

import "github.com/rusilkoirala/student-api/internal/types"

type Storage interface {
	// Auth
	CreateSchool(name, username, passwordHash string) (int64, error)
	GetSchoolByUsername(username string) (types.School, error)

	// Students (scoped to school)
	CreateStudent(name string, email string, age int, classId int, schoolId int) (int64, error)
	GetStudent(id int64, schoolId int) (types.Student, error)
	GetAllStudent(schoolId int) ([]types.Student, error)
	Delete(id int64, schoolId int) error
	UpdateStudent(id int64, name string, email string, age int, classId int, schoolId int) error

	// Teachers (scoped to school)
	CreateTeacher(name string, email string, subject string, schoolId int) (int64, error)
	GetTeacher(id int64, schoolId int) (types.Teacher, error)
	GetAllTeachers(schoolId int) ([]types.Teacher, error)
	DeleteTeacher(id int64, schoolId int) error
	UpdateTeacher(id int64, name string, email string, subject string, schoolId int) error

	// Classes (scoped to school)
	CreateClass(name string, teacherId int, schoolId int) (int64, error)
	GetClass(id int64, schoolId int) (types.Class, error)
	GetAllClasses(schoolId int) ([]types.Class, error)
	DeleteClass(id int64, schoolId int) error
	UpdateClass(id int64, name string, teacherId int, schoolId int) error

	// Stats (scoped to school)
	GetStats(schoolId int) (map[string]int, error)
}
