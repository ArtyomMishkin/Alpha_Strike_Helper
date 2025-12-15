package repository

import (
	"Alpha_Strike_Helper/internal/domain"

	"gorm.io/gorm"
)

// LanceRepository интерфейс для работы с Lance
type LanceRepository interface {
	Create(lance *domain.Lance) error
	GetByID(lanceID uint) (*domain.Lance, error)
	GetByUserID(userID uint) ([]*domain.Lance, error)
	Update(lance *domain.Lance) error
	Delete(lanceID uint) error
	AddMember(member *domain.LanceMember) error
	RemoveMember(lanceID, cardID uint) error
	GetMembers(lanceID uint) ([]domain.LanceMember, error)
}

// StarRepository интерфейс для работы с Star
type StarRepository interface {
	Create(star *domain.Star) error
	GetByID(starID uint) (*domain.Star, error)
	GetByUserID(userID uint) ([]*domain.Star, error)
	Update(star *domain.Star) error
	Delete(starID uint) error
	AddMember(member *domain.StarMember) error
	RemoveMember(starID, cardID uint) error
	GetMembers(starID uint) ([]domain.StarMember, error)
}

// LanceRepositoryImpl реализация LanceRepository
type LanceRepositoryImpl struct {
	db *gorm.DB
}

// StarRepositoryImpl реализация StarRepository
type StarRepositoryImpl struct {
	db *gorm.DB
}

// NewLanceRepository создаёт новый LanceRepository
func NewLanceRepository(db *gorm.DB) LanceRepository {
	return &LanceRepositoryImpl{db: db}
}

// NewStarRepository создаёт новый StarRepository
func NewStarRepository(db *gorm.DB) StarRepository {
	return &StarRepositoryImpl{db: db}
}

// ========== LANCE REPOSITORY ==========

func (r *LanceRepositoryImpl) Create(lance *domain.Lance) error {
	return r.db.Create(lance).Error
}

func (r *LanceRepositoryImpl) GetByID(lanceID uint) (*domain.Lance, error) {
	var lance domain.Lance
	err := r.db.Preload("Members.Card").First(&lance, lanceID).Error
	return &lance, err
}

func (r *LanceRepositoryImpl) GetByUserID(userID uint) ([]*domain.Lance, error) {
	var lances []*domain.Lance
	err := r.db.Preload("Members.Card").Where("user_id = ?", userID).Find(&lances).Error
	return lances, err
}

func (r *LanceRepositoryImpl) Update(lance *domain.Lance) error {
	return r.db.Save(lance).Error
}

func (r *LanceRepositoryImpl) Delete(lanceID uint) error {
	return r.db.Delete(&domain.Lance{}, lanceID).Error
}

func (r *LanceRepositoryImpl) AddMember(member *domain.LanceMember) error {
	return r.db.Create(member).Error
}

func (r *LanceRepositoryImpl) RemoveMember(lanceID, cardID uint) error {
	return r.db.Where("lance_id = ? AND card_id = ?", lanceID, cardID).Delete(&domain.LanceMember{}).Error
}

func (r *LanceRepositoryImpl) GetMembers(lanceID uint) ([]domain.LanceMember, error) {
	var members []domain.LanceMember
	err := r.db.Preload("Card").Where("lance_id = ?", lanceID).Order("position ASC").Find(&members).Error
	return members, err
}

// ========== STAR REPOSITORY ==========

func (r *StarRepositoryImpl) Create(star *domain.Star) error {
	return r.db.Create(star).Error
}

func (r *StarRepositoryImpl) GetByID(starID uint) (*domain.Star, error) {
	var star domain.Star
	err := r.db.Preload("Members.Card").First(&star, starID).Error
	return &star, err
}

func (r *StarRepositoryImpl) GetByUserID(userID uint) ([]*domain.Star, error) {
	var stars []*domain.Star
	err := r.db.Preload("Members.Card").Where("user_id = ?", userID).Find(&stars).Error
	return stars, err
}

func (r *StarRepositoryImpl) Update(star *domain.Star) error {
	return r.db.Save(star).Error
}

func (r *StarRepositoryImpl) Delete(starID uint) error {
	return r.db.Delete(&domain.Star{}, starID).Error
}

func (r *StarRepositoryImpl) AddMember(member *domain.StarMember) error {
	return r.db.Create(member).Error
}

func (r *StarRepositoryImpl) RemoveMember(starID, cardID uint) error {
	return r.db.Where("star_id = ? AND card_id = ?", starID, cardID).Delete(&domain.StarMember{}).Error
}

func (r *StarRepositoryImpl) GetMembers(starID uint) ([]domain.StarMember, error) {
	var members []domain.StarMember
	err := r.db.Preload("Card").Where("star_id = ?", starID).Order("position ASC").Find(&members).Error
	return members, err
}
