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

type SQLXTelegramRecordRepository struct {
	session    *sqlx.Tx
	sqlxMapper *mappers.SqlxTelegramRecordMapper
	logger     *slog.Logger
}

func NewSQLXTelegramRecordRepository(
	session *sqlx.Tx,
	sqlxMapper *mappers.SqlxTelegramRecordMapper,
	logger *slog.Logger,
) repository.TelegramRecordRepository {
	trrLogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "sqlx_telegram_record_repository"),
	)
	return &SQLXTelegramRecordRepository{
		session:    session,
		sqlxMapper: sqlxMapper,
		logger:     trrLogger,
	}
}

func (repo *SQLXTelegramRecordRepository) GetLatestTelegramRecordsByUserTelegramID(
	ctx context.Context,
	userTelegramID uint64,
) (*[]*domain.TelegramRecord, error) {
	repo.logger.DebugContext(ctx, "Started GetLatestTelegramRecordsByUserTelegramID request")
	var records []models.SQLXTelegramRecordModel
	query := `SELECT id, from_user_telegram_id, in_telegram_chat_id, message_text, posted_at, added_at, added_by_user
	FROM "records"."telegram_records" WHERE from_user_telegram_id = $1 LIMIT 5`
	err := repo.session.SelectContext(ctx, &records, query, userTelegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.InfoContext(ctx, "Telegram record not found", slog.Uint64("user_telegram_id", userTelegramID))
			return nil, domain.ErrNoRecordsForThisTelegramID
		}
		repo.logger.ErrorContext(ctx, "Telegram record failed", slog.Any("err", err))
		return nil, repository.ErrDatabaseFailed
	}
	domainRecords := make([]*domain.TelegramRecord, len(records))
	for i, record := range records {
		domainRecords[i] = repo.sqlxMapper.ToDomain(&record)
	}

	return &domainRecords, nil

}

func (repo *SQLXTelegramRecordRepository) CreateTelegramRecord(
	ctx context.Context,
	telegramRecord domain.TelegramRecord,
) error {
	return nil
}

func (repo *SQLXTelegramRecordRepository) CreateTelegramRecords(
	ctx context.Context,
	telegramRecords []domain.TelegramRecord,
) error {
	return nil
}
