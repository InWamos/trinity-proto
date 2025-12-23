package application

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type CreateRandomAdminUser struct {
	passwordHasher            service.PasswordHasher
	uuidGenerator             *service.UUIDGenerator
	transactionManagerFactory database.TransactionManagerFactory
	userRepositoryFactory     repository.UserRepositoryFactory
	logger                    *slog.Logger
}

func NewCreateRandomAdminUser(
	passwordHasher service.PasswordHasher,
	uuidGenerator *service.UUIDGenerator,
	transactionManagerFactory database.TransactionManagerFactory,
	userRepositoryFactory repository.UserRepositoryFactory,
	logger *slog.Logger,
) *CreateRandomAdminUser {
	culogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "create_random_admin_user"),
	)
	return &CreateRandomAdminUser{
		passwordHasher:            passwordHasher,
		uuidGenerator:             uuidGenerator,
		transactionManagerFactory: transactionManagerFactory,
		userRepositoryFactory:     userRepositoryFactory,
		logger:                    culogger,
	}
}

func (interactor *CreateRandomAdminUser) Execute(
	ctx context.Context,
) error {
	interactor.logger.DebugContext(ctx, "Started Create User execution")
	// Generate randomly safe password
	password, err := service.GenerateSafeRandomString(16)
	if err != nil {
		return err
	}

	passwordHashed, err := interactor.passwordHasher.HashPassword(password)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "The password hasher has failed")
		return ErrHashingFailed
	}

	var randomUUID uuid.UUID
	if randomUUID, err = interactor.uuidGenerator.GetUUIDv7(); err != nil {
		interactor.logger.ErrorContext(ctx, "The uuid generator has failed")
		return ErrUUIDGeneration
	}

	newUser := domain.NewUser(randomUUID, "admin", "admin", passwordHashed, domain.RoleAdmin)

	// Execute within a transaction managed by the factory
	transactionManager, err := interactor.transactionManagerFactory.NewTransaction(ctx)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create transaction", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	// Get repository scoped to this transaction
	userRepository := interactor.userRepositoryFactory.CreateUserRepositoryWithTransaction(transactionManager)

	err = userRepository.CreateUser(ctx, *newUser)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create user", slog.Any("err", err))
		if rollbackErr := transactionManager.Rollback(ctx); rollbackErr != nil {
			interactor.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("err", rollbackErr))
		}
		return ErrDatabaseFailed
	}

	if err = transactionManager.Commit(ctx); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	interactor.logger.DebugContext(ctx, "Finished Create User execution")
	interactor.logger.InfoContext(
		ctx,
		"New admin user has been created. You can login now. BTW Change the password after you login",
		slog.String("username", "admin"),
		slog.String("password", password),
	)
	return nil
}
