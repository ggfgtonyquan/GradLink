package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
	"gradlink/internal/db"
)

const (
	SessionCookieName = "gradlink_session"
	SessionTTL        = 30 * 24 * time.Hour
)

type SessionStore struct {
	DB  *gorm.DB
	Now func() time.Time
}

func (s *SessionStore) Create(userID uint) (string, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}

	session := db.Session{
		UserID:    userID,
		TokenHash: HashToken(token),
		ExpiresAt: s.now().Add(SessionTTL),
	}
	if err := s.DB.Create(&session).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (s *SessionStore) FindActive(token string) (*db.Session, error) {
	var session db.Session
	err := s.DB.Preload("User").
		Where("token_hash = ? AND expires_at > ?", HashToken(token), s.now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) Delete(token string) error {
	return s.DB.Where("token_hash = ?", HashToken(token)).Delete(&db.Session{}).Error
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *SessionStore) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}

	return time.Now()
}

func randomToken(byteCount int) (string, error) {
	token := make([]byte, byteCount)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
