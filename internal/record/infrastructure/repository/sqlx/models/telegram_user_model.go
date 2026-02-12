package models

import (
	"time"

	"github.com/google/uuid"
)

// TelegramUserModel represents the sqlx model for the telegram_users table.
type TelegramUserModel struct {
	ID          uuid.UUID `db:"id"`
	TelegramID  uint64    `db:"telegram_id"`
	AddedAt     time.Time `db:"added_at"`
	AddedByUser uuid.UUID `db:"added_by_user"`
}
