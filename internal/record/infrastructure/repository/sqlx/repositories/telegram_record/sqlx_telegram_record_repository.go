package repositories

import (
	"context"
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
	query := `SELECT tr.id, tr.message_telegram_id, tr.from_user_telegram_id, tr.in_telegram_chat_id, tr.message_text, tr.posted_at, tr.added_at, tr.added_by_user
	FROM "records"."telegram_records" tr
	JOIN "records"."telegram_users" tu ON tr.from_user_telegram_id = tu.id
	WHERE tu.telegram_id = $1 
	ORDER BY tr.posted_at DESC
	LIMIT 5`
	err := repo.session.SelectContext(ctx, &records, query, userTelegramID)
	if err != nil {
		repo.logger.ErrorContext(ctx, "Telegram record query failed", slog.Any("err", err))
		return nil, repository.ErrDatabaseFailed
	}

	// Check if no records were found
	if len(records) == 0 {
		repo.logger.InfoContext(ctx, "No telegram records found", slog.Uint64("user_telegram_id", userTelegramID))
		return nil, domain.ErrNoRecordsForThisTelegramID
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
	query := `INSERT INTO "records"."telegram_records" (id, message_telegram_id, from_user_telegram_id, in_telegram_chat_id, message_text, posted_at, added_at, added_by_user)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := repo.session.ExecContext(ctx, query,
		recordModel.ID,
		recordModel.MessageTelegramID,
		recordModel.FromTelegramUserID,
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
					slog.String("user_telegram_id", telegramRecord.FromTelegramUserID.String()),
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
	query := `INSERT INTO "records"."telegram_records" (id, message_telegram_id, from_user_telegram_id, in_telegram_chat_id, message_text, posted_at, added_at, added_by_user)
	VALUES (:id, :message_telegram_id, :from_user_telegram_id, :in_telegram_chat_id, :message_text, :posted_at, :added_at, :added_by_user)`
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
