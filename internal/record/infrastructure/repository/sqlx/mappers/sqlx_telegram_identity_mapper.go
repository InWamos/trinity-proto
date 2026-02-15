package mappers

import (
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
)

type SqlxTelegramIdentityMapper struct{}

func NewSqlxTelegramIdentityMapper() *SqlxTelegramIdentityMapper {
	return &SqlxTelegramIdentityMapper{}
}

func (sm *SqlxTelegramIdentityMapper) ToDomain(inputModel models.TelegramIdentityModel) domain.TelegramIdentity {
	return domain.TelegramIdentity{
		ID:          inputModel.ID,
		UserID:      inputModel.UserID,
		FirstName:   inputModel.FirstName,
		LastName:    inputModel.LastName,
		Username:    inputModel.Username,
		PhoneNumber: inputModel.PhoneNumber,
		Bio:         inputModel.Bio,
		AddedAt:     inputModel.AddedAt,
		AddedByUser: inputModel.AddedByUser,
	}
}

func (sm *SqlxTelegramIdentityMapper) ToModel(inputEntity domain.TelegramIdentity) models.TelegramIdentityModel {
	return models.TelegramIdentityModel{
		ID:          inputEntity.ID,
		UserID:      inputEntity.UserID,
		FirstName:   inputEntity.FirstName,
		LastName:    inputEntity.LastName,
		Username:    inputEntity.Username,
		PhoneNumber: inputEntity.PhoneNumber,
		Bio:         inputEntity.Bio,
		AddedAt:     inputEntity.AddedAt,
		AddedByUser: inputEntity.AddedByUser,
	}
}
