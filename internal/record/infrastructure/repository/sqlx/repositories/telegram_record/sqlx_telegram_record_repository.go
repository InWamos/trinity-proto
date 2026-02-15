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
	"github.com/jackc/pgx/v5/pgconn"
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
) (*[]domain.TelegramRecord, error) {
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
	domainRecords := make([]domain.TelegramRecord, len(records))
	for i, record := range records {
		domainRecords[i] = repo.sqlxMapper.ToDomain(record)
	}

	return &domainRecords, nil
}

func (repo *SQLXTelegramRecordRepository) CreateTelegramRecord(
	ctx context.Context,
	telegramRecord domain.TelegramRecord,
) error {
	repo.logger.DebugContext(
		ctx,
		"Started CreateTelegramRecord request",
		slog.String("record_id", telegramRecord.ID.String()),
	)
	recordModel := repo.sqlxMapper.ToModel(telegramRecord)
	query := `INSERT INTO "records"."telegram_records" (id, from_user_telegram_id, in_telegram_chat_id, message_text, posted_at, added_at, added_by_user)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := repo.session.ExecContext(ctx, query,
		recordModel.ID,
		recordModel.FromUserTelegramID,
		recordModel.InTelegramChatID,
		recordModel.MessageText,
		recordModel.PostedAt,
		recordModel.AddedAt,
		recordModel.AddedByUser,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "fk_telegram_records_user":
				repo.logger.InfoContext(
					ctx,
					"User with this telegram id doesn't exist",
					slog.Uint64("user_telegram_id", telegramRecord.FromUserTelegramID),
				)
				return domain.ErrUnexistentTelegramUserReferenced
			case "unique_telegram_message_id":
				repo.logger.InfoContext(
					ctx,
					"This user has already been added this message",
					slog.Uint64("telegram_message_id", telegramRecord.MessageTelegramID),
					slog.String("added_by_user_id", telegramRecord.AddedByUser.String()),
				)
				return domain.ErrRecordAlreadyExists
			}
		}
		repo.logger.ErrorContext(ctx, "Failed to create telegram record", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	return nil
}

func (repo *SQLXTelegramRecordRepository) CreateTelegramRecords(
	ctx context.Context,
	telegramRecords []domain.TelegramRecord,
) error {
	repo.logger.DebugContext(
		ctx,
		"Started CreateTelegramRecords request",
		slog.Int("record_count", len(telegramRecords)),
	)
	query := `INSERT INTO "records"."telegram_records" (id, from_user_telegram_id, in_telegram_chat_id, message_text, posted_at, added_at, added_by_user)
	VALUES (:id, :from_user_telegram_id, :in_telegram_chat_id, :message_text, :posted_at, :added_at, :added_by_user)`
	recordModels := make([]models.SQLXTelegramRecordModel, len(telegramRecords))
	for i, record := range telegramRecords {
		recordModels[i] = repo.sqlxMapper.ToModel(record)
	}
	_, err := repo.session.NamedExecContext(ctx, query, recordModels)
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to create telegram records", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	return nil
}
