package mappers

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
)

type SqlxTelegramUserMapper struct{}

func NewSqlxTelegramUserMapper() *SqlxTelegramUserMapper {
	return &SqlxTelegramUserMapper{}
}

func (sm *SqlxTelegramUserMapper) ToDomain(inputModel models.TelegramUserModel) domain.TelegramUser {
	return domain.TelegramUser{
		ID:          inputModel.ID,
		TelegramID:  inputModel.TelegramID,
		AddedAt:     inputModel.AddedAt,
		AddedByUser: inputModel.AddedByUser,
	}
}

func (sm *SqlxTelegramUserMapper) ToModel(inputEntity domain.TelegramUser) models.TelegramUserModel {
	return models.TelegramUserModel{
		ID:          inputEntity.ID,
		TelegramID:  inputEntity.TelegramID,
		AddedAt:     inputEntity.AddedAt,
		AddedByUser: inputEntity.AddedByUser,
	}
}
