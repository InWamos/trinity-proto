package application

import (
	"context"
	"errors"
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

type AddTelegramIdentityRequest struct {
	UserID      uuid.UUID
	Username    string
	FirstName   string
	LastName    string
	Bio         string
	PhoneNumber string
}

type AddTelegramIdentityResponse struct {
	ID uuid.UUID
}

type AddTelegramIdentity struct {
	transactionManagerFactory       interfaces.TransactionManagerFactory
	telegramIdentityDomainValidator *service.TelegramModelValidator
	telegramIdentityFactory         repository.TelegramIdentityRepositoryFactory
	logger                          *slog.Logger
}

func NewAddTelegramIdentity(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	telegramIdentityDomainValidator *service.TelegramModelValidator,
	telegramIdentityFactory repository.TelegramIdentityRepositoryFactory,
	logger *slog.Logger,
) *AddTelegramIdentity {
	iLogger := logger.With(
		slog.String("module", "record"),
		slog.String("name", "add_telegram_identity"),
	)
	return &AddTelegramIdentity{
		transactionManagerFactory:       transactionManagerFactory,
		telegramIdentityFactory:         telegramIdentityFactory,
		telegramIdentityDomainValidator: telegramIdentityDomainValidator,
		logger:                          iLogger,
	}
}

func (interactor *AddTelegramIdentity) Execute(
	ctx context.Context,
	input AddTelegramIdentityRequest,
) (*AddTelegramIdentityResponse, error) {
	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	if err := rbac.AuthorizeByRole(idp, userDomain.RoleUser); err != nil {
		return nil, rbac.ErrInsufficientPrivileges
	}

	identityID := uuid.New()
	now := time.Now()

	interactor.logger.DebugContext(
		ctx,
		"Started AddTelegramIdentity execution",
		slog.String("identity_id", identityID.String()),
		slog.String("user_id", input.UserID.String()),
		slog.String("username", input.Username),
	)

	telegramIdentity := &domain.TelegramIdentity{
		ID:          identityID,
		UserID:      input.UserID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Username:    input.Username,
		PhoneNumber: input.PhoneNumber,
		Bio:         input.Bio,
		AddedAt:     now,
		AddedByUser: idp.UserID,
	}

	// Validate the rules before adding to the database
	if err := interactor.telegramIdentityDomainValidator.Validate(telegramIdentity); err != nil {
		return nil, err
	}

	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, application.ErrDatabaseFailed
	}

	telegramIdentityRepository := interactor.telegramIdentityFactory.CreateTelegramIdentityRepositoryWithTransaction(
		transactionManager,
	)

	if err = telegramIdentityRepository.AddIdentity(ctx, telegramIdentity); err != nil {
		switch {
		case errors.Is(err, domain.ErrIdentityAlreadyExists):
			interactor.logger.DebugContext(
				ctx,
				"telegram identity already exists",
				slog.String("user_id", input.UserID.String()),
				slog.String("username", input.Username),
			)
			return nil, err
		case errors.Is(err, domain.ErrUnexistentTelegramUserReferenced):
			interactor.logger.DebugContext(
				ctx,
				"telegram identity linked to unexistent telegram user",
				slog.String("user_id", input.UserID.String()),
				slog.String("username", input.Username),
			)
			return nil, err

		default:
			interactor.logger.ErrorContext(
				ctx,
				"failed to add telegram identity",
				slog.String("identity_id", identityID.String()),
				slog.Any("err", err),
			)
			return nil, application.ErrDatabaseFailed
		}
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return nil, application.ErrDatabaseFailed
	}

	return &AddTelegramIdentityResponse{
		ID: identityID,
	}, nil
}
