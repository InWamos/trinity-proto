package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type GetUserByIDRequest struct {
	ID uuid.UUID
}

type GetUserByIDResponse struct {
	User domain.User
}

type GetUserByID struct {
	transactionManagerFactory database.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewGetUserByID(
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *GetUserByID {
	gulogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "get_user_by_id"),
	)
	return &GetUserByID{
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    gulogger,
	}
}

func (interactor *GetUserByID) Execute(ctx context.Context, input GetUserByIDRequest) (*GetUserByIDResponse, error) {
	interactor.logger.DebugContext(ctx, "Started GetUserByID execution", slog.String("user_id", input.ID.String()))

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	user, err := userRepository.GetUserByID(ctx, input.ID)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to get user", slog.Any("err", err))
		if rollbackErr := transactionManager.Rollback(ctx); rollbackErr != nil {
			interactor.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("err", rollbackErr))
		}
		return nil, ErrUserNotFound
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	interactor.logger.DebugContext(ctx, "Finished GetUserByID execution")
	return &GetUserByIDResponse{User: user}, nil
}
