package gormrepository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/models"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	session    *gorm.DB
	gormMapper *GormMapper
	logger     *slog.Logger
}

func NewGormUserRepository(session *gorm.DB, logger *slog.Logger) repository.UserRepository {
	urlogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "gorm_user_repository"),
	)
	return &GormUserRepository{session: session, logger: urlogger}
}

func (ur *GormUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := gorm.G[models.UserModel](ur.session).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.InfoContext(ctx, "User not found by id", slog.String("user_id", id.String()))
			return domain.User{}, repository.ErrUserNotFound
		}
	}
	return ur.gormMapper.ToDomain(&user), nil
}

func (ur *GormUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := gorm.G[models.UserModel](ur.session).Where("username = ?", username).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.InfoContext(ctx, "User not found by username", slog.String("user_username", username))
			return domain.User{}, repository.ErrUserNotFound
		}
		ur.logger.ErrorContext(
			ctx,
			"Failed to find user by username",
			slog.String("user_username", username),
			slog.Any("err", err),
		)
		return domain.User{}, repository.ErrUserNotFound
	}
	return ur.gormMapper.ToDomain(&user), nil
}

func (ur *GormUserRepository) RemoveUserByID(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := gorm.G[models.UserModel](ur.session).Where("id = ?", id).Delete(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.InfoContext(ctx, "User not found by id", slog.String("id", id.String()))
			return repository.ErrUserNotFound
		}
		ur.logger.ErrorContext(
			ctx,
			"Failed to find user by id",
			slog.String("user_id", id.String()),
			slog.Any("err", err),
		)
		return err
	}
	ur.logger.DebugContext(
		ctx,
		"User has been removed",
		slog.String("user_id", id.String()),
		slog.Int("rows_affected", rowsAffected),
	)
	return nil
}

func (ur *GormUserRepository) ChangeUserRoleByID(ctx context.Context, id uuid.UUID, changeToRole domain.Role) error {
	rowsAffected, err := gorm.G[models.UserModel](ur.session).Where("id = ?", id).Update(ctx, "role", changeToRole)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.InfoContext(ctx, "User not found by id", slog.String("user_id", id.String()))
			return repository.ErrUserNotFound
		}
		ur.logger.ErrorContext(
			ctx,
			"Failed to find user by id",
			slog.String("user_id", id.String()),
			slog.Any("err", err),
		)
		return err
	}
	ur.logger.DebugContext(
		ctx,
		"User role has been updated",
		slog.String("user_id", id.String()),
		slog.Int("rows_affected", rowsAffected),
	)
	return nil
}

func (ur *GormUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	userModel := ur.gormMapper.ToModel(&user)
	err := gorm.G[models.UserModel](ur.session).Create(ctx, &userModel)
	if err != nil {
		ur.logger.ErrorContext(ctx, "Failed to save user record", slog.Any("err", err))
		return repository.ErrUserCreationFailed
	}
	ur.logger.DebugContext(ctx, "User has been created", slog.String("user_id", userModel.ID.String()))
	return nil
}
