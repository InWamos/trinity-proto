package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type PromoteUserRequest struct {
	ID uuid.UUID
}

type PromoteUser struct {
	transactionManager interfaces.TransactionManager
	userRepository     repository.UserRepository
	logger             *slog.Logger
}
