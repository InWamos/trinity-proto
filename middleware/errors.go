package middleware

import "errors"

var (
	ErrMissingToken = errors.New("missing authentication token")
)
