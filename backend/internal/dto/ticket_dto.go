package dto

import (
	"time"
	"tmaster/internal/constants"
	"tmaster/internal/repository/models"
)

type TicketDTO struct {
	Title       string `json:"title" validate:"required,min=5,max=30"`
	Description string `json:"decription" validate:"max:150"`
}

type IDResponse struct {
	ID int `json:"ID"`
}

func TicketToModel(tdto TicketDTO, user_id int) models.Ticket {
	now := time.Now()
	return models.Ticket{
		Title: tdto.Title,
		Description: tdto.Description,
		Status : constants.StatusNew,
		CreatedAt: now,
		UserID: user_id,
	}
}