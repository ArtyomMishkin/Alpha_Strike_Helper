package repository

import (
	"Alpha_Strike_Helper/internal/domain"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

// List получает список карточек с фильтрами и пагинацией
func (r *CardRepository) List(filters map[string]interface{}, page, pageSize int) ([]domain.Card, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	var total int64
	query := r.db

	// Применяем фильтры
	if role, ok := filters["role"].(string); ok && role != "" {
		query = query.Where("role = ?", role)
	}

	if size, ok := filters["size"].(string); ok && size != "" {
		query = query.Where("size = ?", size)
	}

	if faction, ok := filters["faction"].(string); ok && faction != "" {
		query = query.Where("faction = ?", faction)
	}

	if typeVal, ok := filters["type"].(string); ok && typeVal != "" {
		query = query.Where("type = ?", typeVal)
	}

	if techBase, ok := filters["tech_base"].(string); ok && techBase != "" {
		query = query.Where("tech_base = ?", techBase)
	}

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%")
	}

	// PV фильтры (точки за valor)
	if pvMin, ok := filters["pv_min"].(string); ok && pvMin != "" {
		if val, err := strconv.Atoi(pvMin); err == nil {
			query = query.Where("point_value >= ?", val)
		}
	}

	if pvMax, ok := filters["pv_max"].(string); ok && pvMax != "" {
		if val, err := strconv.Atoi(pvMax); err == nil {
			query = query.Where("point_value <= ?", val)
		}
	}

	// Подсчитываем total с фильтрами
	if err := query.Model(&domain.Card{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&cards).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list cards: %w", err)
	}

	return cards, total, nil
}

// Get получает карточку по ID
func (r *CardRepository) Get(id uint) (*domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var card domain.Card
	if err := r.db.First(&card, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("card not found")
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return &card, nil
}

// Create создаёт новую карточку
func (r *CardRepository) Create(card *domain.Card) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Create(card).Error; err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	return nil
}

// Update обновляет карточку
func (r *CardRepository) Update(card *domain.Card) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Save(card).Error; err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	return nil
}

// Delete удаляет карточку
func (r *CardRepository) Delete(id uint) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Delete(&domain.Card{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}

// Search ищет карточки по названию
func (r *CardRepository) Search(query string) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("LOWER(name) LIKE LOWER(?) OR LOWER(model_number) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%").
		Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}

	return cards, nil
}

// GetByRole получает карточки по роли
func (r *CardRepository) GetByRole(role string) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("role = ?", role).Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get cards by role: %w", err)
	}

	return cards, nil
}

// GetByFaction получает карточки по фракции
func (r *CardRepository) GetByFaction(faction string) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("faction = ?", faction).Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get cards by faction: %w", err)
	}

	return cards, nil
}

// GetBySize получает карточки по размеру
func (r *CardRepository) GetBySize(size string) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("size = ?", size).Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get cards by size: %w", err)
	}

	return cards, nil
}

// GetByType получает карточки по типу
func (r *CardRepository) GetByType(typeVal string) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("type = ?", typeVal).Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get cards by type: %w", err)
	}

	return cards, nil
}

// GetAll получает все карточки
func (r *CardRepository) GetAll() ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get all cards: %w", err)
	}

	return cards, nil
}

// Count подсчитывает количество карточек
func (r *CardRepository) Count() (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("database is not initialized")
	}

	var count int64
	if err := r.db.Model(&domain.Card{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count cards: %w", err)
	}

	return count, nil
}

// GetWithin получает карточки в диапазоне PV
func (r *CardRepository) GetWithinPVRange(minPV, maxPV int) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var cards []domain.Card
	if err := r.db.Where("point_value >= ? AND point_value <= ?", minPV, maxPV).
		Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("failed to get cards within PV range: %w", err)
	}

	return cards, nil
}
