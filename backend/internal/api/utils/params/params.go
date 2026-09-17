package params

import (
	"net/http"
	"tmaster/internal/apperrors"
)

type contextKey string

const UserIDKey contextKey = "userID"

func GetUserID(r *http.Request) (int, error) {
	value := r.Context().Value(UserIDKey)

	if value == "" {
		return 0, apperrors.BlankUserID
	}

	userID, ok := value.(int)

	if !ok {
		return 0, apperrors.WrongUserID
	}

	return userID, nil
}