package client

import (
	"context"
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

func (client *UserClient) VerifyCredentials(ctx context.Context, username, password string) error
