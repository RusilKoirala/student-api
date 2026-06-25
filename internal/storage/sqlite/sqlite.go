package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rusilkoirala/student-api/internal/config"
	"github.com/rusilkoirala/student-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schools (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT NOT NULL,
			username      TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS teachers (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			name      TEXT    NOT NULL,
			email     TEXT    NOT NULL,
			subject   TEXT    NOT NULL,
			school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS classes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL,
			teacher_id INTEGER REFERENCES teachers(id) ON DELETE SET NULL,
			school_id  INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS students (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			name      TEXT    NOT NULL,
			email     TEXT    NOT NULL,
			age       INTEGER NOT NULL,
			class_id  INTEGER REFERENCES classes(id) ON DELETE SET NULL,
			school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

// ── Auth ──────────────────────────────────────────────────────────────────────

func (s *Sqlite) CreateSchool(name, username, passwordHash string) (int64, error) {
	result, err := s.Db.Exec(
		"INSERT INTO schools (name, username, password_hash) VALUES (?, ?, ?)",
		name, username, passwordHash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetSchoolByUsername(username string) (types.School, error) {
	var sc types.School
	err := s.Db.QueryRow(
		"SELECT id, name, username, password_hash FROM schools WHERE username = ? LIMIT 1",
		username,
	).Scan(&sc.Id, &sc.Name, &sc.Username, &sc.PasswordHash)
	return sc, err
}

// ── Students ──────────────────────────────────────────────────────────────────

func (s *Sqlite) CreateStudent(name, email string, age, classId, schoolId int) (int64, error) {
	var classArg interface{}
	if classId > 0 {
		classArg = classId
	}
	result, err := s.Db.Exec(
		"INSERT INTO students (name, email, age, class_id, school_id) VALUES (?, ?, ?, ?, ?)",
		name, email, age, classArg, schoolId,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetStudent(id int64, schoolId int) (types.Student, error) {
	var st types.Student
	err := s.Db.QueryRow(
		"SELECT id, name, email, age, COALESCE(class_id,0) FROM students WHERE id=? AND school_id=? LIMIT 1",
		id, schoolId,
	).Scan(&st.Id, &st.Name, &st.Email, &st.Age, &st.ClassId)
	return st, err
}

func (s *Sqlite) GetAllStudent(schoolId int) ([]types.Student, error) {
	rows, err := s.Db.Query(
		"SELECT id, name, email, age, COALESCE(class_id,0) FROM students WHERE school_id=?",
		schoolId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var st types.Student
		if err := rows.Scan(&st.Id, &st.Name, &st.Email, &st.Age, &st.ClassId); err != nil {
			return nil, err
		}
		students = append(students, st)
	}
	if students == nil {
		students = []types.Student{}
	}
	return students, nil
}

func (s *Sqlite) Delete(id int64, schoolId int) error {
	_, err := s.Db.Exec("DELETE FROM students WHERE id=? AND school_id=?", id, schoolId)
	return err
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int, classId int, schoolId int) error {
	var classArg interface{}
	if classId > 0 {
		classArg = classId
	}
	_, err := s.Db.Exec(
		"UPDATE students SET name=?, email=?, age=?, class_id=? WHERE id=? AND school_id=?",
		name, email, age, classArg, id, schoolId,
	)
	return err
}

// ── Teachers ──────────────────────────────────────────────────────────────────

func (s *Sqlite) CreateTeacher(name, email, subject string, schoolId int) (int64, error) {
	result, err := s.Db.Exec(
		"INSERT INTO teachers (name, email, subject, school_id) VALUES (?, ?, ?, ?)",
		name, email, subject, schoolId,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetTeacher(id int64, schoolId int) (types.Teacher, error) {
	var t types.Teacher
	err := s.Db.QueryRow(
		"SELECT id, name, email, subject FROM teachers WHERE id=? AND school_id=? LIMIT 1",
		id, schoolId,
	).Scan(&t.Id, &t.Name, &t.Email, &t.Subject)
	return t, err
}

func (s *Sqlite) GetAllTeachers(schoolId int) ([]types.Teacher, error) {
	rows, err := s.Db.Query(
		"SELECT id, name, email, subject FROM teachers WHERE school_id=?",
		schoolId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []types.Teacher
	for rows.Next() {
		var t types.Teacher
		if err := rows.Scan(&t.Id, &t.Name, &t.Email, &t.Subject); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	if teachers == nil {
		teachers = []types.Teacher{}
	}
	return teachers, nil
}

func (s *Sqlite) DeleteTeacher(id int64, schoolId int) error {
	_, err := s.Db.Exec("DELETE FROM teachers WHERE id=? AND school_id=?", id, schoolId)
	return err
}

func (s *Sqlite) UpdateTeacher(id int64, name string, email string, subject string, schoolId int) error {
	_, err := s.Db.Exec(
		"UPDATE teachers SET name=?, email=?, subject=? WHERE id=? AND school_id=?",
		name, email, subject, id, schoolId,
	)
	return err
}

// ── Classes ───────────────────────────────────────────────────────────────────

func (s *Sqlite) CreateClass(name string, teacherId, schoolId int) (int64, error) {
	var teacherArg interface{}
	if teacherId > 0 {
		teacherArg = teacherId
	}
	result, err := s.Db.Exec(
		"INSERT INTO classes (name, teacher_id, school_id) VALUES (?, ?, ?)",
		name, teacherArg, schoolId,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetClass(id int64, schoolId int) (types.Class, error) {
	var c types.Class
	err := s.Db.QueryRow(`
		SELECT c.id, c.name, COALESCE(c.teacher_id,0), COALESCE(t.name,'')
		FROM classes c
		LEFT JOIN teachers t ON t.id = c.teacher_id
		WHERE c.id=? AND c.school_id=? LIMIT 1`,
		id, schoolId,
	).Scan(&c.Id, &c.Name, &c.TeacherId, &c.TeacherName)
	return c, err
}

func (s *Sqlite) GetAllClasses(schoolId int) ([]types.Class, error) {
	rows, err := s.Db.Query(`
		SELECT c.id, c.name, COALESCE(c.teacher_id,0), COALESCE(t.name,'')
		FROM classes c
		LEFT JOIN teachers t ON t.id = c.teacher_id
		WHERE c.school_id=?`,
		schoolId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []types.Class
	for rows.Next() {
		var c types.Class
		if err := rows.Scan(&c.Id, &c.Name, &c.TeacherId, &c.TeacherName); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	if classes == nil {
		classes = []types.Class{}
	}
	return classes, nil
}

func (s *Sqlite) DeleteClass(id int64, schoolId int) error {
	_, err := s.Db.Exec("DELETE FROM classes WHERE id=? AND school_id=?", id, schoolId)
	return err
}

func (s *Sqlite) UpdateClass(id int64, name string, teacherId int, schoolId int) error {
	var teacherArg interface{}
	if teacherId > 0 {
		teacherArg = teacherId
	}
	_, err := s.Db.Exec(
		"UPDATE classes SET name=?, teacher_id=? WHERE id=? AND school_id=?",
		name, teacherArg, id, schoolId,
	)
	return err
}

// ── Stats ─────────────────────────────────────────────────────────────────────

func (s *Sqlite) GetStats(schoolId int) (map[string]int, error) {
	stats := map[string]int{"students": 0, "teachers": 0, "classes": 0}

	rows, err := s.Db.Query(`
		SELECT 'students', COUNT(*) FROM students WHERE school_id=?
		UNION ALL
		SELECT 'teachers', COUNT(*) FROM teachers WHERE school_id=?
		UNION ALL
		SELECT 'classes',  COUNT(*) FROM classes  WHERE school_id=?`,
		schoolId, schoolId, schoolId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var val int
		if err := rows.Scan(&key, &val); err != nil {
			return nil, err
		}
		stats[key] = val
	}
	return stats, nil
}
