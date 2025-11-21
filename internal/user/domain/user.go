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

func NewUser(username string, displayName string, passwordHash string, role Role) *User {
	randomUUID4 := uuid.New()
	return &User{
		ID:           randomUUID4,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		DeletedAt:    sql.NullTime{Valid: false},
	}
}
