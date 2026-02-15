package service

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/go-playground/validator/v10"
)

type TelegramModelValidator struct {
	validator *validator.Validate
}

func NewTelegramModelValidator(validator *validator.Validate) *TelegramModelValidator {
	return &TelegramModelValidator{validator: validator}
}

func (service *TelegramModelValidator) Validate(model any) error {
	err := service.validator.Struct(model)
	if err != nil {
		return domain.ErrValidationFailed
	}
	return nil
}
