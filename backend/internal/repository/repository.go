package repository

import (
	"context"
	"tmaster/internal/repository/models"
)

type Repository interface {
	HealthCheck()
	CreateTicketRepo(ctx context.Context, ticket models.Ticket) error
}