package application

import (
	"context"
	"errors"

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

// Input DTO.
type createUserRequest struct {
	Username    string
	DisplayName string
	Password    string
	Role        domain.Role
}

type CreateUser struct {
	passwordHasher     service.PasswordHasher
	uuidGenerator      service.UUIDGenerator
	transactionManager interfaces.TransactionManager
	userRepository     repository.UserRepository
}

func NewCreateUser(
	passwordHasher service.PasswordHasher,
	uuidGenerator service.UUIDGenerator,
	transactionManager interfaces.TransactionManager,
	userRepository repository.UserRepository,
) *CreateUser {
	return &CreateUser{
		passwordHasher:     passwordHasher,
		uuidGenerator:      uuidGenerator,
		transactionManager: transactionManager,
		userRepository:     userRepository,
	}
}

func (interactor *CreateUser) Execute(ctx context.Context, input createUserRequest) error {
	passwordHashed, err := interactor.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return ErrHashingFailed
	}
	var randomUUID uuid.UUID
	if randomUUID, err = interactor.uuidGenerator.GetUUIDv7(); err != nil {
		return ErrUUIDGeneration
	}
	newUser := domain.NewUser(randomUUID, input.Username, input.DisplayName, passwordHashed, input.Role)
	return interactor.userRepository.CreateUser(ctx, *newUser)
}
