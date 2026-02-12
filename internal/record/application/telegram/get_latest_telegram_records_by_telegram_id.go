package application

import (
	"context"
	"errors"
	"log/slog"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	userDomain "github.com/InWamos/trinity-proto/internal/user/domain"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/middleware"
)

type GetLatestTelegramRecordsByUserTelegramID struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	telegramRecordFactory     repository.TelegramRecordRepositoryFactory
	logger                    *slog.Logger
}

type GetLatestTelegramRecordsByUserTelegramIDRequest struct {
	UserTelegramID uint64
}

type GetLatestTelegramRecordsByUserTelegramIDResponse struct {
	TelegramRecords *[]domain.TelegramRecord
}

func NewGetLatestTelegramRecordsByUserTelegramID(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	telegramRecordFactory repository.TelegramRecordRepositoryFactory,
	logger *slog.Logger,
) *GetLatestTelegramRecordsByUserTelegramID {
	iLogger := logger.With(
		slog.String("module", "record"),
		slog.String("name", "get_latest_tg_records_by_user_telegram_id"),
	)
	return &GetLatestTelegramRecordsByUserTelegramID{
		transactionManagerFactory: transactionManagerFactory,
		telegramRecordFactory:     telegramRecordFactory,
		logger:                    iLogger,
	}
}

func (interactor *GetLatestTelegramRecordsByUserTelegramID) Execute(
	ctx context.Context,
	input GetLatestTelegramRecordsByUserTelegramIDRequest,
) (*GetLatestTelegramRecordsByUserTelegramIDResponse, error) {
	interactor.logger.DebugContext(ctx, "Started GetUserByID execution", slog.Uint64("user_id", input.UserTelegramID))
	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	if err := rbac.AuthorizeByRole(idp, userDomain.RoleUser); err != nil {
		return nil, rbac.ErrInsufficientPrivileges
	}
	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}
	recordRepository := interactor.telegramRecordFactory.CreateTelegramRecordRepositoryWithTransaction(
		transactionManager,
	)
	records, err := recordRepository.GetLatestTelegramRecordsByUserTelegramID(ctx, input.UserTelegramID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoRecordsForThisTelegramID):
			interactor.logger.ErrorContext(
				ctx,
				"no records for this telegram id",
				slog.Uint64("id", input.UserTelegramID),
				slog.Any("err", err),
			)
			return nil, err
		default:
			interactor.logger.InfoContext(
				ctx,
				"failed to retrieve telegram records by telegram id",
				slog.Any("err", err),
			)
			return nil, ErrDatabaseFailed
		}
	}
	return &GetLatestTelegramRecordsByUserTelegramIDResponse{
		TelegramRecords: records,
	}, nil
}
