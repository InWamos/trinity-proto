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

type AddTelegramUserRequest struct {
	TelegramID uint64 `json:"telegram_id" example:"28736582143"`
}

type AddTelegramUserResponse struct {
	RecordID string `json:"record_id" example:"28736582143"`
}

type AddTelegramUserHandler struct {
	interactor *application.AddTelegramUser
	logger     *slog.Logger
}

func NewAddTelegramUserHandler(
	interactor *application.AddTelegramUser,
	logger *slog.Logger,
) *AddTelegramUserHandler {
	handlerLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "add_telegram_user_handler"),
	)

	return &AddTelegramUserHandler{
		interactor: interactor,
		logger:     handlerLogger,
	}
}

// ServeHTTP handles an HTTP request to add a Telegram user.
//
//	@Summary		Add a new Telegram user
//	@Description	Add a new Telegram user by Telegram ID. This creates a record linking a Telegram user to the system.
//	@Tags			record
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AddTelegramUserRequest	true	"Telegram user request"
//	@Success		201		{object}	AddTelegramUserResponse	"Telegram user added successfully"
//	@Failure		400		"Invalid request format"
//	@Failure		403		"Insufficient privileges"
//	@Failure		409		"You have already added this user"
//	@Failure		422		"This user contains unprocessable fields"
//	@Failure		500		"Internal server error"
//	@Router			/v1/record/telegram/user [post]
func (handler *AddTelegramUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req AddTelegramUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		handler.logger.DebugContext(r.Context(), "invalid request format", slog.Any("err", err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	requestDTO := application.AddTelegramUserRequest{TelegramID: req.TelegramID}
	resp, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		switch {
		case errors.Is(err, rbac.ErrInsufficientPrivileges):
			handler.logger.DebugContext(r.Context(), "Auth error", slog.Any("err", err))
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		case errors.Is(err, domain.ErrUserAlreadyExists):
			handler.logger.DebugContext(r.Context(), "This user is already added", slog.Any("err", err))
			http.Error(w, "You have already added this user", http.StatusConflict)
			return
		case errors.Is(err, domain.ErrValidationFailed):
			handler.logger.DebugContext(
				r.Context(),
				"Validation has failed",
				slog.Any("err", err),
			)
			http.Error(w, "This user contains unprocessable fields", http.StatusUnprocessableEntity)
			return
		default:
			handler.logger.DebugContext(r.Context(), "Database error", slog.Any("err", err))
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}

	response := AddTelegramUserResponse{RecordID: resp.UserID}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
