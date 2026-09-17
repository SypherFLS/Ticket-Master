package apperrors

import "errors"

var (
	WrongSignMethod = errors.New("wrong sign method")
	InvalidToken    = errors.New("invalid token")
	BlankUserID     = errors.New("blank user id")
	WrongUserID     = errors.New("wrong user id")
)
