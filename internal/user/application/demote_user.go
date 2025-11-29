package application

import (
	"context"
	"log/slog"

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
	logger             *slog.Logger
}

func NewDemoteUser(
	transactionManager interfaces.TransactionManager,
	userRepository repository.UserRepository,
	logger *slog.Logger,
) *DemoteUser {
	dulogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "demote_user"),
	)
	return &DemoteUser{
		transactionManager: transactionManager,
		userRepository:     userRepository,
		logger:             dulogger,
	}
}

func (i *DemoteUser) Execute(ctx context.Context, input DemoteUserRequest) error {
	return nil
}
