package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rusilkoirala/student-api/internal/utils/response"
)

type contextKey string

const SchoolIDKey contextKey = "schoolId"
const SchoolNameKey contextKey = "schoolName"

func RequireAuth(jwtSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(
				&authError{"missing or invalid Authorization header"},
			))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, &authError{"unexpected signing method"}
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(
				&authError{"invalid or expired token"},
			))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(
				&authError{"invalid token claims"},
			))
			return
		}

		schoolId, ok := claims["schoolId"].(float64)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(
				&authError{"schoolId missing from token"},
			))
			return
		}

		schoolName, _ := claims["schoolName"].(string)

		ctx := context.WithValue(r.Context(), SchoolIDKey, int(schoolId))
		ctx = context.WithValue(ctx, SchoolNameKey, schoolName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SchoolIDFromCtx pulls school_id out of context — panics with 401 if missing.
func SchoolIDFromCtx(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, ok := r.Context().Value(SchoolIDKey).(int)
	if !ok || id == 0 {
		response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(
			&authError{"not authenticated"},
		))
		return 0, false
	}
	return id, true
}

// ── tiny error type ───────────────────────────────────────────────────────────

type authError struct{ msg string }

func (e *authError) Error() string { return e.msg }
