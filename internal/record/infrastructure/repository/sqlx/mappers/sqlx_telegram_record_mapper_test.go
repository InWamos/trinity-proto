package mappers_test

import (
	"testing"
	"time"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
	"github.com/google/uuid"
)

func TestSqlxTelegramRecordMapper_ToDomain(t *testing.T) {
	mapper := mappers.NewSqlxTelegramRecordMapper()

	// Arrange
	id := uuid.New()
	fromUserID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	postedAt := time.Now().Add(-1 * time.Hour)

	model := models.SQLXTelegramRecordModel{
		ID:                 id,
		MessageTelegramID:  uint64(123456789),
		FromTelegramUserID: fromUserID,
		InTelegramChatID:   int64(-1001234567890),
		MessageText:        "Test message content",
		PostedAt:           postedAt,
		AddedAt:            addedAt,
		AddedByUser:        addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert
	if domainEntity.ID != id {
		t.Errorf("expected ID %v, got %v", id, domainEntity.ID)
	}
	if domainEntity.MessageTelegramID != 123456789 {
		t.Errorf("expected MessageTelegramID 123456789, got %d", domainEntity.MessageTelegramID)
	}
	if domainEntity.FromTelegramUserID != fromUserID {
		t.Errorf("expected FromTelegramUserID %v, got %v", fromUserID, domainEntity.FromTelegramUserID)
	}
	if domainEntity.InTelegramChatID != -1001234567890 {
		t.Errorf("expected InTelegramChatID -1001234567890, got %d", domainEntity.InTelegramChatID)
	}
	if domainEntity.MessageText != "Test message content" {
		t.Errorf("expected MessageText 'Test message content', got %s", domainEntity.MessageText)
	}
	if !domainEntity.PostedAt.Equal(postedAt) {
		t.Errorf("expected PostedAt %v, got %v", postedAt, domainEntity.PostedAt)
	}
	if !domainEntity.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, domainEntity.AddedAt)
	}
	if domainEntity.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, domainEntity.AddedByUser)
	}
}

func TestSqlxTelegramRecordMapper_ToModel(t *testing.T) {
	mapper := mappers.NewSqlxTelegramRecordMapper()

	// Arrange
	id := uuid.New()
	fromUserID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	postedAt := time.Now().Add(-2 * time.Hour)

	domainEntity := domain.TelegramRecord{
		ID:                 id,
		MessageTelegramID:  uint64(987654321),
		FromTelegramUserID: fromUserID,
		InTelegramChatID:   int64(-1009876543210),
		MessageText:        "Another test message",
		PostedAt:           postedAt,
		AddedAt:            addedAt,
		AddedByUser:        addedByUser,
	}

	// Act
	model := mapper.ToModel(domainEntity)

	// Assert
	if model.ID != id {
		t.Errorf("expected ID %v, got %v", id, model.ID)
	}
	if model.MessageTelegramID != 987654321 {
		t.Errorf("expected MessageTelegramID 987654321, got %d", model.MessageTelegramID)
	}
	if model.FromTelegramUserID != fromUserID {
		t.Errorf("expected FromTelegramUserID %v, got %v", fromUserID, model.FromTelegramUserID)
	}
	if model.InTelegramChatID != -1009876543210 {
		t.Errorf("expected InTelegramChatID -1009876543210, got %d", model.InTelegramChatID)
	}
	if model.MessageText != "Another test message" {
		t.Errorf("expected MessageText 'Another test message', got %s", model.MessageText)
	}
	if !model.PostedAt.Equal(postedAt) {
		t.Errorf("expected PostedAt %v, got %v", postedAt, model.PostedAt)
	}
	if !model.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, model.AddedAt)
	}
	if model.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, model.AddedByUser)
	}
}

func TestSqlxTelegramRecordMapper_RoundTrip(t *testing.T) {
	mapper := mappers.NewSqlxTelegramRecordMapper()

	// Arrange
	id := uuid.New()
	fromUserID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	postedAt := time.Now().Add(-3 * time.Hour)

	originalEntity := domain.TelegramRecord{
		ID:                 id,
		MessageTelegramID:  uint64(555666777),
		FromTelegramUserID: fromUserID,
		InTelegramChatID:   int64(-1005556667778),
		MessageText:        "Round-trip test message",
		PostedAt:           postedAt,
		AddedAt:            addedAt,
		AddedByUser:        addedByUser,
	}

	// Act - Convert to model and back to domain
	model := mapper.ToModel(originalEntity)
	resultEntity := mapper.ToDomain(model)

	// Assert - Verify round-trip preserves all data
	if resultEntity.ID != originalEntity.ID {
		t.Errorf("round-trip failed: expected ID %v, got %v", originalEntity.ID, resultEntity.ID)
	}
	if resultEntity.MessageTelegramID != originalEntity.MessageTelegramID {
		t.Errorf(
			"round-trip failed: expected MessageTelegramID %d, got %d",
			originalEntity.MessageTelegramID,
			resultEntity.MessageTelegramID,
		)
	}
	if resultEntity.FromTelegramUserID != originalEntity.FromTelegramUserID {
		t.Errorf(
			"round-trip failed: expected FromTelegramUserID %v, got %v",
			originalEntity.FromTelegramUserID,
			resultEntity.FromTelegramUserID,
		)
	}
	if resultEntity.InTelegramChatID != originalEntity.InTelegramChatID {
		t.Errorf(
			"round-trip failed: expected InTelegramChatID %d, got %d",
			originalEntity.InTelegramChatID,
			resultEntity.InTelegramChatID,
		)
	}
	if resultEntity.MessageText != originalEntity.MessageText {
		t.Errorf(
			"round-trip failed: expected MessageText '%s', got '%s'",
			originalEntity.MessageText,
			resultEntity.MessageText,
		)
	}
	if !resultEntity.PostedAt.Equal(originalEntity.PostedAt) {
		t.Errorf("round-trip failed: expected PostedAt %v, got %v", originalEntity.PostedAt, resultEntity.PostedAt)
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

func TestSqlxTelegramRecordMapper_LongMessageText(t *testing.T) {
	mapper := mappers.NewSqlxTelegramRecordMapper()

	// Arrange - Test with a long message (close to 4096 char limit)
	id := uuid.New()
	fromUserID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	postedAt := time.Now().Add(-1 * time.Hour)

	// Create a message with 4000 characters
	longMessage := ""
	for range 400 {
		longMessage += "0123456789" //nolint: perfsprint // Negligible
	}

	model := models.SQLXTelegramRecordModel{
		ID:                 id,
		MessageTelegramID:  uint64(111222333),
		FromTelegramUserID: fromUserID,
		InTelegramChatID:   int64(-1001112223334),
		MessageText:        longMessage,
		PostedAt:           postedAt,
		AddedAt:            addedAt,
		AddedByUser:        addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert
	if len(domainEntity.MessageText) != 4000 {
		t.Errorf("expected MessageText length 4000, got %d", len(domainEntity.MessageText))
	}
	if domainEntity.MessageText != longMessage {
		t.Errorf("long message text was not preserved correctly")
	}
}

func TestSqlxTelegramRecordMapper_NegativeChatID(t *testing.T) {
	mapper := mappers.NewSqlxTelegramRecordMapper()

	// Arrange - Test with negative chat ID (typical for groups/channels)
	id := uuid.New()
	fromUserID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()
	postedAt := time.Now()

	negativeChatID := int64(-1001234567890)

	model := models.SQLXTelegramRecordModel{
		ID:                 id,
		MessageTelegramID:  uint64(12345),
		FromTelegramUserID: fromUserID,
		InTelegramChatID:   negativeChatID,
		MessageText:        "Message in a group",
		PostedAt:           postedAt,
		AddedAt:            addedAt,
		AddedByUser:        addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert
	if domainEntity.InTelegramChatID != negativeChatID {
		t.Errorf("expected negative InTelegramChatID %d, got %d", negativeChatID, domainEntity.InTelegramChatID)
	}
}
