package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/rusilkoirala/student-api/internal/storage"
	"github.com/rusilkoirala/student-api/internal/types"
	"github.com/rusilkoirala/student-api/internal/utils/response"
)

func Register(store storage.Storage, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("request body is empty")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		schoolId, err := store.CreateSchool(req.SchoolName, req.Username, string(hash))
		if err != nil {
			// username already taken
			response.WriteJson(w, http.StatusConflict, response.GeneralError(errors.New("username already taken")))
			return
		}

		token, err := signToken(int(schoolId), req.SchoolName, jwtSecret)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, types.AuthResponse{
			Token:      token,
			SchoolName: req.SchoolName,
			SchoolId:   int(schoolId),
		})
	}
}

func Login(store storage.Storage, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("request body is empty")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		school, err := store.GetSchoolByUsername(req.Username)
		if err != nil {
			// Don't reveal whether the username exists
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("invalid username or password")))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(school.PasswordHash), []byte(req.Password)); err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("invalid username or password")))
			return
		}

		token, err := signToken(school.Id, school.Name, jwtSecret)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, types.AuthResponse{
			Token:      token,
			SchoolName: school.Name,
			SchoolId:   school.Id,
		})
	}
}

func signToken(schoolId int, schoolName, secret string) (string, error) {
	claims := jwt.MapClaims{
		"schoolId":   schoolId,
		"schoolName": schoolName,
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
