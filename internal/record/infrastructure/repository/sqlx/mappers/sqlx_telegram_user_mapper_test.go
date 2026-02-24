package mappers_test

import (
	"testing"
	"time"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
	"github.com/google/uuid"
)

func TestSqlxTelegramUserMapper_ToDomain(t *testing.T) {
	mapper := mappers.NewSqlxTelegramUserMapper()

	// Arrange
	id := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	telegramID := uint64(123456789)

	model := models.TelegramUserModel{
		ID:          id,
		TelegramID:  telegramID,
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert
	if domainEntity.ID != id {
		t.Errorf("expected ID %v, got %v", id, domainEntity.ID)
	}
	if domainEntity.TelegramID != telegramID {
		t.Errorf("expected TelegramID %d, got %d", telegramID, domainEntity.TelegramID)
	}
	if !domainEntity.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, domainEntity.AddedAt)
	}
	if domainEntity.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, domainEntity.AddedByUser)
	}
}

func TestSqlxTelegramUserMapper_ToModel(t *testing.T) {
	mapper := mappers.NewSqlxTelegramUserMapper()

	// Arrange
	id := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	telegramID := uint64(123456789)

	domainEntity := domain.TelegramUser{
		ID:          id,
		TelegramID:  telegramID,
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act
	model := mapper.ToModel(domainEntity)

	// Assert
	if model.ID != id {
		t.Errorf("expected ID %v, got %v", id, model.ID)
	}
	if model.TelegramID != telegramID {
		t.Errorf("expected TelegramID %d, got %d", telegramID, model.TelegramID)
	}
	if !model.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, model.AddedAt)
	}
	if model.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, model.AddedByUser)
	}
}

func TestSqlxTelegramUserMapper_RoundTrip(t *testing.T) {
	mapper := mappers.NewSqlxTelegramUserMapper()

	// Arrange
	id := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	telegramID := uint64(987654321)

	originalEntity := domain.TelegramUser{
		ID:          id,
		TelegramID:  telegramID,
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act - Convert to model and back to domain
	model := mapper.ToModel(originalEntity)
	resultEntity := mapper.ToDomain(model)

	// Assert - Verify round-trip preserves all data
	if resultEntity.ID != originalEntity.ID {
		t.Errorf("round-trip failed: expected ID %v, got %v", originalEntity.ID, resultEntity.ID)
	}
	if resultEntity.TelegramID != originalEntity.TelegramID {
		t.Errorf(
			"round-trip failed: expected TelegramID %d, got %d",
			originalEntity.TelegramID,
			resultEntity.TelegramID,
		)
	}
	if !resultEntity.AddedAt.Equal(originalEntity.AddedAt) {
		t.Errorf("round-trip failed: expected AddedAt %v, got %v", originalEntity.AddedAt, resultEntity.AddedAt)
	}
	if resultEntity.AddedByUser != originalEntity.AddedByUser {
		t.Errorf(
			"round-trip failed: expected AddedByUser %v, got %v",
			originalEntity.AddedByUser,
			resultEntity.AddedByUser,
		)
	}
}
