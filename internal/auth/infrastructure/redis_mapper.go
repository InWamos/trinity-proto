package infrastructure

import (
	"errors"
	"strconv"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
)

type RedisMapper struct {
}

// SessionToMap converts a domain.Session to a map for Redis storage
// Note: Token is not included as it will be used as the Redis key.
func (rm *RedisMapper) SessionToMap(session domain.Session) map[string]any {
	return map[string]any{
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
// Requires token parameter since it's stored as the Redis key, not in the hash.
func (rm *RedisMapper) MapToSession(data map[string]any, token string) (domain.Session, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return domain.Session{}, errors.New("id is not a string")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.Session{}, errors.New("failed to parse session id")
	}

	userIDStr, ok := data["user_id"].(string)
	if !ok {
		return domain.Session{}, errors.New("user_id is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.Session{}, errors.New("failed to parse user id")
	}

	userRoleStr, ok := data["user_role"].(string)
	if !ok {
		return domain.Session{}, errors.New("user_role is not a string")
	}
	userRole := domain.UserRole(userRoleStr)

	// Handle created_at
	var createdAt time.Time
	switch v := data["created_at"].(type) {
	case string:
		createdAtUnix, parseErr := strconv.ParseInt(v, 10, 64)
		if parseErr != nil {
			return domain.Session{}, errors.New("failed to parse created_at")
		}
		createdAt = time.Unix(createdAtUnix, 0).UTC()
	case int64:
		createdAt = time.Unix(v, 0).UTC()
	}

	// Handle expires_at
	var expiresAt time.Time
	switch v := data["expires_at"].(type) {
	case string:
		expiresAtUnix, parseErr := strconv.ParseInt(v, 10, 64)
		if parseErr != nil {
			return domain.Session{}, errors.New("failed to parse expires_at")
		}
		expiresAt = time.Unix(expiresAtUnix, 0).UTC()
	case int64:
		expiresAt = time.Unix(v, 0).UTC()
	}

	statusStr, ok := data["status"].(string)
	if !ok {
		return domain.Session{}, errors.New("status is not a string")
	}

	ipAddress, ok := data["ip_address"].(string)
	if !ok {
		return domain.Session{}, errors.New("ip_address is not a string")
	}

	userAgent, ok := data["user_agent"].(string)
	if !ok {
		return domain.Session{}, errors.New("user_agent is not a string")
	}

	return domain.Session{
		ID:        id,
		UserID:    userID,
		UserRole:  userRole,
		Status:    domain.SessionStatus(statusStr),
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}
