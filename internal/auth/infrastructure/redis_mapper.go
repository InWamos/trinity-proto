package infrastructure

import (
	"strconv"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
)

type RedisMapper struct {
}

// SessionToMap converts a domain.Session to a map for Redis storage
// Note: Token is not included as it will be used as the Redis key
func (rm *RedisMapper) SessionToMap(session domain.Session) map[string]interface{} {
	return map[string]interface{}{
		"id":         session.ID.String(),
		"user_id":    session.UserID.String(),
		"user_role":  string(session.UserRole),
		"status":     string(session.Status),
		"ip_address": session.IPAddress,
		"user_agent": session.UserAgent,
		"created_at": session.CreatedAt.Unix(),
		"expires_at": session.ExpiresAt.Unix(),
	}
}

// MapToSession converts a map from Redis to a domain.Session
// Requires token parameter since it's stored as the Redis key, not in the hash
func (rm *RedisMapper) MapToSession(data map[string]interface{}, token string) (domain.Session, error) {
	id, err := uuid.Parse(data["id"].(string))
	if err != nil {
		return domain.Session{}, err
	}

	userID, err := uuid.Parse(data["user_id"].(string))
	if err != nil {
		return domain.Session{}, err
	}

	userRole := domain.UserRole(data["user_role"].(string))

	// Handle created_at
	var createdAt time.Time
	switch v := data["created_at"].(type) {
	case string:
		createdAtUnix, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return domain.Session{}, err
		}
		createdAt = time.Unix(createdAtUnix, 0).UTC()
	case int64:
		createdAt = time.Unix(v, 0).UTC()
	}

	// Handle expires_at
	var expiresAt time.Time
	switch v := data["expires_at"].(type) {
	case string:
		expiresAtUnix, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return domain.Session{}, err
		}
		expiresAt = time.Unix(expiresAtUnix, 0).UTC()
	case int64:
		expiresAt = time.Unix(v, 0).UTC()
	}

	return domain.Session{
		ID:        id,
		UserID:    userID,
		UserRole:  userRole,
		Status:    domain.SessionStatus(data["status"].(string)),
		IPAddress: data["ip_address"].(string),
		UserAgent: data["user_agent"].(string),
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}
