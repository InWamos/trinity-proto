package service

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/go-playground/validator/v10"
)

type TelegramRecordService struct {
	validator *validator.Validate
}

func NewTelegramRecordService(validator *validator.Validate) *TelegramRecordService {
	return &TelegramRecordService{validator: validator}
}

func (service *TelegramRecordService) Validate(telegramRecord domain.TelegramRecord) error {
	err := service.validator.Struct(telegramRecord)
	if err != nil {
		return domain.ErrValidationFailed
	}
	return nil
}
