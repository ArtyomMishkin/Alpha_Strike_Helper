package service

import (
	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/repository"
	"Alpha_Strike_Helper/pkg/utils"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService сервис аутентификации
type AuthService struct {
	users   repository.UserRepository
	jwtUtil *utils.JWTService
}

func NewAuthService(users repository.UserRepository, jwtUtil *utils.JWTService) *AuthService {
	return &AuthService{
		users:   users,
		jwtUtil: jwtUtil,
	}
}

// RegisterRequest для регистрации
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest для входа
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse для ответа
type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// Register регистрация нового пользователя
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	if req == nil {
		return nil, errors.New("register request cannot be nil")
	}

	// Проверяем, не существует ли уже пользователь
	existing, err := s.users.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existing, err = s.users.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	// Генерируем JWT
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login вход пользователя
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	if req == nil {
		return nil, errors.New("login request cannot be nil")
	}

	// Ищем пользователя по username
	user, err := s.users.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Обновляем время последнего входа
	s.users.UpdateLastLogin(user.ID)

	// Генерируем JWT
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ValidateToken проверяет JWT токен
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	return s.jwtUtil.ValidateToken(tokenString)
}
