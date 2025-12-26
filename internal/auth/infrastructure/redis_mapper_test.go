package infrastructure_test

import (
	"testing"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionToMap(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}
	sessionID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(24 * time.Hour)

	session := domain.Session{
		ID:        sessionID,
		UserID:    userID,
		UserRole:  domain.User,
		Status:    domain.Active,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		Token:     "test-token-123",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	result := mapper.SessionToMap(session)

	assert.Equal(t, sessionID.String(), result["id"])
	assert.Equal(t, userID.String(), result["user_id"])
	assert.Equal(t, string(domain.User), result["user_role"])
	assert.Equal(t, string(domain.Active), result["status"])
	assert.Equal(t, "192.168.1.1", result["ip_address"])
	assert.Equal(t, "Mozilla/5.0", result["user_agent"])
	assert.Nil(t, result["token"]) // Token should not be in the map
	assert.Equal(t, createdAt.Unix(), result["created_at"])
	assert.Equal(t, expiresAt.Unix(), result["expires_at"])
}

func TestSessionToMapWithRevokedStatus(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}
	sessionID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(24 * time.Hour)

	session := domain.Session{
		ID:        sessionID,
		UserID:    userID,
		UserRole:  domain.Admin,
		Status:    domain.Revoked,
		IPAddress: "10.0.0.1",
		UserAgent: "Chrome/120",
		Token:     "revoked-token",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	result := mapper.SessionToMap(session)

	assert.Equal(t, string(domain.Admin), result["user_role"])
	assert.Equal(t, string(domain.Revoked), result["status"])
}

func TestMapToSession(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}
	sessionID := uuid.New()
	userID := uuid.New()
	createdAtUnix := time.Now().UTC().Unix()
	expiresAtUnix := time.Now().UTC().Add(24 * time.Hour).Unix()
	token := "test-token-123"

	data := map[string]interface{}{
		"id":         sessionID.String(),
		"user_id":    userID.String(),
		"user_role":  string(domain.User),
		"status":     string(domain.Active),
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0",
		"created_at": createdAtUnix,
		"expires_at": expiresAtUnix,
	}

	session, err := mapper.MapToSession(data, token)

	require.NoError(t, err)
	assert.Equal(t, sessionID, session.ID)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, domain.User, session.UserRole)
	assert.Equal(t, domain.Active, session.Status)
	assert.Equal(t, "192.168.1.1", session.IPAddress)
	assert.Equal(t, "Mozilla/5.0", session.UserAgent)
	assert.Equal(t, token, session.Token)
	assert.Equal(t, createdAtUnix, session.CreatedAt.Unix())
	assert.Equal(t, expiresAtUnix, session.ExpiresAt.Unix())
}

func TestMapToSessionWithStringTimestamps(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}
	sessionID := uuid.New()
	userID := uuid.New()
	createdAtUnix := time.Now().UTC().Unix()
	expiresAtUnix := time.Now().UTC().Add(24 * time.Hour).Unix()
	token := "revoked-token"

	data := map[string]interface{}{
		"id":         sessionID.String(),
		"user_id":    userID.String(),
		"user_role":  string(domain.Admin),
		"status":     string(domain.Revoked),
		"ip_address": "10.0.0.1",
		"user_agent": "Safari/600",
		"created_at": createdAtUnix,
		"expires_at": expiresAtUnix,
	}

	session, err := mapper.MapToSession(data, token)

	require.NoError(t, err)
	assert.Equal(t, domain.Admin, session.UserRole)
	assert.Equal(t, domain.Revoked, session.Status)
	assert.Equal(t, token, session.Token)
}

func TestMapToSessionInvalidSessionID(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}

	data := map[string]interface{}{
		"id":         "invalid-uuid",
		"user_id":    uuid.New().String(),
		"user_role":  string(domain.User),
		"status":     string(domain.Active),
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0",
		"created_at": int64(1000),
		"expires_at": int64(2000),
	}

	_, err := mapper.MapToSession(data, "test-token")

	assert.Error(t, err)
}

func TestMapToSessionInvalidUserID(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}

	data := map[string]interface{}{
		"id":         uuid.New().String(),
		"user_id":    "invalid-uuid",
		"user_role":  string(domain.User),
		"status":     string(domain.Active),
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0",
		"created_at": int64(1000),
		"expires_at": int64(2000),
	}

	_, err := mapper.MapToSession(data, "test-token")

	assert.Error(t, err)
}

func TestMapToSessionInvalidCreatedAtTimestamp(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}

	data := map[string]interface{}{
		"id":         uuid.New().String(),
		"user_id":    uuid.New().String(),
		"user_role":  string(domain.User),
		"status":     string(domain.Active),
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0",
		"created_at": "invalid-timestamp",
		"expires_at": int64(2000),
	}

	_, err := mapper.MapToSession(data, "test-token")

	assert.Error(t, err)
}

func TestRoundTrip(t *testing.T) {
	mapper := &infrastructure.RedisMapper{}
	sessionID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(24 * time.Hour)
	token := "round-trip-token"

	originalSession := domain.Session{
		ID:        sessionID,
		UserID:    userID,
		UserRole:  domain.Admin,
		Status:    domain.Active,
		IPAddress: "192.168.1.100",
		UserAgent: "Firefox/121",
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	// Convert to map and back
	data := mapper.SessionToMap(originalSession)
	recoveredSession, err := mapper.MapToSession(data, token)

	require.NoError(t, err)
	assert.Equal(t, originalSession.ID, recoveredSession.ID)
	assert.Equal(t, originalSession.UserID, recoveredSession.UserID)
	assert.Equal(t, originalSession.UserRole, recoveredSession.UserRole)
	assert.Equal(t, originalSession.Status, recoveredSession.Status)
	assert.Equal(t, originalSession.IPAddress, recoveredSession.IPAddress)
	assert.Equal(t, originalSession.UserAgent, recoveredSession.UserAgent)
	assert.Equal(t, originalSession.Token, recoveredSession.Token)
	assert.Equal(t, originalSession.CreatedAt.Unix(), recoveredSession.CreatedAt.Unix())
	assert.Equal(t, originalSession.ExpiresAt.Unix(), recoveredSession.ExpiresAt.Unix())
}
