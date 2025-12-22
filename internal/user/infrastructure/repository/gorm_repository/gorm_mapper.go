package gormrepository

import (
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/models"
	"gorm.io/gorm"
)

type GormMapper struct {
}

func NewGormMapper() *GormMapper {
	return &GormMapper{}
}

func (gm *GormMapper) ToDomain(inputModel *models.UserModel) domain.User {
	return *domain.NewUser(inputModel.ID, inputModel.Username,
		inputModel.DisplayName, inputModel.PasswordHash,
		domain.Role(inputModel.Role))
}

func (gm *GormMapper) ToModel(inputEntity *domain.User) models.UserModel {
	return models.UserModel{
		ID:           inputEntity.ID,
		Username:     inputEntity.Username,
		DisplayName:  inputEntity.Username,
		PasswordHash: inputEntity.PasswordHash,
		Role:         models.Role(inputEntity.Role),
		CreatedAt:    inputEntity.CreatedAt,
		DeletedAt:    gorm.DeletedAt(inputEntity.DeletedAt),
	}
}
