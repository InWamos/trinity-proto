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
	SessionID string
}

type VerifySessionResponse struct {
	Session domain.Session
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

func (vs *VerifySession) Execute(ctx context.Context, request VerifySessionRequest) (domain.Session, error) {
	// Retrieve the session from repository
	session, err := vs.sessionRepository.GetSessionByToken(ctx, request.SessionID)
	if err != nil {
		vs.logger.ErrorContext(ctx, "failed to retrieve session", slog.Any("err", err))
		return domain.Session{}, ErrSessionNotFound
	}

	// Check if session was found
	if session.ID.String() == "" {
		vs.logger.InfoContext(ctx, "session not found", slog.String("token", request.SessionID))
		return domain.Session{}, ErrSessionNotFound
	}

	// Check if session has expired
	if time.Now().UTC().After(session.ExpiresAt) {
		vs.logger.InfoContext(ctx, "session has expired", slog.String("session_id", session.ID.String()))
		return domain.Session{}, ErrSessionExpired
	}

	// Check if session is revoked
	if session.Status == domain.Revoked {
		vs.logger.InfoContext(ctx, "session is revoked", slog.String("session_id", session.ID.String()))
		return domain.Session{}, ErrSessionRevoked
	}

	vs.logger.DebugContext(ctx, "session verified successfully",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
	)

	return session, nil
}
