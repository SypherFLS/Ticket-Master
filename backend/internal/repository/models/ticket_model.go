package models

import (
	"time"
	"tmaster/internal/constants"
)

type Ticket struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	OperatorID  int
	Description string
	Status      constants.TicketStatus `gorm:"not null,default:new"`
	CreatedAt   time.Time
	ClaimedAt   time.Time
	ResolvedAt  time.Time
	UserID      int  `gorm:"not null;index"`
	User        User `gorm:"foreignKey:UserID"`
}
