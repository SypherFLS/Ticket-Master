package params

import (
	"net/http"
	"tmaster/internal/apperrors"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func GetUserRole(r *http.Request) (string, error) {
	value := r.Context().Value(UserRoleKey)

	userRole, ok := value.(string)

	if !ok {
		return "", apperrors.WrongParamType
	}

	return userRole, nil
}

func GetUserID(r *http.Request) (int, error) {
	value := r.Context().Value(UserIDKey)

	if value == "" {
		return 0, apperrors.BlankUserID
	}

	userID, ok := value.(int)

	if !ok {
		return 0, apperrors.WrongParamType
	}

	return userID, nil
}