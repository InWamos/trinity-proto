package application

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
)

type RemoveSessionRequest struct {
	Token string
}

type RemoveSession struct {
	sessionRepository infrastructure.SessionRepository
	logger            *slog.Logger
}

func NewRemoveSession(
	sessionRepository infrastructure.SessionRepository,
	logger *slog.Logger,
) *RemoveSession {
	rsLogger := logger.With(slog.String("module", "auth"), slog.String("name", "remove_session"))
	return &RemoveSession{
		sessionRepository: sessionRepository,
		logger:            rsLogger,
	}
}

func (rsInteractor *RemoveSession) Execute(ctx context.Context, input RemoveSessionRequest) error {
	err := rsInteractor.sessionRepository.RevokeSessionByToken(ctx, input.Token)
	if err != nil {
		rsInteractor.logger.ErrorContext(ctx, "Failed to revoke session", slog.Any("err", err))
		return domain.ErrSessionNotFound
	}
	rsInteractor.logger.DebugContext(ctx, "Session revoked successfully")
	return nil
}
