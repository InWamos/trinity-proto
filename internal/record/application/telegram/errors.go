package application

import "errors"

var (
	ErrDatabaseFailed = errors.New("the database operation has failed")
)
