package application

import (
	"context"
	"errors"
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

func (asInteractor *AddSession) Execute(ctx context.Context, input AddSessionRequest) error {
	// Check whether the provided credentials valid
	response, err := asInteractor.userClient.VerifyCredentials(ctx, input.Username, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrUsernameAbsent):
			asInteractor.logger.InfoContext(ctx, "Username doesn't exist", slog.String("username", input.Username))
			return ErrInvalidCredentials
		case errors.Is(err, client.ErrPasswordMissmatch):
			asInteractor.logger.InfoContext(ctx, "Password didn't match", slog.String("username", input.Username))
			return ErrInvalidCredentials
		default:
			asInteractor.logger.ErrorContext(ctx, "Unexpected error during credential verification")
			return ErrUnexpected
		}
	}
	newSession, err := domain.NewSession(
		response.UserID,
		domain.UserRole(response.UserRole),
		input.IpAddress,
		input.UserAgent,
		"",
		0,
	)
	asInteractor.sessionRepository.CreateSession(ctx, *newSession)
	return nil
}

// TODO add Execute implementation to this interactor
