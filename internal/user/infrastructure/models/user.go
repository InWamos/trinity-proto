package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type UserModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;not null;unique;index:idx_users_id" db:"id"`
	Username     string         `gorm:"unique;size:32" db:"username"`
	DisplayName  string         `db:"display_name"`
	PasswordHash string         `db:"password_hash"`
	Role         Role           `gorm:"type:enum('user', 'admin')" db:"role"`
	CreatedAt    time.Time      `db:"created_at"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at"`
}

// UserModelSqlx is for sqlx repositories
type UserModelSqlx struct {
	ID           uuid.UUID    `db:"id"`
	Username     string       `db:"username"`
	DisplayName  string       `db:"display_name"`
	PasswordHash string       `db:"password_hash"`
	Role         Role         `db:"role"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

func (UserModel) TableName() string {
	return "user.users"
}
