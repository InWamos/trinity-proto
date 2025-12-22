package application

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

// CreateUserRequest Input DTO for interactor.
type CreateUserRequest struct {
	Username    string
	DisplayName string
	Password    string
	Role        domain.Role
}

// CreateUserResponse Output DTO for interactor.
type CreateUserResponse struct {
	UserID uuid.UUID
}

type CreateUser struct {
	passwordHasher            service.PasswordHasher
	uuidGenerator             *service.UUIDGenerator
	transactionManagerFactory database.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewCreateUser(
	passwordHasher service.PasswordHasher,
	uuidGenerator *service.UUIDGenerator,
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *CreateUser {
	culogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "create_user"),
	)
	return &CreateUser{
		passwordHasher:            passwordHasher,
		uuidGenerator:             uuidGenerator,
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    culogger,
	}
}

func (interactor *CreateUser) Execute(ctx context.Context, input CreateUserRequest) (*CreateUserResponse, error) {
	interactor.logger.DebugContext(ctx, "Started Create User execution")

	idp := ctx.Value("IdentityProvider").(*client.UserIdentity)
	if err := service.AuthorizeByRole(idp, domain.RoleAdmin); err != nil {
		return nil, ErrInsufficientPrivileges
	}

	passwordHashed, err := interactor.passwordHasher.HashPassword(input.Password)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "The password hasher has failed")
		return nil, ErrHashingFailed
	}

	var randomUUID uuid.UUID
	if randomUUID, err = interactor.uuidGenerator.GetUUIDv7(); err != nil {
		interactor.logger.ErrorContext(ctx, "The uuid generator has failed")
		return nil, ErrUUIDGeneration
	}

	newUser := domain.NewUser(randomUUID, input.Username, input.DisplayName, passwordHashed, input.Role)

	// Execute within a transaction managed by the factory
	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	// Get repository scoped to this transaction
	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	err = userRepository.CreateUser(ctx, *newUser)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create user", slog.Any("err", err))
		if rollbackErr := transactionManager.Rollback(ctx); rollbackErr != nil {
			interactor.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("err", rollbackErr))
		}
		return nil, ErrDatabaseFailed
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return nil, ErrDatabaseFailed
	}

	interactor.logger.DebugContext(ctx, "Finished Create User execution")
	return &CreateUserResponse{UserID: randomUUID}, nil
}
