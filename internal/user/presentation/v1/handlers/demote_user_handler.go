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

type DemoteUserHandler struct {
	interactor *application.DemoteUser
	logger     *slog.Logger
}

func NewDemoteUserHandler(
	interactor *application.DemoteUser,
	logger *slog.Logger,
) *DemoteUserHandler {
	duhLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "demote_user"),
	)
	return &DemoteUserHandler{
		interactor: interactor,
		logger:     duhLogger,
	}
}

// ServeHTTP handles an HTTP PATCH request to demote a user from admin to regular user.
// @Summary Demote user to regular user
// @Description Change a user's role from admin to user
// @Tags users
// @Produce json
// @Param id path string true "User ID (UUID)" format(uuid)
// @Success 200 {object} SuccessResponse "User demoted successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID format"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users/{id}/demote [patch]
func (handler *DemoteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	requestDTO := application.DemoteUserRequest{
		ID: userID,
	}

	err = handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.ErrorContext(r.Context(), "failed to demote user", slog.Any("err", err))

		// Check if error is due to user not found
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, application.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to demote user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User demoted to regular user successfully",
	})
}
