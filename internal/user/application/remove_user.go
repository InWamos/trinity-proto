package application

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/google/uuid"
)

type RemoveUserRequest struct {
	ID uuid.UUID
}

type RemoveUser struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewRemoveUser(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *RemoveUser {
	rulogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "remove_user"),
	)
	return &RemoveUser{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    rulogger,
	}
}

func (interactor *RemoveUser) Execute(ctx context.Context, input RemoveUserRequest) error {
	interactor.logger.DebugContext(ctx, "Started RemoveUser execution", slog.String("user_id", input.ID.String()))

	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return ErrInsufficientPrivileges
	}

	if err := service.AuthorizeByRole(idp, domain.RoleAdmin); err != nil {
		return ErrInsufficientPrivileges
	}

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	err = userRepository.RemoveUserByID(ctx, input.ID)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to remove user", slog.Any("err", err))
		if rollbackErr := transactionManager.Rollback(ctx); rollbackErr != nil {
			interactor.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("err", rollbackErr))
		}

		// Check if the error is because the user was not found
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return ErrDatabaseFailed
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	interactor.logger.DebugContext(ctx, "Finished RemoveUser execution")
	return nil
}
