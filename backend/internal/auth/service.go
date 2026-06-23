package auth

import (
	"errors"
	"net/mail"
	"strings"

	"gorm.io/gorm"
	"gradlink/internal/db"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDisabledUser       = errors.New("disabled user")
	ErrValidation         = errors.New("validation failed")
)

type Service struct {
	DB       *gorm.DB
	Sessions *SessionStore
}

type RegisterInput struct {
	Email    string
	Password string
	Nickname string
}

func (s *Service) Register(input RegisterInput) (*db.User, string, error) {
	email, nickname, err := normalizeRegistration(input)
	if err != nil {
		return nil, "", err
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, "", err
	}

	user := db.User{
		Email:        email,
		PasswordHash: hash,
		Nickname:     nickname,
		Role:         db.UserRoleUser,
		Status:       db.UserStatusActive,
	}
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, "", ErrValidation
	}

	token, err := s.Sessions.Create(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *Service) Login(email, password string) (*db.User, string, error) {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	var user db.User
	err = s.DB.Where("email = ?", normalized).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", err
	}
	if !CheckPassword(user.PasswordHash, password) {
		return nil, "", ErrInvalidCredentials
	}
	if user.Status != db.UserStatusActive {
		return nil, "", ErrDisabledUser
	}

	token, err := s.Sessions.Create(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func normalizeRegistration(input RegisterInput) (string, string, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return "", "", ErrValidation
	}
	nickname := strings.TrimSpace(input.Nickname)
	if nickname == "" || len(nickname) > 80 || len(input.Password) < 8 {
		return "", "", ErrValidation
	}

	return email, nickname, nil
}

func normalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", ErrValidation
	}
	parsed, err := mail.ParseAddress(normalized)
	if err != nil || parsed.Address != normalized {
		return "", ErrValidation
	}

	return normalized, nil
}
