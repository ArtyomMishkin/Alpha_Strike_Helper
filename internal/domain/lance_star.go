package domain

import "time"

// Lance — отряд Внутренней Сферы (4 меха)
type Lance struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Faction     string         `json:"faction"` // Inner Sphere фракция
	Members     []LanceMember  `gorm:"foreignKey:LanceID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
	User        *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// LanceMember — мех в отряде (Lance)
type LanceMember struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	LanceID   uint       `gorm:"not null;index:idx_lance_member" json:"lance_id"`
	CardID    uint       `gorm:"not null;index:idx_lance_member" json:"card_id"`
	Position  int        `json:"position"` // 1-4 позиция в отряде
	Lance     *Lance     `gorm:"foreignKey:LanceID;constraint:OnDelete:CASCADE" json:"lance,omitempty"`
	Card      *Card      `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"card,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Star — отряд Кланов (5 мехов)
type Star struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	UserID      uint          `gorm:"not null;index" json:"user_id"`
	Name        string        `gorm:"not null" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	Faction     string        `json:"faction"` // Clan фракция
	Members     []StarMember  `gorm:"foreignKey:StarID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
	User        *User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// StarMember — мех в отряде (Star)
type StarMember struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	StarID    uint       `gorm:"not null;index:idx_star_member" json:"star_id"`
	CardID    uint       `gorm:"not null;index:idx_star_member" json:"card_id"`
	Position  int        `json:"position"` // 1-5 позиция в отряде
	Star      *Star      `gorm:"foreignKey:StarID;constraint:OnDelete:CASCADE" json:"star,omitempty"`
	Card      *Card      `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"card,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// LanceStats — расчётные свойства Lance
type LanceStats struct {
	LanceID       uint    `json:"lance_id"`
	LanceName     string  `json:"lance_name"`
	UnitCount     int     `json:"unit_count"`      // Кол-во мехов в отряде (макс 4)
	TotalTonnage  int     `json:"total_tonnage"`   // Сумма tonnage
	AvgTonnage    float64 `json:"avg_tonnage"`     // Среднее tonnage
	TotalPV       int     `json:"total_pv"`        // Сумма PointValue
	AvgPV         float64 `json:"avg_pv"`          // Среднее PointValue
	AvgTMM        float64 `json:"avg_tmm"`         // Среднее TMM
	TotalDamage   int     `json:"total_damage"`    // Среднее damage (short range)
	TechBases     []string `json:"tech_bases"`     // Какие tech bases представлены
	Roles         []string `json:"roles"`          // Какие роли представлены
}

// StarStats — расчётные свойства Star
type StarStats struct {
	StarID        uint    `json:"star_id"`
	StarName      string  `json:"star_name"`
	UnitCount     int     `json:"unit_count"`      // Кол-во мехов в отряде (макс 5)
	TotalTonnage  int     `json:"total_tonnage"`   // Сумма tonnage
	AvgTonnage    float64 `json:"avg_tonnage"`     // Среднее tonnage
	TotalPV       int     `json:"total_pv"`        // Сумма PointValue
	AvgPV         float64 `json:"avg_pv"`          // Среднее PointValue
	AvgTMM        float64 `json:"avg_tmm"`         // Среднее TMM
	TotalDamage   int     `json:"total_damage"`    // Среднее damage (short range)
	TechBases     []string `json:"tech_bases"`     // Какие tech bases представлены
	Roles         []string `json:"roles"`          // Какие роли представлены
}

// TableName указывает имя таблицы для GORM
func (Lance) TableName() string {
	return "lances"
}

func (LanceMember) TableName() string {
	return "lance_members"
}

func (Star) TableName() string {
	return "stars"
}

func (StarMember) TableName() string {
	return "star_members"
}
