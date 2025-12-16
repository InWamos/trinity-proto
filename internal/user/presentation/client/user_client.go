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
	return &UserClient{validateUserCredentialsInteractor: validateUserCredentialsInteractor, logger: logger}
}

func (uClient *UserClient) VerifyCredentials(ctx context.Context, username, password string) error {
	interactorRequest := application.ValidateUserCredentialsRequest{Username: username, Password: password}
	if err := uClient.validateUserCredentialsInteractor.Execute(ctx, interactorRequest); err != nil {
		if errors.Is(err, application.ErrPasswordMismatch) {
			return client.ErrPasswordMissmatch
		}
		if errors.Is(err, application.ErrUsernameAbsent) {
			return client.ErrUsernameAbsent
		}
		return client.ErrUnexpectedError
	}
	return nil
}
