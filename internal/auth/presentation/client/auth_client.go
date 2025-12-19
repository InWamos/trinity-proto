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

func (ac *AuthClient) ValidateSession(ctx context.Context, token string) (client.UserIdentity, error) {
	interactorRequest := application.VerifySessionRequest{Session_id: token}
	session, err := ac.verifySessionInteractor.Execute(ctx, interactorRequest)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrSessionNotFound):
			ac.logger.InfoContext(ctx, "session not found or expired",
				slog.String("token", token))
			return client.UserIdentity{}, client.ErrSessionInvalid

		case errors.Is(err, application.ErrSessionExpired):
			ac.logger.InfoContext(ctx, "session has expired",
				slog.String("token", token))
			return client.UserIdentity{}, client.ErrSessionExpired

		case errors.Is(err, application.ErrSessionRevoked):
			ac.logger.InfoContext(ctx, "session has been revoked",
				slog.String("token", token))
			return client.UserIdentity{}, client.ErrSessionRevoked

		default:
			ac.logger.ErrorContext(ctx, "unexpected error during session verification",
				slog.Any("err", err))
			return client.UserIdentity{}, client.ErrUnexpectedError
		}
	}

	ac.logger.DebugContext(ctx, "Session successfully verified",
		slog.String("token", token),
		slog.String("user_id", session.UserID.String()),
		slog.String("user_role", string(session.UserRole)))

	return client.UserIdentity{
		UserID:   session.UserID,
		UserRole: client.UserRole(session.UserRole),
	}, nil
}
