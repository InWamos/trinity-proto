package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
)

type ValidateUserCredentialsRequest struct {
	Username string
	Password string
}

type ValidateUserCredentials struct {
	transactionManagerFactory database.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewValidateUserCredentials(
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *ValidateUserCredentials {
	vucLogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "validate_user_credentials"),
	)
	return &ValidateUserCredentials{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    vucLogger,
	}
}
