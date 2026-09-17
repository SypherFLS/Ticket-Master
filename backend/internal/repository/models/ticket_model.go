package models

type TicketStatus string

const (
	StatusNew     TicketStatus = "new"
	StatusPending TicketStatus = "pending"
	StatusClose   TicketStatus = "close"
)

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string	
	Priority    string `gorm:"not null"`
	Status      string `gorm:"not null,default:new"`
	UserID      uint   `gorm:"not null;index"`
	User        User   `gorm:"foreignKey:UserID"`
}
