package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	application "github.com/InWamos/trinity-proto/internal/record/application/telegram/identity"
	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/google/uuid"
)

type AddTelegramIdentityRequest struct {
	UserID      uuid.UUID `json:"telegram_id"           example:"20d8a06c-2fac-4643-ba78-7da267a7fe51"`
	Username    string    `json:"telegram_username"     example:"user1235"`
	FirstName   string    `json:"telegram_first_name"   example:"John"`
	LastName    string    `json:"telegram_last_name"    example:"Doe"`
	Bio         string    `json:"telegram_bio"          example:"Hi! I am using Whatsapp"`
	PhoneNumber string    `json:"telegram_phone_number" example:"+11234567890 (Use e164 format)"`
}

type AddTelegramIdentityResponse struct {
	RecordID string `json:"record_id" example:"28736582143"`
}

type AddTelegramIdentityHandler struct {
	interactor *application.AddTelegramIdentity
	logger     *slog.Logger
}

func NewAddTelegramIdentityHandler(
	interactor *application.AddTelegramIdentity,
	logger *slog.Logger,
) *AddTelegramIdentityHandler {
	handlerLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "add_telegram_identity_handler"),
	)

	return &AddTelegramIdentityHandler{
		interactor: interactor,
		logger:     handlerLogger,
	}
}

func (handler *AddTelegramIdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req AddTelegramIdentityRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		handler.logger.DebugContext(r.Context(), "invalid request format", slog.Any("err", err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	requestDTO := application.AddTelegramIdentityRequest{
		UserID:      req.UserID,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Bio:         req.Bio,
		PhoneNumber: req.PhoneNumber,
	}
	resp, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		switch {
		case errors.Is(err, rbac.ErrInsufficientPrivileges):
			handler.logger.DebugContext(r.Context(), "Auth error", slog.Any("err", err))
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		case errors.Is(err, domain.ErrIdentityAlreadyExists):
			handler.logger.DebugContext(r.Context(), "This identity is already added", slog.Any("err", err))
			http.Error(w, "You have already added this identity", http.StatusConflict)
			return
		case errors.Is(err, domain.ErrUnexistentTelegramUserReferenced):
			handler.logger.DebugContext(
				r.Context(),
				"This identity references user that hasn't been added yet",
				slog.Any("err", err),
			)
			http.Error(w, "This identity references user that hasn't been added yet", http.StatusConflict)
			return
		default:
			handler.logger.DebugContext(r.Context(), "Database error", slog.Any("err", err))
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}

	response := AddTelegramIdentityResponse{RecordID: resp.ID.String()}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
