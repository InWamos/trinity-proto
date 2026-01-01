package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRecordAlreadyExists        = errors.New("record already exists")
	ErrNoRecordsForThisTelegramID = errors.New("no records for this telegram id")
)

// TelegramRecord to see all possible parseable fields of a message, see
// this beautiful implementation of MTProto client
// https://github.com/KurimuzonAkuma/kurigram/blob/dev/pyrogram/types/messages_and_media/message.py
type TelegramRecord struct {
	ID                 uuid.UUID `validate:"required,uuid"`
	FromUserTelegramID uint64    `validate:"required,gt=0,lte=300000000000"`
	InTelegramChatID   int64     `validate:"required"`
	MessageText        string    `validate:"required,max=4096"`
	PostedAt           time.Time `validate:"required"`
	AddedAt            time.Time `validate:"required"`
	AddedByUser        uuid.UUID `validate:"required,uuid"`
}
