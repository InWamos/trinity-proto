package service

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/go-playground/validator/v10"
)

type TelegramUserService struct {
	validator *validator.Validate
}

func NewTelegramUserService(validator *validator.Validate) *TelegramUserService {
	return &TelegramUserService{validator: validator}
}

func (service *TelegramUserService) Validate(telegramUser domain.TelegramUser) error {
	err := service.validator.Struct(telegramUser)
	if err != nil {
		return domain.ErrValidationFailed
	}
	return nil
}
