package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/auth/domain"
)

// LogoutResponse represents the response from the Logout endpoint
//
//	@Description	Logout response with the message
type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

type LogoutHandler struct {
	interactor *application.RemoveSession
	logger     *slog.Logger
}

// NewLogoutHandler builds a new LogoutHandler.
func NewLogoutHandler(
	interactor *application.RemoveSession,
	logger *slog.Logger,
) *LogoutHandler {
	lhLogger := logger.With(slog.String("component", "handler"), slog.String("name", "logout"))
	return &LogoutHandler{
		interactor: interactor,
		logger:     lhLogger,
	}
}

// ServeHTTP handles an HTTP request to logout a user.
//
//	@Summary		User logout
//	@Description	Revoke the current session token and clear the session
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LogoutResponse	"Logout successful"
//	@Failure		401	string		"Invalid request (no active session)"
//	@Failure		500	string		"Server error"
//	@Router			/v1/auth/logout [post]
func (handler *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		handler.logger.DebugContext(r.Context(), "failed to get session token from cookie", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "No active session"})
		return
	}

	// Revoke the session
	requestDTO := application.RemoveSessionRequest{
		Token: cookie.Value,
	}
	err = handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSessionNotFound):
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Negative MaxAge deletes the cookie
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
