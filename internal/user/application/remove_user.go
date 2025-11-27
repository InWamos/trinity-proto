package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type RemoveUserRequest struct {
	ID uuid.UUID
}

type RemoveUser struct {
	transactionManager interfaces.TransactionManager
	userRepository     repository.UserRepository
}
