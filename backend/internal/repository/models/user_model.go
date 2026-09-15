package models

type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleOperator UserRole = "operator"
)

type User struct {
	ID           uint     `gorm:"primaryKey"`
	Name         string   `gorm:"not null;unique"`
	Email        string   `gorm:"not null;unique"`
	Role         UserRole `gorm:"not null;default:user"`
	PasswordHash string

	Tickets []Ticket
}
