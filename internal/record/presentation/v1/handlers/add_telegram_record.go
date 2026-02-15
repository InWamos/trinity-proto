package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	application "github.com/InWamos/trinity-proto/internal/record/application/telegram/record"
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
)

type AddTelegramRecordRequest struct {
	MessageTelegramID  uint64    `json:"message_telegram_id"   example:"28736582143"`
	FromUserTelegramID uint64    `json:"from_user_telegram_id" example:"28736582143"`
	InTelegramChatID   int64     `json:"in_telegram_chat_id"   example:"123456789"`
	MessageText        string    `json:"message_text"          example:"Hello world!"`
	PostedAt           time.Time `json:"posted_at"             example:"2024-01-15T10:30:00Z"`
}

type AddTelegramRecordResponse struct {
	RecordID string `json:"record_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type AddTelegramRecordHandler struct {
	interactor *application.AddTelegramRecord
	logger     *slog.Logger
}

func NewAddTelegramRecordHandler(
	interactor *application.AddTelegramRecord,
	logger *slog.Logger,
) *AddTelegramRecordHandler {
	handlerLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "add_telegram_record_handler"),
	)

	return &AddTelegramRecordHandler{
		interactor: interactor,
		logger:     handlerLogger,
	}
}

func (handler *AddTelegramRecordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req AddTelegramRecordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		handler.logger.DebugContext(r.Context(), "invalid request format", slog.Any("err", err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	requestDTO := application.AddTelegramRecordRequest{
		MessageTelegramID:  req.MessageTelegramID,
		FromUserTelegramID: req.FromUserTelegramID,
		InTelegramChatID:   req.InTelegramChatID,
		MessageText:        req.MessageText,
		PostedAt:           req.PostedAt,
	}

	resp, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil { //nolint: dupl //Has to be a duplicate
		switch {
		case errors.Is(err, rbac.ErrInsufficientPrivileges):
			handler.logger.DebugContext(r.Context(), "Auth error", slog.Any("err", err))
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		case errors.Is(err, domain.ErrValidationFailed):
			handler.logger.DebugContext(
				r.Context(),
				"Validation has failed",
				slog.Any("err", err),
			)
			http.Error(w, "Record contains unprocessable fields", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrUnexistentTelegramUserReferenced):
			handler.logger.DebugContext(
				r.Context(),
				"Unexistent user referenced",
				slog.Any("err", err),
			)
			http.Error(w, "This record references user that hasn't been added yet", http.StatusConflict)
			return
		case errors.Is(err, domain.ErrRecordAlreadyExists):
			handler.logger.DebugContext(
				r.Context(),
				"This record already exists",
				slog.Any("err", err),
			)
			http.Error(w, "This record already exists", http.StatusConflict)
			return
		default:
			handler.logger.ErrorContext(r.Context(), "Database error", slog.Any("err", err))
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}

	response := AddTelegramRecordResponse{RecordID: resp.RecordID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
