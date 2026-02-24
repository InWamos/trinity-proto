package messageprocessors_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/auth/infrastructure/messaging/messageprocessors"
	"github.com/InWamos/trinity-proto/internal/events/pb"
)

// mockSessionRepository is a test double for infrastructure.SessionRepository
type mockSessionRepository struct {
	revokeAllCalled bool
	revokeAllUserID uuid.UUID
	revokeAllErr    error
}

func (m *mockSessionRepository) RevokeAllSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	m.revokeAllCalled = true
	m.revokeAllUserID = userID
	return m.revokeAllErr
}

func (m *mockSessionRepository) GetSessionByToken(ctx context.Context, token string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (m *mockSessionRepository) RevokeSessionByToken(ctx context.Context, token string) error {
	return nil
}

func (m *mockSessionRepository) GetAllSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	return nil, nil
}

func (m *mockSessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	return nil
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func newProcessor(repo *mockSessionRepository) *messageprocessors.RemoveUserProcessor {
	interactor := application.NewRemoveAllSessionsByUserID(repo, newTestLogger())
	return messageprocessors.NewRemoveUserProcessor(interactor, newTestLogger())
}

func TestRemoveUserProcessor_Process_Success(t *testing.T) {
	repo := &mockSessionRepository{}
	processor := newProcessor(repo)

	userID := uuid.New()
	event := &pb.UserRemovedEvent{UserId: userID.String()}

	err := processor.Process(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repo.revokeAllCalled {
		t.Error("expected RevokeAllSessionsByUserID to be called")
	}
	if repo.revokeAllUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, repo.revokeAllUserID)
	}
}

func TestRemoveUserProcessor_Process_InvalidUserID(t *testing.T) {
	repo := &mockSessionRepository{}
	processor := newProcessor(repo)

	event := &pb.UserRemovedEvent{UserId: "not-a-valid-uuid"}

	err := processor.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
	if repo.revokeAllCalled {
		t.Error("expected RevokeAllSessionsByUserID NOT to be called on invalid UUID")
	}
}

func TestRemoveUserProcessor_Process_RepositoryError(t *testing.T) {
	repoErr := errors.New("redis connection lost")
	repo := &mockSessionRepository{revokeAllErr: repoErr}
	processor := newProcessor(repo)

	userID := uuid.New()
	event := &pb.UserRemovedEvent{UserId: userID.String()}

	err := processor.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error when repository fails, got nil")
	}
	if !repo.revokeAllCalled {
		t.Error("expected RevokeAllSessionsByUserID to be called even on error")
	}
}
