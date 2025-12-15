package repository

import (
	"Alpha_Strike_Helper/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

// CollectionRepository это структура (не интерфейс!) для работы с коллекциями
type CollectionRepository struct {
	db *gorm.DB
}

// NewCollectionRepository создаёт новый репозиторий коллекций
func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create создаёт новую коллекцию
func (r *CollectionRepository) Create(collection *domain.Collection) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Create(collection).Error; err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// Get получает коллекцию по ID
func (r *CollectionRepository) Get(id uint) (*domain.Collection, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var collection domain.Collection
	if err := r.db.First(&collection, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &collection, nil
}

// GetByUserID получает все коллекции пользователя
func (r *CollectionRepository) GetByUserID(userID uint) ([]domain.Collection, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var collections []domain.Collection
	if err := r.db.Where("user_id = ?", userID).Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	return collections, nil
}

// Update обновляет коллекцию
func (r *CollectionRepository) Update(collection *domain.Collection) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Save(collection).Error; err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	return nil
}

// Delete удаляет коллекцию
func (r *CollectionRepository) Delete(id uint) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := r.db.Delete(&domain.Collection{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// AddCard добавляет карточку в коллекцию
func (r *CollectionRepository) AddCard(collectionID, cardID uint) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Проверяем что коллекция существует
	var collection domain.Collection
	if err := r.db.First(&collection, collectionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("collection not found")
		}
		return fmt.Errorf("failed to check collection: %w", err)
	}

	// Проверяем что карточка существует
	var card domain.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("card not found")
		}
		return fmt.Errorf("failed to check card: %w", err)
	}

	// Добавляем карточку в коллекцию
	// Предполагаем что есть many-to-many отношение
	if err := r.db.Model(&collection).Association("Cards").Append(&card); err != nil {
		return fmt.Errorf("failed to add card to collection: %w", err)
	}

	return nil
}

// RemoveCard удаляет карточку из коллекции
func (r *CollectionRepository) RemoveCard(collectionID, cardID uint) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	var collection domain.Collection
	if err := r.db.First(&collection, collectionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("collection not found")
		}
		return fmt.Errorf("failed to check collection: %w", err)
	}

	var card domain.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("card not found")
		}
		return fmt.Errorf("failed to check card: %w", err)
	}

	if err := r.db.Model(&collection).Association("Cards").Delete(&card); err != nil {
		return fmt.Errorf("failed to remove card from collection: %w", err)
	}

	return nil
}

// GetCards получает все карточки коллекции
func (r *CollectionRepository) GetCards(collectionID uint) ([]domain.Card, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var collection domain.Collection
	if err := r.db.First(&collection, collectionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	var cards []domain.Card
	if err := r.db.Model(&collection).Association("Cards").Find(&cards); err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	return cards, nil
}

// List получает все коллекции с пагинацией
func (r *CollectionRepository) List(page, pageSize int) ([]domain.Collection, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database is not initialized")
	}

	var collections []domain.Collection
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&domain.Collection{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count collections: %w", err)
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&collections).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list collections: %w", err)
	}

	return collections, total, nil
}
