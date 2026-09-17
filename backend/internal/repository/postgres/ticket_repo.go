package postgres

import (
	"context"
	"tmaster/internal/repository/models"
)

func (r *Repo) GetTicketsRepo(ctx context.Context, user_id int, limit int, offset int) ([]models.Ticket, error) {
	var tickets []models.Ticket

	resp := r.db.Model(&models.Ticket{}).Where("user_id = ?", user_id).Limit(limit).Offset(offset).Find(&tickets)

	return tickets, resp.Error
}

func (r *Repo) CreateTicketRepo(ctx context.Context, ticket models.Ticket) error{
	return r.db.WithContext(ctx).Create(ticket).Error
}
