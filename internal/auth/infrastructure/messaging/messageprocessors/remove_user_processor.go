package messageprocessors

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/events/pb"
)

type RemoveUserProcessor struct {
	interactor *application.RemoveAllSessionsByUserID
	logger     *slog.Logger
}

func NewRemoveUserProcessor(
	removeAllSessions *application.RemoveAllSessionsByUserID,
	logger *slog.Logger,
) *RemoveUserProcessor {
	return &RemoveUserProcessor{
		interactor: removeAllSessions,
		logger:     logger.With(slog.String("component", "remove_user_processor")),
	}
}

// Process handles UserRemovedEvent by removing all user sessions.
func (p *RemoveUserProcessor) Process(ctx context.Context, event *pb.UserRemovedEvent) error {
	// Parse user ID
	userID, err := uuid.Parse(event.GetUserId())
	if err != nil {
		p.logger.ErrorContext(ctx, "Invalid user_id in event",
			slog.String("user_id", event.GetUserId()),
			slog.Any("err", err))
		return err
	}

	// Call interactor to remove all sessions
	_, err = p.interactor.Execute(ctx,
		application.RemoveAllSessionsByUserIDRequest{UserID: userID})
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to remove user sessions",
			slog.String("user_id", event.GetUserId()),
			slog.Any("err", err))
		return err
	}

	p.logger.DebugContext(ctx, "User sessions removed",
		slog.String("user_id", event.GetUserId()))
	return nil
}
