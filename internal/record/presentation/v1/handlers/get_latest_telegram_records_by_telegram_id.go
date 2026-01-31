package handlers

import (
	"log/slog"

	application "github.com/InWamos/trinity-proto/internal/record/application/telegram"
)

type GetLatestTelegramRecordsByTelegramIDHandler struct {
	interactor *application.GetLatestTelegramRecordsByUserTelegramID
	logger     *slog.Logger
}

func NewGetLatestTelegramRecordsByTelegramID(
	interactor *application.GetLatestTelegramRecordsByUserTelegramID,
	logger *slog.Logger,
) *GetLatestTelegramRecordsByTelegramIDHandler {
	handlerLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "get_latest_telegram_records_by_user_id"),
	)

	return &GetLatestTelegramRecordsByTelegramIDHandler{
		interactor: interactor,
		logger:     handlerLogger,
	}
}
