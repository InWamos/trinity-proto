package application

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type DemoteUserRequest struct {
	ID uuid.UUID
}

type DemoteUser struct {
	transactionManagerFactory database.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewDemoteUser(
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *DemoteUser {
	dulogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "demote_user"),
	)
	return &DemoteUser{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    dulogger,
	}
}

func (interactor *DemoteUser) Execute(ctx context.Context, input DemoteUserRequest) error {
	interactor.logger.DebugContext(ctx, "Started DemoteUser execution", slog.String("user_id", input.ID.String()))

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	err = userRepository.ChangeUserRoleByID(ctx, input.ID, domain.RoleUser)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to demote user", slog.Any("err", err))
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

	interactor.logger.DebugContext(ctx, "Finished DemoteUser execution")
	return nil
}
