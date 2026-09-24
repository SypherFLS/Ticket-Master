package dto

import "tmaster/internal/constants"

type RequestUser struct {
	Email string
	Role  constants.UserRole
}

type PaginationData struct {
	Limit  int
	Offset int
}
