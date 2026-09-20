package models

import "time"

type TicketStatus string

const (
	StatusNew     TicketStatus = "new"
	StatusPending TicketStatus = "pending"
	StatusClosed   TicketStatus = "closed"
)

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	OperatorID  uint
	Description string
	Status      TicketStatus `gorm:"not null,default:new"`
	CreatedAt   time.Time
	ClaimedAt   time.Time
	ResolvedAt  time.Time
	UserID      uint `gorm:"not null;index"`
	User        User `gorm:"foreignKey:UserID"`
}
