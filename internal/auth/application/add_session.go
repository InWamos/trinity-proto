package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/user/client"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnexpected         = errors.New("unexpected error")
)

const DefaultSessionDuration = 24 * time.Hour

type AddSessionRequest struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

type AddSessionResponse struct {
	Session domain.Session
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

func (asInteractor *AddSession) Execute(ctx context.Context, input AddSessionRequest) (AddSessionResponse, error) {
	// Check whether the provided credentials valid
	response, err := asInteractor.userClient.VerifyCredentials(ctx, input.Username, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrUsernameAbsent):
			asInteractor.logger.InfoContext(ctx, "Username doesn't exist", slog.String("username", input.Username))
			return AddSessionResponse{}, ErrInvalidCredentials
		case errors.Is(err, client.ErrPasswordMissmatch):
			asInteractor.logger.InfoContext(ctx, "Password didn't match", slog.String("username", input.Username))
			return AddSessionResponse{}, ErrInvalidCredentials
		default:
			asInteractor.logger.ErrorContext(ctx, "Unexpected error during credential verification")
			return AddSessionResponse{}, ErrUnexpected
		}
	}

	newSession, err := domain.NewSession(
		response.UserID,
		domain.UserRole(response.UserRole),
		input.IPAddress,
		input.UserAgent,
		DefaultSessionDuration,
	)
	if err != nil {
		asInteractor.logger.ErrorContext(ctx, "Failed to create new session", slog.Any("err", err))
		return AddSessionResponse{}, ErrUnexpected
	}

	err = asInteractor.sessionRepository.CreateSession(ctx, *newSession)
	if err != nil {
		asInteractor.logger.ErrorContext(ctx, "Failed to save session", slog.Any("err", err))
		return AddSessionResponse{}, ErrUnexpected
	}

	asInteractor.logger.InfoContext(ctx, "Session created successfully",
		slog.String("user_id", response.UserID.String()),
		slog.String("session_id", newSession.ID.String()),
	)

	return AddSessionResponse{Session: *newSession}, nil
}
