package infrastructure

import (
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(id uuid.UUID) (domain.User, error)
	GetUserByUsername(username string) (domain.User, error)
	RemoveUserByID(id uuid.UUID) error
	ChangeUserRoleById(id uuid.UUID, changeToRole domain.Role) error
	CreateUser(user domain.User) error
}
