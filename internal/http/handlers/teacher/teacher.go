package teacher

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rusilkoirala/student-api/internal/http/middleware"
	"github.com/rusilkoirala/student-api/internal/storage"
	"github.com/rusilkoirala/student-api/internal/types"
	"github.com/rusilkoirala/student-api/internal/utils/response"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}

		var t types.Teacher
		err := json.NewDecoder(r.Body).Decode(&t)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(t); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		lastId, err := storage.CreateTeacher(t.Name, t.Email, t.Subject, schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Teacher created", slog.Int64("id", lastId))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}

		intId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		t, err := storage.GetTeacher(intId, schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, t)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}

		teachers, err := storage.GetAllTeachers(schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, teachers)
	}
}

func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}

		intId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err = storage.DeleteTeacher(intId, schoolId); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, map[string]int64{"deleted": intId})
	}
}
