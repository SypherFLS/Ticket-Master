package apperrors

import "errors"

var (
	WrongSignMethod    = errors.New("wrong sign method")
	InvalidToken       = errors.New("invalid token")
	BlankUserID        = errors.New("blank user id")
	WrongParamType     = errors.New("wrong param type")
	TicketQueueIsEmpty = errors.New("ticket queue is empty")
	NoPermission       = errors.New("no permisssion")
	NothingChanged     = errors.New("nothing changed")
)
