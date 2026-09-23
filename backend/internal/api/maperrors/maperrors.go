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
	case errors.Is(err, apperrors.TicketQueueIsEmpty):
		return http.StatusOK // because it is not a real service error
	case errors.Is(err, apperrors.NoPermission):
		return http.StatusForbidden
	case errors.Is(err, apperrors.NothingChanged):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
