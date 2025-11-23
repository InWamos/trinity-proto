package application

import (
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure"
)

type CreateUserRequest struct {
	Username    string
	DisplayName string
	Password    string
	Role        domain.Role
}

type CreateUser struct {
	passwordHasher     service.PasswordHasher
	transactionManager interfaces.TransactionManager
	userRepository     infrastructure.UserRepository
}
