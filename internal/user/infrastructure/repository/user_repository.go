package repository

import (
	"context"
	"errors"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user was not found")
	ErrUserCreationFailed = errors.New("failed to save user")
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	RemoveUserByID(ctx context.Context, id uuid.UUID) error
	ChangeUserRoleByID(ctx context.Context, id uuid.UUID, changeToRole domain.Role) error
	CreateUser(ctx context.Context, user domain.User) error
}

type UserRepositoryFactory interface {
	CreateUserRepositoryWithTransaction(tm database.TransactionManager) UserRepository
}
