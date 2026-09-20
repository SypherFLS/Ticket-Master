package models

import (
	"time"
	"tmaster/internal/constants"
)



type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	OperatorID  uint
	Description string
	Status      constants.TicketStatus `gorm:"not null,default:new"`
	CreatedAt   time.Time
	ClaimedAt   time.Time
	ResolvedAt  time.Time
	UserID      uint `gorm:"not null;index"`
	User        User `gorm:"foreignKey:UserID"`
}
