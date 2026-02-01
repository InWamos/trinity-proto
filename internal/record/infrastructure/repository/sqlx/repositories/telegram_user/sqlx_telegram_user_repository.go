package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramUserRepository struct {
	session    *sqlx.Tx
	sqlxMapper *mappers.SqlxTelegramUserMapper
	logger     *slog.Logger
}

func NewSQLXTelegramUserRepository(
	session *sqlx.Tx,
	sqlxMapper *mappers.SqlxTelegramUserMapper,
	logger *slog.Logger,
) repository.TelegramUserRepository {
	trrLogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "sqlx_telegram_user_repository"),
	)
	return &SQLXTelegramUserRepository{
		session:    session,
		sqlxMapper: sqlxMapper,
		logger:     trrLogger,
	}
}

func (repo *SQLXTelegramUserRepository) GetByTelegramID(
	ctx context.Context,
	telegramID uint64,
) (*domain.TelegramUser, error) {
	repo.logger.DebugContext(ctx, "Started GetByTelegramID request", slog.Uint64("telegram_id", telegramID))
	var userModel models.TelegramUserModel
	query := `SELECT id, telegram_id, telegram_user_identity_id, added_at, added_by_user
	FROM "records"."telegram_users" WHERE telegram_id = $1`
	err := repo.session.GetContext(ctx, &userModel, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.InfoContext(ctx, "Telegram user not found", slog.Uint64("telegram_id", telegramID))
			return nil, domain.ErrUserNotFound
		}
		repo.logger.ErrorContext(ctx, "Failed to get telegram user", slog.Any("err", err))
		return nil, repository.ErrDatabaseFailed
	}
	user := repo.sqlxMapper.ToDomain(userModel)
	return &user, nil
}

func (repo *SQLXTelegramUserRepository) AddUser(ctx context.Context, user *domain.TelegramUser) error {
	repo.logger.DebugContext(ctx, "Started AddUser request", slog.String("user_id", user.ID.String()))
	userModel := repo.sqlxMapper.ToModel(*user)
	query := `INSERT INTO "records"."telegram_users" (id, telegram_id, telegram_user_identity_id, added_at, added_by_user)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.session.ExecContext(ctx, query,
		userModel.ID,
		userModel.TelegramID,
		userModel.TelegramIdentityID,
		userModel.AddedAt,
		userModel.AddedByUser,
	)
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to add telegram user", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	return nil
}

func (repo *SQLXTelegramUserRepository) DeleteUserByTelegramID(ctx context.Context, telegramID uint64) error {
	repo.logger.DebugContext(ctx, "Started DeleteUserByTelegramID request", slog.Uint64("telegram_id", telegramID))
	query := `DELETE FROM "records"."telegram_users" WHERE telegram_id = $1`
	result, err := repo.session.ExecContext(ctx, query, telegramID)
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to delete telegram user", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to get rows affected", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	if rowsAffected == 0 {
		repo.logger.InfoContext(ctx, "Telegram user not found for deletion", slog.Uint64("telegram_id", telegramID))
		return domain.ErrUserNotFound
	}
	return nil
}
