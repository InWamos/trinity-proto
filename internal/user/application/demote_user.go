package application

import (
	"context"

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

func NewDemoteUser(
	transactionManager interfaces.TransactionManager,
	userRepository repository.UserRepository,
) *DemoteUser {
	return &DemoteUser{
		transactionManager: transactionManager,
		userRepository:     userRepository,
	}
}

func (i *DemoteUser) Execute(ctx context.Context, input DemoteUserRequest) error {
	return nil
}
