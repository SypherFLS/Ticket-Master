package models

type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleOperator UserRole = "operator"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uint     `gorm:"primaryKey"`
	Name         string   `gorm:"not null;unique"`
	Email        string   `gorm:"not null;unique,index"`
	Role         UserRole `gorm:"not null;default:user"`
	PasswordHash string

	Tickets []Ticket
}

type LoginResult struct {
	ID           int
	PasswordHash string
}
