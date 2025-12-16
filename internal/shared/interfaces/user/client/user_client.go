package client

import "context"

type UserClient interface {
	VerifyCredentials(ctx context.Context, username, password string) error
}
