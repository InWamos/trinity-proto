package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestAddTelegramIdentity_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// First, add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(55555555555),
	}
	userResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, telegramUserReq)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	if userResp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create telegram user: status %d, response: %s", userResp.StatusCode, string(userRespBody))
	}

	// Extract the user ID from response
	var userResponse map[string]string
	if err := json.Unmarshal(userRespBody, &userResponse); err != nil {
		t.Fatalf("failed to unmarshal user response: %v", err)
	}

	userID := userResponse["record_id"]

	// Now add an identity for that user
	identityReq := map[string]interface{}{
		"telegram_id":           userID,
		"telegram_username":     "testuser123",
		"telegram_first_name":   "John",
		"telegram_last_name":    "Doe",
		"telegram_bio":          "This is a test bio",
		"telegram_phone_number": "+12025550123",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusCreated, resp.StatusCode, string(respBody))
	}

	// Verify response structure
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["record_id"] == "" {
		t.Errorf("expected non-empty record_id in response")
	}
}

func TestAddTelegramIdentity_NonexistentUser(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Try to add identity for a user that doesn't exist
	nonexistentUserID := uuid.New().String()

	identityReq := map[string]interface{}{
		"telegram_id":           nonexistentUserID,
		"telegram_username":     "testuser456",
		"telegram_first_name":   "Jane",
		"telegram_last_name":    "Smith",
		"telegram_bio":          "Another test bio",
		"telegram_phone_number": "+12025550124",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert conflict (user doesn't exist)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusConflict, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramIdentity_Duplicate(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(66666666666),
	}
	userResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, telegramUserReq)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Add identity first time
	identityReq := map[string]interface{}{
		"telegram_id":           userID,
		"telegram_username":     "duplicateuser",
		"telegram_first_name":   "Duplicate",
		"telegram_last_name":    "User",
		"telegram_bio":          "Test duplicate",
		"telegram_phone_number": "+12025550125",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("initial identity creation failed with status %d", resp.StatusCode)
	}

	// Try to add the same identity again
	resp = MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert conflict
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusConflict, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramIdentity_InvalidUsername_TooShort(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(77777777777),
	}
	userResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, telegramUserReq)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Try to add identity with username too short (min is 4)
	identityReq := map[string]interface{}{
		"telegram_id":           userID,
		"telegram_username":     "abc", // Too short
		"telegram_first_name":   "Test",
		"telegram_last_name":    "User",
		"telegram_bio":          "Test",
		"telegram_phone_number": "+12025550126",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramIdentity_InvalidPhoneNumber(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(88888888888),
	}
	userResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, telegramUserReq)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Try to add identity with invalid phone number (not e164 format)
	identityReq := map[string]interface{}{
		"telegram_id":           userID,
		"telegram_username":     "validuser123",
		"telegram_first_name":   "Valid",
		"telegram_last_name":    "User",
		"telegram_bio":          "Test",
		"telegram_phone_number": "123456", // Invalid format
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramIdentity_EmptyFirstName(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(99999999999),
	}
	userResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, telegramUserReq)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Try to add identity with empty first name (required, min=1)
	identityReq := map[string]interface{}{
		"telegram_id":           userID,
		"telegram_username":     "testuser789",
		"telegram_first_name":   "", // Empty, but required
		"telegram_last_name":    "User",
		"telegram_bio":          "Test",
		"telegram_phone_number": "+12025550127",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL), token, identityReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramIdentity_Unauthorized(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Don't provide a token
	identityReq := map[string]interface{}{
		"telegram_id":           uuid.New().String(),
		"telegram_username":     "testuser",
		"telegram_first_name":   "Test",
		"telegram_last_name":    "User",
		"telegram_bio":          "Test",
		"telegram_phone_number": "+12025550128",
	}

	bodyJSON, _ := json.Marshal(identityReq)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/record/telegram/identity", baseURL),
		"application/json",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert unauthorized
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d or %d, got %d. Response: %s. RequestBody: %s",
			http.StatusUnauthorized, http.StatusForbidden, resp.StatusCode, string(respBody), string(bodyJSON))
	}
}
