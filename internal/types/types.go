package types

import (
	"encoding/json"
	"strconv"
)

type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age   int    `json:"age" validate:"required"`
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

	// Handle id as either string or number
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
