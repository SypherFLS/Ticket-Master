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

type UserTicketResponse struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      constants.TicketStatus `json:"status"`
	OperatorID  *int                   `json:"operator_id"`
	CreatedAt   time.Time              `json:"created_at"`
	ClaimedAt   *time.Time             `json:"claimed_at"`
}

func ModelToUserTicket(data models.Ticket) UserTicketResponse {
	return UserTicketResponse{
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
		OperatorID:  data.OperatorID,
		CreatedAt:   data.CreatedAt,
		ClaimedAt:   data.ClaimedAt,
	}
}

func ManyMTUT(data []models.Ticket) []UserTicketResponse {
	resp := []UserTicketResponse{}
	for _, r := range data {
		resp = append(resp, ModelToUserTicket(r))
	}
	return resp
}

type OperatorTicketResponse struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"decription"`
	Status      constants.TicketStatus `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	ClaimedAt   *time.Time             `json:"claimed_at"`
	UserID      int                    `json:"user_id"`
}

func ModelToOperatorTicket(data models.Ticket) OperatorTicketResponse {
	return OperatorTicketResponse{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		ClaimedAt:   data.ClaimedAt,
		UserID:      data.UserID,
	}
}

func ManyMTOT(data []models.Ticket) []OperatorTicketResponse {
	resp := []OperatorTicketResponse{}
	for _, r := range data {
		resp = append(resp, ModelToOperatorTicket(r))
	}

	return resp
}

// type AdminTicketResponse struct {

// }

// func ModelToAdminTicket(data models.Ticket) AdminTicketResponse

func TicketToModel(tdto TicketDTO, user_id int) models.Ticket {
	now := time.Now()
	return models.Ticket{
		Title:       tdto.Title,
		Description: tdto.Description,
		Status:      constants.StatusNew,
		CreatedAt:   now,
		UserID:      user_id,
	}
}
