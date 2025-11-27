package infrastructure

import (
	"context"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	RemoveUserByID(ctx context.Context, id uuid.UUID) error
	ChangeUserRoleById(ctx context.Context, id uuid.UUID, changeToRole domain.Role) error
	CreateUser(ctx context.Context, user domain.User) error
}
