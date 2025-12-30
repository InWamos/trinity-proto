package repositories

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramRecordRepository struct {
	session    *sqlx.Tx
	sqlxMapper *mappers.SqlxTelegramRecordMapper
	logger     *slog.Logger
}

func NewSQLXTelegramRecordRepository(session *sqlx.Tx, sqlxMapper *mappers.SqlxTelegramRecordMapper, logger *slog.Logger) repository.TelegramRecordRepository {
	trrlogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "sqlx_telegram_record_repository"),
	)
	return &SQLXTelegramRecordRepository{
		session:    session,
		sqlxMapper: sqlxMapper,
		logger:     trrlogger,
	}
}
