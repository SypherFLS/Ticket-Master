package helpers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tmaster/internal/apperrors"
	"tmaster/internal/constants"
	"tmaster/internal/dto"
)

func GetAllParams(c *gin.Context) (int, constants.UserRole, error) {
	userID, err := GetUserID(c)
	userRole, erro := GetUserRole(c)

	if err != nil {
		return 0, "", err
	} else if erro != nil {
		return 0, "", erro
	}

	return userID, userRole, nil
}

func GetUserRole(c *gin.Context) (constants.UserRole, error) {
	value, exists := c.Get(constants.UserRoleKey)

	userRole, ok := value.(string)

	if !ok || !exists {
		return "", apperrors.WrongParamType
	}

	return constants.UserRole(userRole), nil
}

func GetUserID(c *gin.Context) (int, error) {
	value, exists := c.Get(constants.UserIDKey)

	if !exists {
		return 0, apperrors.BlankUserID
	}

	userID, ok := value.(int)

	if !ok {
		return 0, apperrors.WrongParamType
	}

	return userID, nil
}

func GetPagData(c *gin.Context) dto.PaginationData {
	return dto.PaginationData{
		Limit:  GetLimit(c),
		Offset: GetOffset(c),
	}
}

func GetLimit(c *gin.Context) int {
	raw_lim := c.Request.URL.Query().Get("limit")

	if raw_lim == "" {
		return 10
	}

	lim, err := strconv.Atoi(raw_lim)

	if err != nil {
		return 10
	}

	return lim
}

func GetOffset(c *gin.Context) int {
	raw_off := c.Request.URL.Query().Get("offset")

	if raw_off == "" {
		return 0
	}

	off, err := strconv.Atoi(raw_off)

	if err != nil {
		return 0
	}

	return off
}
