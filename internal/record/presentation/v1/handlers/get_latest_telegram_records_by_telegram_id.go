package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	application "github.com/InWamos/trinity-proto/internal/record/application/telegram"
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
)

// GetLatestTelegramRecordsByTelegramIDResponse represents the response from the GetLatestTelegramRecordsByTelegramID endpoint
//
//	@Description	information response
type GetLatestTelegramRecordsByTelegramIDResponse struct {
	TelegramID uint64                  `json:"telegram_id" example:"428736582143"`
	Records    []domain.TelegramRecord `json:"records"`
}
type GetLatestTelegramRecordsByTelegramIDRequest struct {
	TelegramID uint64 `json:"telegram_id" example:"28736582143"`
}

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

// ServeHTTP handles an HTTP request to get the latest Telegram records by Telegram ID.
//
//	@Summary		Get latest Telegram records
//	@Description	Get the latest Telegram records for a specific Telegram user ID.
//	@Tags			record
//	@Accept			json
//	@Produce		json
//	@Param			request	body		GetLatestTelegramRecordsByTelegramIDRequest		true	"Telegram ID request"
//	@Success		200		{object}	GetLatestTelegramRecordsByTelegramIDResponse	"Latest records retrieved successfully"
//	@Failure		400		"Invalid telegram ID format"
//	@Failure		403		"Insufficient privileges"
//	@Failure		404		"Telegram ID not found"
//	@Failure		500		"Internal server error"
//	@Router			/v1/record/telegram/{telegram_id}/records [get]
func (handler *GetLatestTelegramRecordsByTelegramIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req GetLatestTelegramRecordsByTelegramIDRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		handler.logger.DebugContext(r.Context(), "invalid telegram ID format", slog.Any("err", err))
		http.Error(w, "Invalid telegram ID format", http.StatusBadRequest)
		return
	}
	requestDTO := application.GetLatestTelegramRecordsByUserTelegramIDRequest{UserTelegramID: req.TelegramID}
	resp, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		switch {
		case errors.Is(err, rbac.ErrInsufficientPrivileges):
			handler.logger.DebugContext(r.Context(), "Auth error", slog.Any("err", err))
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		case errors.Is(err, domain.ErrNoRecordsForThisTelegramID):
			handler.logger.DebugContext(r.Context(), "telegram user not found by ID", slog.Any("err", err))
			http.Error(w, "Telegram ID not found", http.StatusNotFound)
			return
		default:
			handler.logger.DebugContext(r.Context(), "Database error", slog.Any("err", err))
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}
	responce := GetLatestTelegramRecordsByTelegramIDResponse{TelegramID: req.TelegramID, Records: *resp.TelegramRecords}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responce) //nolint:musttag // Linter error, struct contains json tags
}
