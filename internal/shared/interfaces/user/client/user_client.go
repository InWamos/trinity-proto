package client

import (
	"context"
	"errors"
)

var (
	ErrUsernameAbsent    = errors.New("username record absent")
	ErrPasswordMissmatch = errors.New("password missmatch")
	ErrUnexpectedError   = errors.New("unexpected error occured")
)

type UserClient interface {
	VerifyCredentials(ctx context.Context, username, password string) error
}
