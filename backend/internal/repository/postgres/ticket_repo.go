package postgres

import (
	"context"
	"errors"
	"time"
	"tmaster/internal/apperrors"
	"tmaster/internal/constants"
	"tmaster/internal/repository/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repo) GetOwnTicketsRepo(ctx context.Context, user_id int, limit int, offset int) ([]models.Ticket, error) {
	var tickets []models.Ticket

	resp := r.db.Model(&models.Ticket{}).Where("user_id = ?", user_id).Limit(limit).Offset(offset).Find(&tickets)

	return tickets, resp.Error
}

func (r *Repo) CreateTicketRepo(ctx context.Context, ticket models.Ticket) error{
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *Repo) GetNewTickets(ctx context.Context, limit int) ([]models.Ticket, error) {
    var tickets []models.Ticket

    res := r.db.WithContext(ctx).
        Where("status = ?", constants.StatusNew).
        Order("created_at ASC, id ASC").
        Limit(limit).
        Find(&tickets)

    return tickets, res.Error
}

func (r *Repo) ClaimNextTicketRepo(ctx context.Context, operatorID uint) (*models.Ticket, error) {
	var ticket models.Ticket 
	err := r.db.WithContext(ctx).Transaction(func (tx *gorm.DB) error {
		result:=tx.Where("status = ?", constants.StatusNew).
		Order("created_at ASC, id ASC").
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options: "SKIP LOCKED",
		}).
		First(&ticket)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apperrors.TicketQueueIsEmpty
		}

		if result.Error != nil {
			return result.Error
		}

		now := time.Now()

		res := tx.Model(&models.Ticket{}).Updates(map[string]any{
			"status" : constants.StatusPending,
			"operator_id": operatorID,
            "claimed_at": now,
		})
		if res.Error != nil {
			return res.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ticket, nil
} 

func (r *Repo) CloseTicketRepo(ctx context.Context, ticket_id uint) error {
	return r.db.WithContext(ctx).Model(&models.Ticket{}).Update("status", constants.StatusClosed).Error
} 