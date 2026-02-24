package application

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	"github.com/InWamos/trinity-proto/internal/shared/authorization/rbac"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/middleware"
)

type GetAllSessionsByUserIDResponse struct {
	Sessions []domain.Session
}

type GetAllSessionsByUserID struct {
	sessionRepository infrastructure.SessionRepository
	logger            *slog.Logger
}

func NewGetAllSessionsByUserID(
	sessionRepository infrastructure.SessionRepository,
	logger *slog.Logger,
) *GetAllSessionsByUserID {
	gasLogger := logger.With(slog.String("module", "auth"), slog.String("name", "get_all_sessions_by_user_id"))
	return &GetAllSessionsByUserID{
		sessionRepository: sessionRepository,
		logger:            gasLogger,
	}
}

func (gas *GetAllSessionsByUserID) Execute(
	ctx context.Context,
) (GetAllSessionsByUserIDResponse, error) {
	idp, ok := ctx.Value(middleware.IdentityProviderKey).(*client.UserIdentity)
	if !ok || idp == nil {
		return GetAllSessionsByUserIDResponse{}, rbac.ErrInsufficientPrivileges
	}
	userID := idp.UserID
	// Retrieve all sessions for the user
	sessions, err := gas.sessionRepository.GetAllSessionsByUserID(ctx, userID)
	if err != nil {
		gas.logger.ErrorContext(
			ctx,
			"failed to retrieve sessions for user",
			slog.String("user_id", userID.String()),
			slog.Any("err", err),
		)
		return GetAllSessionsByUserIDResponse{}, err
	}

	gas.logger.DebugContext(ctx, "retrieved all sessions for user",
		slog.String("user_id", userID.String()),
		slog.Int("session_count", len(sessions)),
	)

	return GetAllSessionsByUserIDResponse{Sessions: sessions}, nil
}
