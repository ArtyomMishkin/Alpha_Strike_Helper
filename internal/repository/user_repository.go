package repository

import (
	"Alpha_Strike_Helper/internal/domain"
	"errors"
	"gorm.io/gorm"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	UpdateLastLogin(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.Username == "" || user.Email == "" {
		return errors.New("username and email are required")
	}
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == 0 {
		return errors.New("user ID is required for update")
	}
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("last_login", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
