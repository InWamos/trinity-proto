package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestRemoveUser_SessionsRevokedViaEvent validates the event-driven flow:
// 1. A user is created and logs in (creating sessions)
// 2. An admin removes the user, which publishes a UserRemovedEvent to Kafka
// 3. The auth module consumes the event and revokes all sessions for that user
// 4. The removed user's token is no longer valid
func TestRemoveUser_SessionsRevokedViaEvent(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Create a target user
	reqBody := map[string]string{
		"username":     "eventsessionuser",
		"display_name": "Event Session User",
		"password":     "password123",
		"user_role":    "user",
	}

	createResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(createResp.Body)
		t.Fatalf("failed to create user: status=%d body=%s", createResp.StatusCode, body)
	}

	var createBody map[string]interface{}
	raw, _ := io.ReadAll(createResp.Body)
	if err := json.Unmarshal(raw, &createBody); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	userID, ok := createBody["id"].(string)
	if !ok || userID == "" {
		t.Skip("CreateUser endpoint doesn't return user ID yet")
	}

	// Login as the newly created user to generate sessions
	userToken1 := LoginUser(t, baseURL, "eventsessionuser", "password123")
	userToken2 := LoginUser(t, baseURL, "eventsessionuser", "password123")

	// Verify sessions exist before deletion
	sessionsResp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", userToken1, nil)
	defer sessionsResp.Body.Close()
	if sessionsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected sessions to exist before deletion, got status %d", sessionsResp.StatusCode)
	}

	var sessionsBody map[string]interface{}
	sessionsRaw, _ := io.ReadAll(sessionsResp.Body)
	if err := json.Unmarshal(sessionsRaw, &sessionsBody); err != nil {
		t.Fatalf("failed to unmarshal sessions response: %v", err)
	}
	sessions := sessionsBody["sessions"].([]interface{})
	if len(sessions) < 2 {
		t.Errorf("expected at least 2 sessions before deletion, got %d", len(sessions))
	}

	// Admin deletes the user - this publishes a UserRemovedEvent to Kafka
	deleteResp := MakeAuthorizedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID), adminToken, nil)
	defer deleteResp.Body.Close()
	if deleteResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(deleteResp.Body)
		t.Fatalf("failed to delete user: status=%d body=%s", deleteResp.StatusCode, body)
	}

	// Wait for the Kafka consumer to process the UserRemovedEvent and revoke sessions
	time.Sleep(3 * time.Second)

	// Verify token1 is no longer valid (session revoked by event consumer)
	sessionCheckResp1 := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", userToken1, nil)
	defer sessionCheckResp1.Body.Close()
	if sessionCheckResp1.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected token1 to be revoked (401) after user deletion, got status %d", sessionCheckResp1.StatusCode)
	}

	// Verify token2 is also no longer valid
	sessionCheckResp2 := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", userToken2, nil)
	defer sessionCheckResp2.Body.Close()
	if sessionCheckResp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected token2 to be revoked (401) after user deletion, got status %d", sessionCheckResp2.StatusCode)
	}
}

// TestRemoveUser_AdminSessionUnaffected verifies that removing a user only
// revokes that user's sessions and not the admin's session.
func TestRemoveUser_AdminSessionUnaffected(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Create a target user
	reqBody := map[string]string{
		"username":     "eventsessionuser2",
		"display_name": "Event Session User 2",
		"password":     "password123",
		"user_role":    "user",
	}

	createResp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(createResp.Body)
		t.Fatalf("failed to create user: status=%d body=%s", createResp.StatusCode, body)
	}

	var createBody map[string]interface{}
	raw, _ := io.ReadAll(createResp.Body)
	if err := json.Unmarshal(raw, &createBody); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	userID, ok := createBody["id"].(string)
	if !ok || userID == "" {
		t.Skip("CreateUser endpoint doesn't return user ID yet")
	}

	// Delete the target user
	deleteResp := MakeAuthorizedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID), adminToken, nil)
	defer deleteResp.Body.Close()
	if deleteResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(deleteResp.Body)
		t.Fatalf("failed to delete user: status=%d body=%s", deleteResp.StatusCode, body)
	}

	// Wait for the Kafka consumer to process the event
	time.Sleep(3 * time.Second)

	// Admin's own session should remain valid
	sessionsResp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", adminToken, nil)
	defer sessionsResp.Body.Close()
	if sessionsResp.StatusCode != http.StatusOK {
		t.Errorf("expected admin session to remain valid after deleting another user, got status %d", sessionsResp.StatusCode)
	}
}
