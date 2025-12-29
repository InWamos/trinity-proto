package domain

import "github.com/google/uuid"

type TelegramAttachment struct {
	ID                   uuid.UUID
	TelegramAttachmentID uint64
	FileName             string
	StorageKey           string
	FileSizeType         uint64
	MimeType             string
	AddedByUser          uuid.UUID
}
