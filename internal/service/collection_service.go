package service

import (
	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/repository"
	"fmt"
)

type CollectionService struct {
	collectionRepo *repository.CollectionRepository
	cardRepo       *repository.CardRepository
}

func NewCollectionService(collectionRepo *repository.CollectionRepository, cardRepo *repository.CardRepository) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
		cardRepo:       cardRepo,
	}
}

// ListCollections получает список коллекций пользователя
func (s *CollectionService) ListCollections(userID uint) ([]domain.Collection, error) {
	if s.collectionRepo == nil {
		return nil, fmt.Errorf("collection repository is nil")
	}

	collections, err := s.collectionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	return collections, nil
}

// GetCollection получает коллекцию по ID
func (s *CollectionService) GetCollection(id uint) (*domain.Collection, error) {
	if s.collectionRepo == nil {
		return nil, fmt.Errorf("collection repository is nil")
	}

	collection, err := s.collectionRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return collection, nil
}

// CreateCollection создаёт новую коллекцию
func (s *CollectionService) CreateCollection(collection *domain.Collection) error {
	if s.collectionRepo == nil {
		return fmt.Errorf("collection repository is nil")
	}

	if err := s.collectionRepo.Create(collection); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// UpdateCollection обновляет коллекцию
func (s *CollectionService) UpdateCollection(collection *domain.Collection) error {
	if s.collectionRepo == nil {
		return fmt.Errorf("collection repository is nil")
	}

	if err := s.collectionRepo.Update(collection); err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	return nil
}

// DeleteCollection удаляет коллекцию
func (s *CollectionService) DeleteCollection(id uint) error {
	if s.collectionRepo == nil {
		return fmt.Errorf("collection repository is nil")
	}

	if err := s.collectionRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// AddCardToCollection добавляет карточку в коллекцию
func (s *CollectionService) AddCardToCollection(collectionID, cardID uint) error {
	if s.collectionRepo == nil {
		return fmt.Errorf("collection repository is nil")
	}

	// Проверяем что карточка существует
	if s.cardRepo != nil {
		_, err := s.cardRepo.Get(cardID)
		if err != nil {
			return fmt.Errorf("card not found: %w", err)
		}
	}

	return s.collectionRepo.AddCard(collectionID, cardID)
}

// RemoveCardFromCollection удаляет карточку из коллекции
func (s *CollectionService) RemoveCardFromCollection(collectionID, cardID uint) error {
	if s.collectionRepo == nil {
		return fmt.Errorf("collection repository is nil")
	}

	return s.collectionRepo.RemoveCard(collectionID, cardID)
}

// GetCollectionCards получает все карточки коллекции
func (s *CollectionService) GetCollectionCards(collectionID uint) ([]domain.Card, error) {
	if s.collectionRepo == nil {
		return nil, fmt.Errorf("collection repository is nil")
	}

	return s.collectionRepo.GetCards(collectionID)
}

// ListCollectionsWithPagination получает коллекции с пагинацией
func (s *CollectionService) ListCollectionsWithPagination(page, pageSize int) ([]domain.Collection, int64, error) {
	if s.collectionRepo == nil {
		return nil, 0, fmt.Errorf("collection repository is nil")
	}

	collections, total, err := s.collectionRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list collections: %w", err)
	}

	return collections, total, nil
}
