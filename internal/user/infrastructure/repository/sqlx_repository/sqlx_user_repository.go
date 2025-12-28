package sqlxrepository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/models"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SqlxUserRepository struct {
	session    *sqlx.Tx
	sqlxMapper *SqlxMapper
	logger     *slog.Logger
}

func NewSqlxUserRepository(session *sqlx.Tx, logger *slog.Logger) repository.UserRepository {
	urlogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "sqlx_user_repository"),
	)
	return &SqlxUserRepository{
		session:    session,
		sqlxMapper: NewSqlxMapper(),
		logger:     urlogger,
	}
}

func (ur *SqlxUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ur.logger.DebugContext(ctx, "Started GetUserByID request")

	var user models.UserModelSqlx
	query := `SELECT id, username, display_name, password_hash, user_role, created_at, deleted_at 
			  FROM "user".users WHERE id = $1 AND deleted_at IS NULL`

	err := ur.session.GetContext(ctx, &user, query, id)
	ur.logger.DebugContext(ctx, "Finished GetUserByID request")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ur.logger.InfoContext(ctx, "User not found by id", slog.String("user_id", id.String()))
			return domain.User{}, repository.ErrUserNotFound
		}
		ur.logger.ErrorContext(
			ctx,
			"Failed to find user by id",
			slog.String("user_id", id.String()),
			slog.Any("err", err),
		)
		return domain.User{}, err
	}

	return ur.sqlxMapper.ToDomain(&user), nil
}

func (ur *SqlxUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	ur.logger.DebugContext(ctx, "Started GetUserByUsername request")

	var user models.UserModelSqlx
	query := `SELECT id, username, display_name, password_hash, user_role, created_at, deleted_at 
			  FROM "user".users WHERE username = $1 AND deleted_at IS NULL`

	err := ur.session.GetContext(ctx, &user, query, username)
	ur.logger.DebugContext(ctx, "Finished GetUserByUsername request")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ur.logger.InfoContext(ctx, "User not found by username", slog.String("user_username", username))
			return domain.User{}, repository.ErrUserNotFound
		}
		ur.logger.ErrorContext(
			ctx,
			"Failed to find user by username",
			slog.String("user_username", username),
			slog.Any("err", err),
		)
		return domain.User{}, err
	}

	return ur.sqlxMapper.ToDomain(&user), nil
}

func (ur *SqlxUserRepository) RemoveUserByID(ctx context.Context, id uuid.UUID) error {
	ur.logger.DebugContext(ctx, "Started RemoveUserByID request")

	// Soft delete
	query := `UPDATE "user".users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := ur.session.ExecContext(ctx, query, id)

	ur.logger.DebugContext(ctx, "Finished RemoveUserByID request")

	if err != nil {
		ur.logger.ErrorContext(
			ctx,
			"Failed to remove user by id",
			slog.String("user_id", id.String()),
			slog.Any("err", err),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ur.logger.ErrorContext(ctx, "Failed to get rows affected", slog.Any("err", err))
		return err
	}

	if rowsAffected == 0 {
		ur.logger.InfoContext(ctx, "User not found by id", slog.String("id", id.String()))
		return repository.ErrUserNotFound
	}

	ur.logger.DebugContext(
		ctx,
		"User has been removed",
		slog.String("user_id", id.String()),
		slog.Int64("rows_affected", rowsAffected),
	)
	return nil
}

func (ur *SqlxUserRepository) ChangeUserRoleByID(ctx context.Context, id uuid.UUID, changeToRole domain.Role) error {
	ur.logger.DebugContext(ctx, "Started ChangeUserRoleByID request")

	query := `UPDATE "user".users SET user_role = $2 WHERE id = $1 AND deleted_at IS NULL`
	result, err := ur.session.ExecContext(ctx, query, id, changeToRole)

	ur.logger.DebugContext(ctx, "Finished ChangeUserRoleByID request")

	if err != nil {
		ur.logger.ErrorContext(
			ctx,
			"Failed to change user role by id",
			slog.String("user_id", id.String()),
			slog.Any("err", err),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ur.logger.ErrorContext(ctx, "Failed to get rows affected", slog.Any("err", err))
		return err
	}

	if rowsAffected == 0 {
		ur.logger.InfoContext(ctx, "User not found by id", slog.String("user_id", id.String()))
		return repository.ErrUserNotFound
	}

	ur.logger.DebugContext(
		ctx,
		"User role has been updated",
		slog.String("user_id", id.String()),
		slog.Int64("rows_affected", rowsAffected),
	)
	return nil
}

func (ur *SqlxUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	ur.logger.DebugContext(ctx, "Started CreateUser request")

	userModel := ur.sqlxMapper.ToModel(&user)
	query := `INSERT INTO "user".users (id, username, display_name, password_hash, user_role, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := ur.session.ExecContext(
		ctx,
		query,
		userModel.ID,
		userModel.Username,
		userModel.DisplayName,
		userModel.PasswordHash,
		userModel.UserRole,
		userModel.CreatedAt,
	)

	ur.logger.DebugContext(ctx, "Finished CreateUser request")

	if err != nil {
		ur.logger.ErrorContext(ctx, "Failed to save user record", slog.Any("err", err))
		return repository.ErrUserCreationFailed
	}

	ur.logger.DebugContext(ctx, "User has been created", slog.String("user_id", userModel.ID.String()))
	return nil
}
