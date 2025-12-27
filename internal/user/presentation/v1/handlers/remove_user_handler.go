//nolint:dupl // Intended to be similar to other handlers
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/google/uuid"
)

type RemoveUserHandler struct {
	interactor *application.RemoveUser
	logger     *slog.Logger
}

func NewRemoveUserHandler(
	interactor *application.RemoveUser,
	logger *slog.Logger,
) *RemoveUserHandler {
	ruhLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "remove_user"),
	)
	return &RemoveUserHandler{
		interactor: interactor,
		logger:     ruhLogger,
	}
}

// ServeHTTP handles an HTTP DELETE request to remove a user.
//
//	@Summary		Delete a user
//	@Description	Permanently remove a user from the system
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string			true	"User ID (UUID)"	format(uuid)
//	@Success		200	{object}	SuccessResponse	"User deleted successfully"
//	@Failure		400	{object}	ErrorResponse	"Invalid user ID format"
//	@Failure		404	{object}	ErrorResponse	"User not found"
//	@Failure		500	{object}	ErrorResponse	"Server error"
//	@Router			/v1/users/{id} [delete]
func (handler *RemoveUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from path parameter
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		handler.logger.DebugContext(r.Context(), "missing user ID in path")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User ID is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handler.logger.DebugContext(r.Context(), "invalid user ID format", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID format"})
		return
	}

	requestDTO := application.RemoveUserRequest{
		ID: userID,
	}

	err = handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.ErrorContext(r.Context(), "failed to remove user", slog.Any("err", err))
		switch {
		case errors.Is(err, application.ErrInsufficientPrivileges):
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Insufficient privileges"})
			return
		case errors.Is(err, application.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User removed successfully",
	})
}
