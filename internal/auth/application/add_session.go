package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
)

type AddSessionRequest struct {
	Username string
	Password string
}

type AddSession struct {
	sessionRepository infrastructure.SessionRepository
	logger            *slog.Logger
}

func NewAddSession(sessionRepository infrastructure.SessionRepository, logger *slog.Logger) *AddSession {
	asLogger := logger.With(slog.String("component", "interactor"), slog.String("name", "add_session"))
	return &AddSession{
		sessionRepository: sessionRepository,
		logger:            asLogger,
	}
}
