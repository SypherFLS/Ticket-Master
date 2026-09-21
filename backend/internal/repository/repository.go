package repository

import (
	"context"
	"tmaster/internal/constants"
	"tmaster/internal/dto"
	"tmaster/internal/repository/models"
)

type Repository interface {
	RegisterRepo(ctx context.Context, user models.User) error
	GetLoginDataRepo(ctx context.Context, email string) (models.LoginResult, error)
	SetRoleRepo(ctx context.Context,email string, role constants.UserRole) error
	
	GetOwnTicketsRepo(ctx context.Context, user_id int, pagData dto.PaginationData) ([]models.Ticket, error)
	GetNewTickets(ctx context.Context, limit int) ([]models.Ticket, error)
	ClaimNextTicketRepo(ctx context.Context, operatorID int) (*models.Ticket, error)
	CreateTicketRepo(ctx context.Context, ticket models.Ticket) error
	CloseTicketRepo(ctx context.Context, ticket_id uint, operator_id int) error
}