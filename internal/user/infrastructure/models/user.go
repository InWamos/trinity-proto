package models

import (
	"time"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type role string

const (
	RoleUser  role = "user"
	RoleAdmin role = "admin"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null;unique;index:idx_users_id"`
	Username     string    `gorm:"unique;size:32"`
	DisplayName  string
	PasswordHash string
	Role         role `gorm:"type:enum('user', 'admin')"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (model *UserModel) ToEntity() *domain.User {
	return domain.NewUser(model.ID, model.Username, model.DisplayName, model.PasswordHash, domain.Role(model.Role))
}
