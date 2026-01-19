package application //nolint:dupl //Interactors should be nearly the same

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/google/uuid"
)

type PromoteUserRequest struct {
	ID uuid.UUID
}

type PromoteUser struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewPromoteUser(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *PromoteUser {
	pulogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "promote_user"),
	)
	return &PromoteUser{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    pulogger,
	}
}

func (interactor *PromoteUser) Execute(ctx context.Context, input PromoteUserRequest) error {
	interactor.logger.DebugContext(ctx, "Started PromoteUser execution", slog.String("user_id", input.ID.String()))

	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return ErrInsufficientPrivileges
	}

	if err := rbac.AuthorizeByRole(idp, domain.RoleAdmin); err != nil {
		return ErrInsufficientPrivileges
	}

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	err = userRepository.ChangeUserRoleByID(ctx, input.ID, domain.RoleAdmin)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to promote user", slog.Any("err", err))
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

	interactor.logger.DebugContext(ctx, "Finished PromoteUser execution")
	return nil
}
