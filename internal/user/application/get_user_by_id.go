package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
)

type GetUserByIDRequest struct {
	ID uuid.UUID
}

type GetUserByIDResponce struct {
	User domain.User
}

type GetUserByID struct {
	userRepository repository.UserRepository
	logger         *slog.Logger
}
