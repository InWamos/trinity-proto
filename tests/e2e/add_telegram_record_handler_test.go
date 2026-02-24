package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddTelegramRecord_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// First, add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(11111111111),
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	if userResp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create telegram user: status %d, response: %s", userResp.StatusCode, string(userRespBody))
	}

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Now add a telegram record
	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(987654321),
		"from_user_telegram_id": userID,
		"in_telegram_chat_id":   int64(123456789),
		"message_text":          "Hello, this is a test message!",
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
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

func TestAddTelegramRecord_NonexistentUser(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Try to add a record for a user that doesn't exist
	nonexistentUserID := uuid.New().String()

	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(123456789),
		"from_user_telegram_id": nonexistentUserID,
		"in_telegram_chat_id":   int64(987654321),
		"message_text":          "This message references a nonexistent user",
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert conflict (user doesn't exist)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusConflict, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramRecord_Duplicate(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(22222222222),
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Add a record
	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(111222333),
		"from_user_telegram_id": userID,
		"in_telegram_chat_id":   int64(444555666),
		"message_text":          "Duplicate test message",
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("initial record creation failed with status %d", resp.StatusCode)
	}

	// Try to add the same record again (same message_telegram_id)
	resp = MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert conflict
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusConflict, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramRecord_InvalidMessageTelegramID_Zero(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(33333333333),
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Try to add record with invalid message_telegram_id (0)
	recordReq := map[string]interface{}{
		"message_telegram_id":   0, // Invalid: must be > 0
		"from_user_telegram_id": userID,
		"in_telegram_chat_id":   int64(123456789),
		"message_text":          "Test message",
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusUnprocessableEntity,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestAddTelegramRecord_MessageTextTooLong(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(44444444444),
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Create a message that's too long (max is 4096 characters)
	longMessage := ""
	for i := 0; i < 5000; i++ {
		longMessage += "a"
	}

	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(555666777),
		"from_user_telegram_id": userID,
		"in_telegram_chat_id":   int64(888999000),
		"message_text":          longMessage,
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusUnprocessableEntity,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestAddTelegramRecord_EmptyMessageText(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user with a unique ID based on timestamp to avoid conflicts
	telegramUserReq := map[string]interface{}{
		"telegram_id": uint64(55555000000 + time.Now().UnixNano()%1000000), // Unique ID for each run
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userRespBody, _ := io.ReadAll(userResp.Body)
	userResp.Body.Close()

	var userResponse map[string]string
	json.Unmarshal(userRespBody, &userResponse)
	userID := userResponse["record_id"]

	// Try to add record with empty message text (required)
	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(666000000 + time.Now().UnixNano()%1000000), // Unique ID for each run
		"from_user_telegram_id": userID,
		"in_telegram_chat_id":   int64(999000111),
		"message_text":          "", // Empty, but required
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusUnprocessableEntity,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestAddTelegramRecord_Unauthorized(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Don't provide a token
	recordReq := map[string]interface{}{
		"message_telegram_id":   uint64(123456789),
		"from_user_telegram_id": uuid.New().String(),
		"in_telegram_chat_id":   int64(987654321),
		"message_text":          "Unauthorized test message",
		"posted_at":             time.Now().Format(time.RFC3339),
	}

	bodyJSON, _ := json.Marshal(recordReq)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL),
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

func TestAddTelegramRecord_InvalidRequestFormat(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Send invalid JSON (missing required fields)
	recordReq := map[string]interface{}{
		"invalid_field": "value",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL), token, recordReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert bad request or unprocessable entity
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d or %d, got %d. Response: %s",
			http.StatusBadRequest, http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}
