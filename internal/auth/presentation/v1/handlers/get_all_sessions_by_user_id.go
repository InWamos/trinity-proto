package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/auth/domain"
)

// GetAllSessionsByUserIDResponse represents the response from the GetAllSessionsByUserIDResponse endpoint
//
//	@Description	Response with a list of sessions
type GetAllSessionsByUserIDResponse struct {
	Sessions []map[string]any `json:"sessions"`
}

type GetAllSessionsByUserIDHandler struct {
	interactor *application.GetAllSessionsByUserID
	logger     *slog.Logger
}

// NewGetAllSessionsByUserIDHandler builds a new GetAllSessionsByUserIDHandler.
func NewGetAllSessionsByUserIDHandler(
	interactor *application.GetAllSessionsByUserID,
	logger *slog.Logger,
) *GetAllSessionsByUserIDHandler {
	gasLogger := logger.With(slog.String("component", "handler"), slog.String("name", "get_all_sessions_by_user_id"))
	return &GetAllSessionsByUserIDHandler{
		interactor: interactor,
		logger:     gasLogger,
	}
}

// ServeHTTP handles an HTTP request to get all sessions for the authenticated user.
//
//	@Summary		Get all user sessions
//	@Description	Retrieve all active sessions for the authenticated user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetAllSessionsByUserIDResponse	"Sessions retrieved successfully"
//	@Failure		401	{object}	ErrorResponse					"Unauthorized (no valid session)"
//	@Failure		500	{object}	ErrorResponse					"Server error"
//	@Router			/v1/auth/sessions [get]
func (handler *GetAllSessionsByUserIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := handler.interactor.Execute(r.Context())
	if err != nil {
		handler.logger.ErrorContext(r.Context(), "failed to get sessions", slog.Any("err", err))

		switch {
		case errors.Is(err, domain.ErrSessionNotFound):
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "No sessions found"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "The server was unable to complete your request. Please try again later",
			})
		}
		return
	}

	// Convert sessions to map format
	sessionsList := make([]map[string]any, len(response.Sessions))
	for i, session := range response.Sessions {
		sessionsList[i] = map[string]any{
			"id":         session.ID.String(),
			"user_id":    session.UserID.String(),
			"user_role":  session.UserRole,
			"ip_address": session.IPAddress,
			"user_agent": session.UserAgent,
			"created_at": session.CreatedAt,
			"expires_at": session.ExpiresAt,
			"status":     session.Status,
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetAllSessionsByUserIDResponse{
		Sessions: sessionsList,
	})
}
