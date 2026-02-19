package mappers_test

import (
	"testing"
	"time"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/models"
	"github.com/google/uuid"
)

func TestSqlxTelegramIdentityMapper_ToDomain(t *testing.T) {
	mapper := mappers.NewSqlxTelegramIdentityMapper()

	// Arrange
	id := uuid.New()
	userID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()

	model := models.TelegramIdentityModel{
		ID:          id,
		UserID:      userID,
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "johndoe",
		PhoneNumber: "+1234567890",
		Bio:         "Test bio",
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert
	if domainEntity.ID != id {
		t.Errorf("expected ID %v, got %v", id, domainEntity.ID)
	}
	if domainEntity.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, domainEntity.UserID)
	}
	if domainEntity.FirstName != "John" {
		t.Errorf("expected FirstName 'John', got %s", domainEntity.FirstName)
	}
	if domainEntity.LastName != "Doe" {
		t.Errorf("expected LastName 'Doe', got %s", domainEntity.LastName)
	}
	if domainEntity.Username != "johndoe" {
		t.Errorf("expected Username 'johndoe', got %s", domainEntity.Username)
	}
	if domainEntity.PhoneNumber != "+1234567890" {
		t.Errorf("expected PhoneNumber '+1234567890', got %s", domainEntity.PhoneNumber)
	}
	if domainEntity.Bio != "Test bio" {
		t.Errorf("expected Bio 'Test bio', got %s", domainEntity.Bio)
	}
	if !domainEntity.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, domainEntity.AddedAt)
	}
	if domainEntity.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, domainEntity.AddedByUser)
	}
}

func TestSqlxTelegramIdentityMapper_ToModel(t *testing.T) {
	mapper := mappers.NewSqlxTelegramIdentityMapper()

	// Arrange
	id := uuid.New()
	userID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()

	domainEntity := domain.TelegramIdentity{
		ID:          id,
		UserID:      userID,
		FirstName:   "Jane",
		LastName:    "Smith",
		Username:    "janesmith",
		PhoneNumber: "+0987654321",
		Bio:         "Another test bio",
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act
	model := mapper.ToModel(domainEntity)

	// Assert
	if model.ID != id {
		t.Errorf("expected ID %v, got %v", id, model.ID)
	}
	if model.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, model.UserID)
	}
	if model.FirstName != "Jane" {
		t.Errorf("expected FirstName 'Jane', got %s", model.FirstName)
	}
	if model.LastName != "Smith" {
		t.Errorf("expected LastName 'Smith', got %s", model.LastName)
	}
	if model.Username != "janesmith" {
		t.Errorf("expected Username 'janesmith', got %s", model.Username)
	}
	if model.PhoneNumber != "+0987654321" {
		t.Errorf("expected PhoneNumber '+0987654321', got %s", model.PhoneNumber)
	}
	if model.Bio != "Another test bio" {
		t.Errorf("expected Bio 'Another test bio', got %s", model.Bio)
	}
	if !model.AddedAt.Equal(addedAt) {
		t.Errorf("expected AddedAt %v, got %v", addedAt, model.AddedAt)
	}
	if model.AddedByUser != addedByUser {
		t.Errorf("expected AddedByUser %v, got %v", addedByUser, model.AddedByUser)
	}
}

func TestSqlxTelegramIdentityMapper_RoundTrip(t *testing.T) {
	mapper := mappers.NewSqlxTelegramIdentityMapper()

	// Arrange
	id := uuid.New()
	userID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()

	originalEntity := domain.TelegramIdentity{
		ID:          id,
		UserID:      userID,
		FirstName:   "Alice",
		LastName:    "Johnson",
		Username:    "alicej",
		PhoneNumber: "+1122334455",
		Bio:         "Round-trip test",
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
	if resultEntity.UserID != originalEntity.UserID {
		t.Errorf("round-trip failed: expected UserID %v, got %v", originalEntity.UserID, resultEntity.UserID)
	}
	if resultEntity.FirstName != originalEntity.FirstName {
		t.Errorf(
			"round-trip failed: expected FirstName '%s', got '%s'",
			originalEntity.FirstName,
			resultEntity.FirstName,
		)
	}
	if resultEntity.LastName != originalEntity.LastName {
		t.Errorf("round-trip failed: expected LastName '%s', got '%s'", originalEntity.LastName, resultEntity.LastName)
	}
	if resultEntity.Username != originalEntity.Username {
		t.Errorf("round-trip failed: expected Username '%s', got '%s'", originalEntity.Username, resultEntity.Username)
	}
	if resultEntity.PhoneNumber != originalEntity.PhoneNumber {
		t.Errorf(
			"round-trip failed: expected PhoneNumber '%s', got '%s'",
			originalEntity.PhoneNumber,
			resultEntity.PhoneNumber,
		)
	}
	if resultEntity.Bio != originalEntity.Bio {
		t.Errorf("round-trip failed: expected Bio '%s', got '%s'", originalEntity.Bio, resultEntity.Bio)
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

func TestSqlxTelegramIdentityMapper_EmptyOptionalFields(t *testing.T) {
	mapper := mappers.NewSqlxTelegramIdentityMapper()

	// Arrange - Test with empty optional fields
	id := uuid.New()
	userID := uuid.New()
	addedByUser := uuid.New()
	addedAt := time.Now()

	model := models.TelegramIdentityModel{
		ID:          id,
		UserID:      userID,
		FirstName:   "Bob",
		LastName:    "", // Optional field - empty
		Username:    "bobby",
		PhoneNumber: "", // Optional field - empty
		Bio:         "", // Optional field - empty
		AddedAt:     addedAt,
		AddedByUser: addedByUser,
	}

	// Act
	domainEntity := mapper.ToDomain(model)

	// Assert - Empty strings should be preserved
	if domainEntity.LastName != "" {
		t.Errorf("expected empty LastName, got %s", domainEntity.LastName)
	}
	if domainEntity.PhoneNumber != "" {
		t.Errorf("expected empty PhoneNumber, got %s", domainEntity.PhoneNumber)
	}
	if domainEntity.Bio != "" {
		t.Errorf("expected empty Bio, got %s", domainEntity.Bio)
	}
	// Required fields should still be present
	if domainEntity.FirstName != "Bob" {
		t.Errorf("expected FirstName 'Bob', got %s", domainEntity.FirstName)
	}
	if domainEntity.Username != "bobby" {
		t.Errorf("expected Username 'bobby', got %s", domainEntity.Username)
	}
}
