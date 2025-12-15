package domain

import "time"

// Collection модель коллекции карточек пользователя
type Collection struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	UserID      uint              `gorm:"not null;index" json:"user_id"`
	Name        string            `gorm:"not null" json:"name"`
	Description string            `gorm:"type:text" json:"description"`
	Cards       []CollectionCard  `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"cards,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	User        *User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// CollectionCard связь Many-to-Many между Collection и Card
type CollectionCard struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	CollectionID uint        `gorm:"not null;index:idx_collection_card" json:"collection_id"`
	CardID       uint        `gorm:"not null;index:idx_collection_card" json:"card_id"`
	Quantity     int         `gorm:"default:1" json:"quantity"`
	Collection   *Collection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"collection,omitempty"`
	Card         *Card       `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"card,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

// TableName указывает имя таблицы для GORM
func (Collection) TableName() string {
	return "collections"
}

func (CollectionCard) TableName() string {
	return "collection_cards"
}
