package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Token     string
	ExpiresAt time.Time
}

func NewSession(userID uuid.UUID, ipAddress string, userAgent string, token string) *Session {
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Token:     token,
	}
}
