package dto

import "tmaster/internal/constants"

type RequsetUser struct {
	ID   int
	Role constants.UserRole
}

type PaginationData struct {
	Limit  int
	Offset int
}
