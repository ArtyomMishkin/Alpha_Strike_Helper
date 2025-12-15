package domain

import "time"

// User модель пользователя
type User struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Username      string        `gorm:"unique;not null;index" json:"username"`
	Email         string        `gorm:"unique;not null;index" json:"email"`
	PasswordHash  string        `gorm:"not null" json:"-"` // Не отправляем пароль в JSON
	Collections   []Collection  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"collections,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	LastLogin     *time.Time    `json:"last_login,omitempty"`
}

// TableName указывает имя таблицы для GORM
func (User) TableName() string {
	return "users"
}
