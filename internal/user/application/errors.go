package application

import "errors"

var (
	ErrHashingFailed           = errors.New("password hashing failed")
	ErrUUIDGeneration          = errors.New("UUID generation failed")
	ErrUserWithIDAlreadyExists = errors.New("this uuid is already in the database")
	ErrDatabaseFailed          = errors.New("the database operation has failed")
	ErrUsernameAbsent          = errors.New("this username is absent")
	ErrPasswordMismatch        = errors.New("password didn't match")
	ErrNoUserIdentityProvided  = errors.New("no user identity provided")
)
