package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/user/client"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnexpected         = errors.New("unexpected error")
)

type AddSessionRequest struct {
	Username  string
	Password  string
	IpAddress string
	UserAgent string
}

type AddSession struct {
	sessionRepository infrastructure.SessionRepository
	userClient        client.UserClient
	logger            *slog.Logger
}

func NewAddSession(
	sessionRepository infrastructure.SessionRepository,
	userClient client.UserClient,
	logger *slog.Logger,
) *AddSession {
	asLogger := logger.With(slog.String("module", "auth"), slog.String("name", "add_session"))
	return &AddSession{
		sessionRepository: sessionRepository,
		userClient:        userClient,
		logger:            asLogger,
	}
}
