package client

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
)

type AuthClient struct {
	logger                  *slog.Logger
	verifySessionInteractor *application.VerifySession
}

func NewAuthClient(
	verifySessionInteractor *application.VerifySession,
	logger *slog.Logger,
) client.AuthClient {
	acLogger := logger.With(slog.String("component", "auth_client"))
	return &AuthClient{verifySessionInteractor: verifySessionInteractor, logger: acLogger}
}

func (ac *AuthClient) ValidateSession(ctx context.Context, token string) error {
	interactorRequest := application.VerifySessionRequest{Session_id: token}
	err := ac.verifySessionInteractor.Execute(ctx, interactorRequest)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrSessionNotFound):
			ac.logger.InfoContext(ctx, "session not found or expired",
				slog.String("token", token))
			return client.ErrSessionInvalid

		case errors.Is(err, application.ErrSessionExpired):
			ac.logger.InfoContext(ctx, "session has expired",
				slog.String("token", token))
			return client.ErrSessionExpired

		case errors.Is(err, application.ErrSessionRevoked):
			ac.logger.InfoContext(ctx, "session has been revoked",
				slog.String("token", token))
			return client.ErrSessionRevoked

		default:
			ac.logger.ErrorContext(ctx, "unexpected error during session verification",
				slog.Any("err", err))
			return client.ErrUnexpectedError
		}
	}

	ac.logger.DebugContext(ctx, "Session successfully verified",
		slog.String("token", token))
	return nil
}
