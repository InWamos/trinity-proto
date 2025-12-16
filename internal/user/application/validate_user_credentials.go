package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/application/service"
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
	passwordHasher            service.PasswordHasher
	logger                    *slog.Logger
}

func NewValidateUserCredentials(
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	passwordHasher service.PasswordHasher,
	logger *slog.Logger,
) *ValidateUserCredentials {
	vucLogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "validate_user_credentials"),
	)
	return &ValidateUserCredentials{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		passwordHasher:            passwordHasher,
		logger:                    vucLogger,
	}
}

func (interactor *ValidateUserCredentials) Execute(
	ctx context.Context,
	input ValidateUserCredentialsRequest,
) error {
	interactor.logger.DebugContext(
		ctx,
		"Started ValidateUserCredentials execution",
		slog.String("user_id", input.Username),
	)

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return ErrDatabaseFailed
	}
	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)
	user, err := userRepository.GetUserByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUsernameAbsent
		} else {
			return ErrDatabaseFailed
		}
	}
	if err := interactor.passwordHasher.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		return ErrPasswordMismatch
	} else {
		return nil
	}
}
