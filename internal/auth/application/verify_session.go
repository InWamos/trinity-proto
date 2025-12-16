package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/user/client"
)

type VerifySessionRequest struct {
	Session_id string
}

type VerifySession struct {
	userClient client.UserClient
	logger     *slog.Logger
}
