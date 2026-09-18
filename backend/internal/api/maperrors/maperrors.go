package maperrors

import (
	"errors"
	"net/http"
	"tmaster/internal/apperrors"
)

func StatusFromErr(err error) int {
	switch {
	case errors.Is(err, apperrors.WrongSignMethod) || errors.Is(err, apperrors.InvalidToken) || errors.Is(err, apperrors.BlankUserID) || errors.Is(err, apperrors.WrongParamType):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}