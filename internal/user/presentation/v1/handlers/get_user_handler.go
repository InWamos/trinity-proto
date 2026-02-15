package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/google/uuid"
)

// GetUserResponse represents the response from the GetUser endpoint
//
//	@Description	User information response
type GetUserResponse struct {
	ID          string    `json:"id"           example:"019b1a49-dbf6-74d6-97bf-2d7e57d30c75"`
	Username    string    `json:"username"     example:"johndoe"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	UserRole    string    `json:"user_role"    example:"user"                                 enums:"user,admin"`
	CreatedAt   time.Time `json:"created_at"   example:"2025-12-14T00:36:46.545Z"`
}

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

// ServeHTTP handles an HTTP GET request to retrieve a user by ID.
//
//	@Summary		Get user by ID
//	@Description	Retrieve a user's information by their ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string			true	"User ID (UUID)"	format(uuid)
//	@Success		200	{object}	GetUserResponse	"User found"
//	@Failure		400	{object}	ErrorResponse	"Invalid user ID format"
//	@Failure		404	{object}	ErrorResponse	"User not found"
//	@Router			/v1/users/{id} [get]
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
		handler.logger.ErrorContext(r.Context(), "failed to get user by ID", slog.Any("err", err))
		switch {
		case errors.Is(err, rbac.ErrInsufficientPrivileges):
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
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":           response.User.ID,
		"username":     response.User.Username,
		"display_name": response.User.DisplayName,
		"user_role":    response.User.Role,
		"created_at":   response.User.CreatedAt,
	})
}
