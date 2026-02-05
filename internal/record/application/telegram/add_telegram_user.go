package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/domain/telegram/service"
	userDomain "github.com/InWamos/trinity-proto/internal/user/domain"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/google/uuid"
)

type AddTelegramUserRequest struct {
	TelegramID uint64
}

type AddTelegramUserResponse struct {
	UserID string
}

type AddTelegramUser struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	telegramDomainValidator   *service.TelegramModelValidator
	telegramUserFactory       repository.TelegramUserRepositoryFactory
	logger                    *slog.Logger
}

func NewAddTelegramUser(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	telegramDomainValidator *service.TelegramModelValidator,
	telegramUserFactory repository.TelegramUserRepositoryFactory,
	logger *slog.Logger,
) *AddTelegramUser {
	iLogger := logger.With(
		slog.String("module", "record"),
		slog.String("name", "add_telegram_user"),
	)
	return &AddTelegramUser{
		transactionManagerFactory: transactionManagerFactory,
		telegramUserFactory:       telegramUserFactory,
		telegramDomainValidator:   telegramDomainValidator,
		logger:                    iLogger,
	}
}

func (interactor *AddTelegramUser) Execute(
	ctx context.Context,
	input AddTelegramUserRequest,
) (*AddTelegramUserResponse, error) {
	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	if err := rbac.AuthorizeByRole(idp, userDomain.RoleUser); err != nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	userID := uuid.New()
	now := time.Now()

	interactor.logger.DebugContext(
		ctx,
		"Started AddTelegramUser execution",
		slog.String("user_id", userID.String()),
		slog.Uint64("telegram_id", input.TelegramID),
	)

	telegramUser := &domain.TelegramUser{
		ID:          userID,
		TelegramID:  input.TelegramID,
		AddedAt:     now,
		AddedByUser: idp.UserID,
	}
	// Validate the rules before adding to the database
	if err := interactor.telegramDomainValidator.Validate(telegramUser); err != nil {
		return nil, err
	}

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	telegramUserRepository := interactor.telegramUserFactory.CreateTelegramUserRepositoryWithTransaction(
		transactionManager,
	)

	if err = telegramUserRepository.AddUser(ctx, telegramUser); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			interactor.logger.WarnContext(
				ctx,
				"telegram user already exists",
				slog.Uint64("telegram_id", input.TelegramID),
			)
			return nil, err
		default:
			interactor.logger.ErrorContext(
				ctx,
				"failed to add telegram user",
				slog.String("user_id", userID.String()),
				slog.Any("err", err),
			)
			return nil, ErrDatabaseFailed
		}
	}
	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	return &AddTelegramUserResponse{
		UserID: userID.String(),
	}, nil
}
