package client

import "context"

type AuthClient interface {
	ValidateSession(ctx context.Context, session string) error
}
