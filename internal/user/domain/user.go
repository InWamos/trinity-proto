package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Role string

// All roles Enum.
const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           uuid.UUID
	Username     string
	DisplayName  string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
}

func NewUser(uuid7 uuid.UUID, username string, displayName string, passwordHash string, role Role) *User {
	// Consider using uuidv7 for a better indexing
	return &User{
		ID:           uuid7,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		DeletedAt:    sql.NullTime{Valid: false},
	}
}
