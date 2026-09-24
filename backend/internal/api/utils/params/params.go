package params

import (
	"net/http"
	"strconv"
	"tmaster/internal/apperrors"
	"tmaster/internal/constants"
	"tmaster/internal/dto"
)

func GetAllParams(r *http.Request) (int, constants.UserRole, error) {
	userID, err := GetUserID(r)
	userRole, erro := GetUserRole(r)

	if err != nil {
		return 0, "", err
	} else if erro != nil {
		return 0, "", erro
	}

	return userID, userRole, nil
}

func GetUserRole(r *http.Request) (constants.UserRole, error) {
	value := r.Context().Value(constants.UserRoleKey)

	userRole, ok := value.(string)

	if !ok {
		return "", apperrors.WrongParamType
	}

	return constants.UserRole(userRole), nil
}

func GetUserID(r *http.Request) (int, error) {
	value := r.Context().Value(constants.UserIDKey)

	if value == "" {
		return 0, apperrors.BlankUserID
	}

	userID, ok := value.(int)

	if !ok {
		return 0, apperrors.WrongParamType
	}

	return userID, nil
}

func GetPagData(r *http.Request) dto.PaginationData {
	return dto.PaginationData{
		Limit:  GetLimit(r),
		Offset: GetOffset(r),
	}
}

func GetLimit(r *http.Request) int {
	raw_lim := r.URL.Query().Get("limit")

	if raw_lim == "" {
		return 10
	}

	lim, err := strconv.Atoi(raw_lim)

	if err != nil {
		return 10
	}

	return lim
}

func GetOffset(r *http.Request) int {
	raw_off := r.URL.Query().Get("offset")

	if raw_off == "" {
		return 0
	}

	off, err := strconv.Atoi(raw_off)

	if err != nil {
		return 0
	}

	return off
}
