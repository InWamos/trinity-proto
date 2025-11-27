package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure"
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
	userRepository     infrastructure.UserRepository
}

func NewCreateUser(
	passwordHasher service.PasswordHasher,
	uuidGenerator service.UUIDGenerator,
	transactionManager interfaces.TransactionManager,
	userRepository infrastructure.UserRepository,
) *CreateUser {
	return &CreateUser{
		passwordHasher:     passwordHasher,
		uuidGenerator:      uuidGenerator,
		transactionManager: transactionManager,
		userRepository:     userRepository,
	}
}

