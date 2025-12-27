package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// UserModelSqlx is for sqlx repositories.
type UserModelSqlx struct {
	ID           uuid.UUID    `db:"id"`
	Username     string       `db:"username"`
	DisplayName  string       `db:"display_name"`
	PasswordHash string       `db:"password_hash"`
	Role         Role         `db:"role"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}
