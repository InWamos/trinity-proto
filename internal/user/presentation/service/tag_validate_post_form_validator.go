package service

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type TagValidatePostFormValidator struct {
	validator *validator.Validate
}

func NewTagValidatePostFormValidator(validator *validator.Validate) PostFormValidator {
	return &TagValidatePostFormValidator{
		validator: validator,
	}
}

func (val *TagValidatePostFormValidator) ValidateBody(body io.ReadCloser, form any) error {
	if err := json.NewDecoder(body).Decode(&form); err != nil {
		return err
	}
	if err := val.validator.Struct(form); err != nil {
		return err
	}
	return nil
}
