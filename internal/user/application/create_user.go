package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

var (
	ErrHashingFailed           = errors.New("password hashing failed")
	ErrUUIDGeneration          = errors.New("UUID generation failed")
	ErrUserWithIDAlreadyExists = errors.New("this uuid is already in the database")
	ErrDatabaseFailed          = errors.New("the database operation has failed")
)

// CreateUserRequest Input DTO for interactor.
type CreateUserRequest struct {
	Username    string
	DisplayName string
	Password    string
	Role        domain.Role
}

type CreateUser struct {
	passwordHasher     service.PasswordHasher
	uuidGenerator      *service.UUIDGenerator
	transactionManager interfaces.TransactionManager
	userRepository     repository.UserRepository
	logger             *slog.Logger
}

func NewCreateUser(
	passwordHasher service.PasswordHasher,
	uuidGenerator *service.UUIDGenerator,
	transactionManager interfaces.TransactionManager,
	userRepository repository.UserRepository,
	logger *slog.Logger,
) *CreateUser {
	culogger := logger.With(
		slog.String("component", "interactor"),
		slog.String("name", "create_user"),
	)
	return &CreateUser{
		passwordHasher:     passwordHasher,
		uuidGenerator:      uuidGenerator,
		transactionManager: transactionManager,
		userRepository:     userRepository,
		logger:             culogger,
	}
}

func (interactor *CreateUser) Execute(ctx context.Context, input CreateUserRequest) error {
	interactor.logger.DebugContext(ctx, "Started Create User execution")
	passwordHashed, err := interactor.passwordHasher.HashPassword(input.Password)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "The password hasher has failed")
		return ErrHashingFailed
	}
	var randomUUID uuid.UUID
	if randomUUID, err = interactor.uuidGenerator.GetUUIDv7(); err != nil {
		interactor.logger.ErrorContext(ctx, "The uuid generator has failed")
		return ErrUUIDGeneration
	}
	newUser := domain.NewUser(randomUUID, input.Username, input.DisplayName, passwordHashed, input.Role)
	err = interactor.userRepository.CreateUser(ctx, *newUser)
	if err != nil {
		interactor.logger.ErrorContext(ctx, "failed to create user", slog.Any("err", err))
		return ErrDatabaseFailed
	}

	if err = interactor.transactionManager.Commit(); err != nil {
		interactor.logger.ErrorContext(ctx, "failed to commit", slog.Any("err", err))
		return ErrDatabaseFailed
	}
	interactor.logger.DebugContext(ctx, "Finished Create User execution")
	return nil
}
