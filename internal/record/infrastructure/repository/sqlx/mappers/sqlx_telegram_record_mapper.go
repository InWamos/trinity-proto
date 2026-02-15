package mappers

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
)

type SqlxTelegramRecordMapper struct{}

func NewSqlxTelegramRecordMapper() *SqlxTelegramRecordMapper {
	return &SqlxTelegramRecordMapper{}
}

func (sm *SqlxTelegramRecordMapper) ToDomain(inputModel models.SQLXTelegramRecordModel) domain.TelegramRecord {
	return domain.TelegramRecord{
		ID:                 inputModel.ID,
		MessageTelegramID:  inputModel.MessageTelegramID,
		FromTelegramUserID: inputModel.FromTelegramUserID,
		InTelegramChatID:   inputModel.InTelegramChatID,
		MessageText:        inputModel.MessageText,
		PostedAt:           inputModel.PostedAt,
		AddedAt:            inputModel.AddedAt,
		AddedByUser:        inputModel.AddedByUser,
	}
}

func (sm *SqlxTelegramRecordMapper) ToModel(inputEntity domain.TelegramRecord) models.SQLXTelegramRecordModel {
	return models.SQLXTelegramRecordModel{
		ID:                 inputEntity.ID,
		MessageTelegramID:  inputEntity.MessageTelegramID,
		FromTelegramUserID: inputEntity.FromTelegramUserID,
		InTelegramChatID:   inputEntity.InTelegramChatID,
		MessageText:        inputEntity.MessageText,
		PostedAt:           inputEntity.PostedAt,
		AddedAt:            inputEntity.AddedAt,
		AddedByUser:        inputEntity.AddedByUser,
	}
}
