package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestAddTelegramUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Login as admin user (users need at least user role to add telegram users)
	token := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request to add a telegram user
	reqBody := map[string]interface{}{
		"telegram_id": uint64(98765432101), // Use a large unique telegram ID
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	defer resp.Body.Close()

	// Read response body
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

func TestAddTelegramUser_Duplicate(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user first time
	reqBody := map[string]interface{}{
		"telegram_id": uint64(12345678901),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("initial user creation failed with status %d", resp.StatusCode)
	}

	// Try to add the same telegram user again
	resp = MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert conflict
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusConflict, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramUser_InvalidTelegramID_Zero(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Try to add telegram user with invalid ID (0)
	reqBody := map[string]interface{}{
		"telegram_id": 0,
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramUser_InvalidTelegramID_TooLarge(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Try to add telegram user with ID exceeding maximum (300000000000)
	reqBody := map[string]interface{}{
		"telegram_id": uint64(400000000000),
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert unprocessable entity
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}

func TestAddTelegramUser_Unauthorized(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Don't provide a token - make unauthenticated request
	reqBody := map[string]interface{}{
		"telegram_id": uint64(11111111111),
	}

	bodyJSON, _ := json.Marshal(reqBody)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
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

func TestAddTelegramUser_InvalidRequestFormat(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Send invalid JSON (missing required field)
	reqBody := map[string]interface{}{
		"invalid_field": "value",
	}

	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL), token, reqBody)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert bad request or unprocessable entity
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d or %d, got %d. Response: %s",
			http.StatusBadRequest, http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
	}
}
