package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure"
	"github.com/google/uuid"
)

type PromoteUserRequest struct {
	ID uuid.UUID
}

type PromoteUser struct {
	transactionManager interfaces.TransactionManager
	userRepository     infrastructure.UserRepository
}
