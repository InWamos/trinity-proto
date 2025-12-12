package models

import (
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
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null;unique;index:idx_users_id"`
	Username     string    `gorm:"unique;size:32"`
	DisplayName  string
	PasswordHash string
	Role         Role `gorm:"type:enum('user', 'admin')"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (UserModel) TableName() string {
	return "user.users"
}
