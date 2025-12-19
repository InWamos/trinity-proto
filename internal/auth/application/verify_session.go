package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionRevoked  = errors.New("session revoked")
)

type VerifySessionRequest struct {
	Session_id string
}

type VerifySession struct {
	sessionRepository infrastructure.SessionRepository
	logger            *slog.Logger
}

func NewVerifySession(
	sessionRepository infrastructure.SessionRepository,
	logger *slog.Logger,
) *VerifySession {
	vsLogger := logger.With(slog.String("module", "auth"), slog.String("name", "verify_session"))
	return &VerifySession{
		sessionRepository: sessionRepository,
		logger:            vsLogger,
	}
}

func (vs *VerifySession) Execute(ctx context.Context, request VerifySessionRequest) error {
	// Retrieve the session from repository
	session, err := vs.sessionRepository.GetSessionByToken(ctx, request.Session_id)
	if err != nil {
		vs.logger.ErrorContext(ctx, "failed to retrieve session", slog.Any("err", err))
		return ErrSessionNotFound
	}

	// Check if session was found
	if session.ID.String() == "" {
		vs.logger.InfoContext(ctx, "session not found", slog.String("token", request.Session_id))
		return ErrSessionNotFound
	}

	// Check if session has expired
	if time.Now().UTC().After(session.ExpiresAt) {
		vs.logger.InfoContext(ctx, "session has expired", slog.String("session_id", session.ID.String()))
		return ErrSessionExpired
	}

	// Check if session is revoked
	if session.Status == domain.Revoked {
		vs.logger.InfoContext(ctx, "session is revoked", slog.String("session_id", session.ID.String()))
		return ErrSessionRevoked
	}

	vs.logger.DebugContext(ctx, "session verified successfully",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
	)

	return nil
}
