package types

import (
	"encoding/json"
	"strconv"
)

// ── Student ──────────────────────────────────────────────────────────────────

type Student struct {
	Id      int    `json:"id"`
	Name    string `json:"name"    validate:"required"`
	Email   string `json:"email"   validate:"required"`
	Age     int    `json:"age"     validate:"required"`
	ClassId int    `json:"classId"`
}

func (s *Student) UnmarshalJSON(data []byte) error {
	type Alias Student
	aux := &struct {
		Id interface{} `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch v := aux.Id.(type) {
	case float64:
		s.Id = int(v)
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		s.Id = id
	}
	return nil
}

// ── School ────────────────────────────────────────────────────────────────────

type School struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // never serialised
}

type RegisterRequest struct {
	SchoolName string `json:"schoolName" validate:"required"`
	Username   string `json:"username"   validate:"required,min=3"`
	Password   string `json:"password"   validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token      string `json:"token"`
	SchoolName string `json:"schoolName"`
	SchoolId   int    `json:"schoolId"`
}

type Teacher struct {
	Id      int    `json:"id"`
	Name    string `json:"name"    validate:"required"`
	Email   string `json:"email"   validate:"required"`
	Subject string `json:"subject" validate:"required"`
}

// ── Class ────────────────────────────────────────────────────────────────────

type Class struct {
	Id        int    `json:"id"`
	Name      string `json:"name"      validate:"required"`
	TeacherId int    `json:"teacherId"`
	// populated on read, not stored directly
	TeacherName string `json:"teacherName,omitempty"`
}
