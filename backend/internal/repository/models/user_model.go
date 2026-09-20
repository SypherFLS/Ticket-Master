package models

import "tmaster/internal/constants"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null;unique"`
	PasswordHash string
	Email        string   `gorm:"not null;unique"`
	Role         constants.UserRole `gorm:"not null;default:user"`

	Tickets []Ticket
}

type LoginResult struct {
	ID           int
	UserRole     string
	PasswordHash string
}
