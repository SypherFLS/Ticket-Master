package models

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Priority    string `gorm:"not null"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID"`
}
