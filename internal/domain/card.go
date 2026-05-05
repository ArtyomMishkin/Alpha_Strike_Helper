package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Card — основная модель для работы с БД и JSON
type Card struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ModelNumber string `gorm:"uniqueIndex" json:"model_number"`
	Name        string `gorm:"not null;index" json:"name"`
	// UnitType is the catalog category (e.g. BATTLEMECH / INFANTRY / COMBAT-VEHICLE).
	// It allows storing multiple unit kinds in the same `cards` table.
	UnitType               string                      `gorm:"index" json:"unit_type"`
	Type                   string                      `json:"type"` // Medium, Heavy, Assault
	Size                   string                      `json:"size"`
	Move                   string                      `json:"move"`
	TMM                    int                         `json:"tmm"`
	PointValue             int                         `json:"point_value"`
	Armor                  int                         `json:"armor"`
	Structure              int                         `json:"structure"`
	DamageShort            string                      `json:"damage_short"`
	DamageMedium           string                      `json:"damage_medium"`
	DamageLong             string                      `json:"damage_long"`
	Overheat               int                         `json:"overheat"`
	Abilities              datatypes.JSONSlice[string] `gorm:"type:jsonb;default:'[]'" json:"abilities"`
	Tonnage                string                      `json:"tonnage"`
	TechBase               string                      `json:"tech_base"`
	Role                   string                      `json:"role"`
	Source                 string                      `json:"source"`
	Faction                string                      `json:"faction"`
	FactionGroup           string                      `json:"faction_group"`
	Era                    string                      `json:"era"`
	AvailableFactions      datatypes.JSONSlice[string] `gorm:"type:jsonb;default:'[]'" json:"available_factions"`
	AvailableFactionGroups datatypes.JSONSlice[string] `gorm:"type:jsonb;default:'[]'" json:"available_faction_groups"`
	AvailableEras          datatypes.JSONSlice[string] `gorm:"type:jsonb;default:'[]'" json:"available_eras"`
	FactionEraAvailability datatypes.JSON              `gorm:"type:jsonb;default:'{}'" json:"faction_era_availability"`
	ImageURL               string                      `json:"image_url"`
	CardURL                string                      `json:"card_url"`
	CreatedAt              time.Time                   `json:"created_at"`
	UpdatedAt              time.Time                   `json:"updated_at"`
}

// TableName указывает имя таблицы для GORM
func (Card) TableName() string {
	return "cards"
}
