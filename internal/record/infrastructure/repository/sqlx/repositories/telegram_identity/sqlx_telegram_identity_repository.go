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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramIdentityRepository struct {
	session    *sqlx.Tx
	sqlxMapper *mappers.SqlxTelegramIdentityMapper
	logger     *slog.Logger
}

func NewSQLXTelegramIdentityRepository(
	session *sqlx.Tx,
	sqlxMapper *mappers.SqlxTelegramIdentityMapper,
	logger *slog.Logger,
) repository.TelegramIdentityRepository {
	tirLogger := logger.With(
		slog.String("component", "repository"),
		slog.String("name", "sqlx_telegram_identity_repository"),
	)
	return &SQLXTelegramIdentityRepository{
		session:    session,
		sqlxMapper: sqlxMapper,
		logger:     tirLogger,
	}
}

func (repo *SQLXTelegramIdentityRepository) AddIdentity(ctx context.Context, identity *domain.TelegramIdentity) error {
	repo.logger.DebugContext(ctx, "Started AddIdentity request", slog.String("identity_id", identity.ID.String()))
	identityModel := repo.sqlxMapper.ToModel(*identity)
	query := `INSERT INTO "records"."telegram_identities" (id, user_id, first_name, last_name, username, phone_number, bio, added_at, added_by_user)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := repo.session.ExecContext(ctx, query,
		identityModel.ID,
		identityModel.UserID,
		identityModel.FirstName,
		identityModel.LastName,
		identityModel.Username,
		identityModel.PhoneNumber,
		identityModel.Bio,
		identityModel.AddedAt,
		identityModel.AddedByUser,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "fk_telegram_identities_user":
				repo.logger.InfoContext(
					ctx,
					"Telegram user with that user id not found",
					slog.String("telegram_user_id", identityModel.UserID.String()),
					slog.String("constraint_name", pgErr.ConstraintName),
				)
				return domain.ErrUnexistentTelegramUserReferenced
			case "unique_telegram_identity_per_user":
				repo.logger.InfoContext(
					ctx,
					"This identity already exists",
					slog.String("telegram_user_id", identityModel.ID.String()),
					slog.String("constraint_name", pgErr.ConstraintName),
				)
				return domain.ErrIdentityAlreadyExists
			}
		}
		repo.logger.ErrorContext(ctx, "Failed to add telegram identity", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	return nil
}

func (repo *SQLXTelegramIdentityRepository) RemoveIdentityByID(ctx context.Context, identityID uuid.UUID) error {
	repo.logger.DebugContext(ctx, "Started RemoveIdentityByID request", slog.String("identity_id", identityID.String()))
	query := `DELETE FROM "records"."telegram_identities" WHERE id = $1`
	result, err := repo.session.ExecContext(ctx, query, identityID)
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to delete telegram identity", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.ErrorContext(ctx, "Failed to get rows affected", slog.Any("err", err))
		return repository.ErrDatabaseFailed
	}
	if rowsAffected == 0 {
		repo.logger.InfoContext(
			ctx,
			"Telegram identity not found for deletion",
			slog.String("identity_id", identityID.String()),
		)
		return repository.ErrIdentityNotFound
	}
	return nil
}

func (repo *SQLXTelegramIdentityRepository) GetIdentityByID(
	ctx context.Context,
	identityID uuid.UUID,
) (*domain.TelegramIdentity, error) {
	repo.logger.DebugContext(ctx, "Started GetIdentityByID request", slog.String("identity_id", identityID.String()))
	var identityModel models.TelegramIdentityModel
	query := `SELECT id, first_name, last_name, username, phone_number, bio, added_at, added_by_user
	FROM "records"."telegram_identities" WHERE id = $1`
	err := repo.session.GetContext(ctx, &identityModel, query, identityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.InfoContext(ctx, "Telegram identity not found", slog.String("identity_id", identityID.String()))
			return nil, repository.ErrIdentityNotFound
		}
		repo.logger.ErrorContext(ctx, "Failed to get telegram identity", slog.Any("err", err))
		return nil, repository.ErrDatabaseFailed
	}
	identity := repo.sqlxMapper.ToDomain(identityModel)
	return &identity, nil
}
