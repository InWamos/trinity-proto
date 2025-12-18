package domain

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type SessionStatus string
type UserRole string

const (
	Active  SessionStatus = "active"
	Revoked SessionStatus = "revoked"
	Admin   UserRole      = "admin"
	User    UserRole      = "user"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserRole  UserRole
	Status    SessionStatus
	IPAddress string
	UserAgent string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewSession(
	userID uuid.UUID,
	userRole UserRole,
	ipAddress string,
	userAgent string,
	token string,
	duration time.Duration,
) (*Session, error) {
	token, err := generateToken(32)
	if err != nil {
		return &Session{}, err
	}
	createdAt := time.Now().UTC()
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		UserRole:  userRole,
		Status:    Active,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(duration),
	}, nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Use URL-safe base64 encoding (no padding)
	return base64.URLEncoding.EncodeToString(bytes), nil
}
