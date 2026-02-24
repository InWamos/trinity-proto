package application

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	"github.com/google/uuid"
)

type RemoveAllSessionsByUserIDRequest struct {
	UserID uuid.UUID
}

type RemoveAllSessionsByUserIDResponse struct {
	Message string
}

type RemoveAllSessionsByUserID struct {
	sessionRepository infrastructure.SessionRepository
	logger            *slog.Logger
}

func NewRemoveAllSessionsByUserID(
	sessionRepository infrastructure.SessionRepository,
	logger *slog.Logger,
) *RemoveAllSessionsByUserID {
	rasLogger := logger.With(slog.String("module", "auth"), slog.String("name", "remove_all_sessions_by_user_id"))
	return &RemoveAllSessionsByUserID{
		sessionRepository: sessionRepository,
		logger:            rasLogger,
	}
}

func (ras *RemoveAllSessionsByUserID) Execute(
	ctx context.Context,
	input RemoveAllSessionsByUserIDRequest,
) (RemoveAllSessionsByUserIDResponse, error) {
	err := ras.sessionRepository.RevokeAllSessionsByUserID(ctx, input.UserID)
	if err != nil {
		ras.logger.ErrorContext(
			ctx,
			"Failed to revoke all sessions for user",
			slog.String("user_id", input.UserID.String()),
			slog.Any("err", err),
		)
		return RemoveAllSessionsByUserIDResponse{}, ErrUnexpected
	}

	ras.logger.InfoContext(
		ctx,
		"All sessions revoked successfully for user",
		slog.String("user_id", input.UserID.String()),
	)

	return RemoveAllSessionsByUserIDResponse{Message: "All sessions revoked successfully"}, nil
}
