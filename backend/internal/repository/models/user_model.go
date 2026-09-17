package models

type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleOperator UserRole = "operator"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null;unique"`
	PasswordHash string
	Email        string   `gorm:"not null;unique"`
	Role         UserRole `gorm:"not null;default:user"`

	Tickets []Ticket
}

type LoginResult struct {
	ID           int
	PasswordHash string
}
