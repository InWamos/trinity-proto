package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type DemoteUserRequest struct {
	ID uuid.UUID
}

type DemoteUser struct {
	transactionManager interfaces.TransactionManager
	userRepository     repository.UserRepository
}
