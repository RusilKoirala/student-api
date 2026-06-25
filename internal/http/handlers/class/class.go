package class

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

		var c types.Class
		err := json.NewDecoder(r.Body).Decode(&c)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(c); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		lastId, err := storage.CreateClass(c.Name, c.TeacherId, schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Class created", slog.Int64("id", lastId))
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

		c, err := storage.GetClass(intId, schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, c)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}

		classes, err := storage.GetAllClasses(schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, classes)
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

		if err = storage.DeleteClass(intId, schoolId); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, map[string]int64{"deleted": intId})
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
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

		var c types.Class
		err = json.NewDecoder(r.Body).Decode(&c)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(c); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		if err = storage.UpdateClass(intId, c.Name, c.TeacherId, schoolId); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Class updated", slog.Int64("id", intId))
		response.WriteJson(w, http.StatusOK, map[string]int64{"updated": intId})
	}
}
