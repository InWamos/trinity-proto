package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	Active  SessionStatus = "active"
	Revoked SessionStatus = "revoked"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    SessionStatus
	IPAddress string
	UserAgent string
	Token     string
	ExpiresAt time.Time
}

func NewSession(userID uuid.UUID, ipAddress string, userAgent string, token string) *Session {
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    Active,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Token:     token,
	}
}
