package sqlxrepository

import (
	"database/sql"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/models"
)

type SqlxMapper struct{}

func NewSqlxMapper() *SqlxMapper {
	return &SqlxMapper{}
}

func (sm *SqlxMapper) ToDomain(inputModel *models.UserModelSqlx) domain.User {
	return *domain.NewUser(
		inputModel.ID,
		inputModel.Username,
		inputModel.DisplayName,
		inputModel.PasswordHash,
		domain.Role(inputModel.UserRole),
	)
}

func (sm *SqlxMapper) ToModel(inputEntity *domain.User) models.UserModelSqlx {
	deletedAt := sql.NullTime{} //nolint:exhaustruct // Required fields set after
	if inputEntity.DeletedAt.Valid {
		deletedAt.Time = inputEntity.DeletedAt.Time
		deletedAt.Valid = true
	}

	return models.UserModelSqlx{
		ID:           inputEntity.ID,
		Username:     inputEntity.Username,
		DisplayName:  inputEntity.DisplayName,
		PasswordHash: inputEntity.PasswordHash,
		UserRole:     models.UserRole(inputEntity.Role),
		CreatedAt:    inputEntity.CreatedAt,
		DeletedAt:    deletedAt,
	}
}
