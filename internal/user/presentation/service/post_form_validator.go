package service

import "io"

type PostFormValidator interface {
	ValidateBody(body io.ReadCloser, form any) error
}
