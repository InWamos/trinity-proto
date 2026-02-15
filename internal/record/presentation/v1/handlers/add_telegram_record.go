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
	"github.com/google/uuid"
)

// AddTelegramRecordRequest represents the request payload for adding a new telegram record.
type AddTelegramRecordRequest struct {
	MessageTelegramID  uint64    `json:"message_telegram_id"   example:"28736582143"`
	FromUserTelegramID uuid.UUID `json:"from_user_telegram_id" example:"cf6e273b-ac6e-43f1-abba-d8009ffc1b3f"`
	InTelegramChatID   int64     `json:"in_telegram_chat_id"   example:"123456789"`
	MessageText        string    `json:"message_text"          example:"Hello world!"`
	PostedAt           time.Time `json:"posted_at"             example:"2024-01-15T10:30:00Z"`
}

// AddTelegramRecordResponse represents the response payload after successfully adding a telegram record.
type AddTelegramRecordResponse struct {
	RecordID string `json:"record_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// AddTelegramRecordHandler handles HTTP requests for adding telegram records.
type AddTelegramRecordHandler struct {
	interactor *application.AddTelegramRecord
	logger     *slog.Logger
}

// NewAddTelegramRecordHandler creates a new instance of AddTelegramRecordHandler.
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

// ServeHTTP handles POST requests to add a new telegram record.
//
//	@Summary		Add a new telegram record
//	@Description	Creates a new telegram record with the provided message details
//	@Tags			record
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AddTelegramRecordRequest	true	"Record details"
//	@Success		201		{object}	AddTelegramRecordResponse
//	@Failure		400		{string}	string	"Invalid request format"
//	@Failure		403		{string}	string	"Insufficient privileges"
//	@Failure		409		{string}	string	"Record already exists or user not found"
//	@Failure		422		{string}	string	"Record contains unprocessable fields"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/v1/record/telegram [post]
//	@Security		Bearer
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
