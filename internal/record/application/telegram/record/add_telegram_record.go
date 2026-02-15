package record

import (
	"context"
	"log/slog"
	"time"

	application "github.com/InWamos/trinity-proto/internal/record/application/telegram"
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/domain/telegram/service"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	userDomain "github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/google/uuid"
)

type AddTelegramRecordRequest struct {
	MessageTelegramID  uint64
	FromUserTelegramID uint64
	InTelegramChatID   int64
	MessageText        string
	PostedAt           time.Time
}

type AddTelegramRecordResponse struct {
	RecordID string
}

type AddTelegramRecord struct {
	transactionManagerFactory       interfaces.TransactionManagerFactory
	telegramDomainValidator         *service.TelegramModelValidator
	telegramRecordRepositoryFactory repository.TelegramRecordRepositoryFactory
	logger                          *slog.Logger
}

func NewAddTelegramRecord(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	telegramDomainValidator *service.TelegramModelValidator,
	telegramRecordRepositoryFactory repository.TelegramRecordRepositoryFactory,
	logger *slog.Logger,
) *AddTelegramRecord {
	iLogger := logger.With(
		slog.String("module", "record"),
		slog.String("name", "add_telegram_record"),
	)
	return &AddTelegramRecord{
		transactionManagerFactory:       transactionManagerFactory,
		telegramRecordRepositoryFactory: telegramRecordRepositoryFactory,
		telegramDomainValidator:         telegramDomainValidator,
		logger:                          iLogger,
	}
}

func (interactor *AddTelegramRecord) Execute(
	ctx context.Context,
	input AddTelegramRecordRequest,
) (*AddTelegramRecordResponse, error) {
	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	if err := rbac.AuthorizeByRole(idp, userDomain.RoleUser); err != nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	recordID := uuid.New()
	now := time.Now()

	interactor.logger.DebugContext(
		ctx,
		"Started AddTelegramRecord execution",
		slog.String("record_id", recordID.String()),
		slog.Uint64("message_telegram_id", input.MessageTelegramID),
	)

	telegramRecord := domain.TelegramRecord{
		ID:                 recordID,
		MessageTelegramID:  input.MessageTelegramID,
		FromUserTelegramID: input.FromUserTelegramID,
		InTelegramChatID:   input.InTelegramChatID,
		MessageText:        input.MessageText,
		PostedAt:           input.PostedAt,
		AddedAt:            now,
		AddedByUser:        idp.UserID,
	}

	// Validate the rules before adding to the database
	if err := interactor.telegramDomainValidator.Validate(&telegramRecord); err != nil {
		return nil, err
	}

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, application.ErrDatabaseFailed
	}

	telegramRecordRepository := interactor.telegramRecordRepositoryFactory.CreateTelegramRecordRepositoryWithTransaction(
		transactionManager,
	)

	if err = telegramRecordRepository.CreateTelegramRecord(ctx, telegramRecord); err != nil {
		interactor.logger.InfoContext(
			ctx,
			"failed to add telegram record",
			slog.String("record_id", recordID.String()),
			slog.Any("err", err),
		)
		if rollbackErr := transactionManager.Rollback(ctx); rollbackErr != nil {
			interactor.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("err", rollbackErr))
		}
		return nil, err
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return nil, application.ErrDatabaseFailed
	}

	interactor.logger.DebugContext(ctx, "Finished AddTelegramRecord execution")
	return &AddTelegramRecordResponse{
		RecordID: recordID.String(),
	}, nil
}
