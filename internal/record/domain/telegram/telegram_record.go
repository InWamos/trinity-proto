package domain

import (
	"time"

	"github.com/google/uuid"
)

// TelegramRecord to see all possible parseable fields of a message, see
// this beautiful implementation of MTProto client
// https://github.com/KurimuzonAkuma/kurigram/blob/dev/pyrogram/types/messages_and_media/message.py
type TelegramRecord struct {
	ID          uuid.UUID `validate:"required,uuid"`
	Author      uuid.UUID `validate:"required,uuid"`
	MessageText string    `validate:"required,min=1,max=4096"`
	Attachments []uuid.UUID
	PostedAt    time.Time `validate:"required"`
	AddedAt     time.Time `validate:"required"`
	AddedByUser uuid.UUID `validate:"required,uuid"`
}
