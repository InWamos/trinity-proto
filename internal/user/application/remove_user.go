package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure"
	"github.com/google/uuid"
)

type RemoveUserRequest struct {
	id uuid.UUID
}

type RemoveUser struct {
	transactionManager interfaces.TransactionManager
	userRepository     infrastructure.UserRepository
}
