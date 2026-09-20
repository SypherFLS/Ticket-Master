package models

import "time"

type TicketStatus string

const (
	StatusNew     TicketStatus = "new"
	StatusPending TicketStatus = "pending"
	StatusClose   TicketStatus = "close"
)

type Ticket struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"not null"`
	WorkingUserID uint
	Description   string
	Status        TicketStatus `gorm:"not null,default:new"`
	CreatedAt     time.Time
	ClaimedAT     time.Time
	ResolvedAT    time.Time
	UserID        uint `gorm:"not null;index"`
	User          User `gorm:"foreignKey:UserID"`
}
