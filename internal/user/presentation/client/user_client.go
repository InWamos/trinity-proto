package client

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/user/client"
	"github.com/InWamos/trinity-proto/internal/user/application"
)

type UserClient struct {
	validateUserCredentialsInteractor *application.ValidateUserCredentials
	logger                            *slog.Logger
}

func NewUserClient(
	validateUserCredentialsInteractor *application.ValidateUserCredentials,
	logger *slog.Logger,
) client.UserClient {
	ucLogger := logger.With(slog.String("component", "user_client"))
	return &UserClient{validateUserCredentialsInteractor: validateUserCredentialsInteractor, logger: ucLogger}
}

func (uClient *UserClient) VerifyCredentials(ctx context.Context, username, password string) error {
	interactorRequest := application.ValidateUserCredentialsRequest{Username: username, Password: password}
	if err := uClient.validateUserCredentialsInteractor.Execute(ctx, interactorRequest); err != nil {
		switch {
		case errors.Is(err, application.ErrUsernameAbsent):
			uClient.logger.InfoContext(ctx, "login attempt with non-existent username",
				slog.String("username", username))
			return client.ErrUsernameAbsent

		case errors.Is(err, application.ErrPasswordMismatch):
			uClient.logger.InfoContext(ctx, "invalid password attempt",
				slog.String("username", username))
			return client.ErrPasswordMissmatch

		case errors.Is(err, application.ErrDatabaseFailed):
			uClient.logger.ErrorContext(ctx, "database error during credential verification",
				slog.Any("err", err))
			return client.ErrUnexpectedError

		default:
			uClient.logger.ErrorContext(ctx, "unexpected error during credential verification",
				slog.Any("err", err))
			return client.ErrUnexpectedError
		}
	}
	return nil
}
