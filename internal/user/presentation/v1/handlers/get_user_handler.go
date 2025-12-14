package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/google/uuid"
)

type GetUserHandler struct {
	interactor *application.GetUserByID
	logger     *slog.Logger
}

func NewGetUserHandler(
	interactor *application.GetUserByID,
	logger *slog.Logger,
) *GetUserHandler {
	guhLogger := logger.With(
		slog.String("component", "handler"),
		slog.String("name", "get_user"),
	)
	return &GetUserHandler{
		interactor: interactor,
		logger:     guhLogger,
	}
}

// ServeHTTP handles an HTTP request to the GET /api/v1/users/{id} endpoint.
func (handler *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	requestDTO := application.GetUserByIDRequest{
		ID: userID,
	}

	response, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.ErrorContext(r.Context(), "failed to get user", slog.Any("err", err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           response.User.ID,
		"username":     response.User.Username,
		"display_name": response.User.DisplayName,
		"role":         response.User.Role,
		"created_at":   response.User.CreatedAt,
	})
}
