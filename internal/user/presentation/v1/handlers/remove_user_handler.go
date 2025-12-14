package handlers

import (
	"database/sql"
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

// ServeHTTP handles an HTTP request to the DELETE /api/v1/users/{id} endpoint.
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
		
		// Check if error is due to user not found
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, application.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to remove user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User removed successfully",
	})
}
