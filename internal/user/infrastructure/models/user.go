package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// UserModelSqlx is for sqlx repositories.
type UserModelSqlx struct {
	ID           uuid.UUID    `db:"id"`
	Username     string       `db:"username"`
	DisplayName  string       `db:"display_name"`
	PasswordHash string       `db:"password_hash"`
	UserRole     UserRole     `db:"user_role"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}
