package models

import (
	"time"

	"github.com/google/uuid"
)

type SQLXTelegramRecordModel struct {
	ID                 uuid.UUID `db:"id"`
	MessageTelegramID  uint64    `db:"message_telegram_id"`
	FromUserTelegramID uint64    `db:"from_user_telegram_id"`
	InTelegramChatID   int64     `db:"in_telegram_chat_id"`
	MessageText        string    `db:"message_text"`
	PostedAt           time.Time `db:"posted_at"`
	AddedAt            time.Time `db:"added_at"`
	AddedByUser        uuid.UUID `db:"added_by_user"`
}
