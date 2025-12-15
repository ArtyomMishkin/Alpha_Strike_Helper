package service

import (
	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/repository"
	"fmt"
)

type CardService struct {
	cardRepo *repository.CardRepository
}

func NewCardService(repo *repository.CardRepository) *CardService {
	return &CardService{cardRepo: repo}
}

// ListCards получает список карточек с фильтрами
func (s *CardService) ListCards(filters map[string]interface{}, page, pageSize int) ([]domain.Card, int64, error) {
	if s.cardRepo == nil {
		return nil, 0, fmt.Errorf("card repository is nil")
	}

	cards, total, err := s.cardRepo.List(filters, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cards: %w", err)
	}

	return cards, total, nil
}

// GetCard получает карточку по ID
func (s *CardService) GetCard(id uint) (*domain.Card, error) {
	if s.cardRepo == nil {
		return nil, fmt.Errorf("card repository is nil")
	}

	card, err := s.cardRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return card, nil
}

// CreateCard создаёт новую карточку
func (s *CardService) CreateCard(card *domain.Card) error {
	if s.cardRepo == nil {
		return fmt.Errorf("card repository is nil")
	}

	if err := s.cardRepo.Create(card); err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	return nil
}

// UpdateCard обновляет карточку
func (s *CardService) UpdateCard(card *domain.Card) error {
	if s.cardRepo == nil {
		return fmt.Errorf("card repository is nil")
	}

	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	return nil
}

// DeleteCard удаляет карточку
func (s *CardService) DeleteCard(id uint) error {
	if s.cardRepo == nil {
		return fmt.Errorf("card repository is nil")
	}

	if err := s.cardRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}

// SearchCards ищет карточки по названию
func (s *CardService) SearchCards(query string) ([]domain.Card, error) {
	if s.cardRepo == nil {
		return nil, fmt.Errorf("card repository is nil")
	}

	cards, err := s.cardRepo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}

	return cards, nil
}

// GetByRole получает карточки по роли
func (s *CardService) GetByRole(role string) ([]domain.Card, error) {
	if s.cardRepo == nil {
		return nil, fmt.Errorf("card repository is nil")
	}

	return s.cardRepo.GetByRole(role)
}

// GetByFaction получает карточки по фракции
func (s *CardService) GetByFaction(faction string) ([]domain.Card, error) {
	if s.cardRepo == nil {
		return nil, fmt.Errorf("card repository is nil")
	}

	return s.cardRepo.GetByFaction(faction)
}

// GetBySize получает карточки по размеру
func (s *CardService) GetBySize(size string) ([]domain.Card, error) {
	if s.cardRepo == nil {
		return nil, fmt.Errorf("card repository is nil")
	}

	return s.cardRepo.GetBySize(size)
}
