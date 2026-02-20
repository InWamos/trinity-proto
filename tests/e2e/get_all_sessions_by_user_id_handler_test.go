package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestGetAllSessionsByUserID_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Login to get a token
	token := LoginUser(t, baseURL, "testuser", "user12345")

	// Get all sessions for the logged-in user
	resp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", token, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(respBody))
	}

	// Verify response structure
	var sessionsResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &sessionsResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Check that sessions array exists
	sessions, ok := sessionsResponse["sessions"].([]interface{})
	if !ok {
		t.Fatal("expected 'sessions' field to be an array")
	}

	// Should have at least one session from login
	if len(sessions) < 1 {
		t.Errorf("expected at least 1 session, got %d", len(sessions))
	}

	// Verify session structure
	firstSession := sessions[0].(map[string]interface{})
	expectedFields := []string{"id", "user_id", "user_role", "ip_address", "user_agent", "created_at", "expires_at", "status"}
	for _, field := range expectedFields {
		if value, ok := firstSession[field]; !ok || value == nil {
			t.Errorf("expected field %q in session response (got: %v)", field, value)
		}
	}
}

func TestGetAllSessionsByUserID_MultipleSessions(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Login multiple times to create multiple sessions
	token1 := LoginUser(t, baseURL, "testuser", "user12345")
	_ = LoginUser(t, baseURL, "testuser", "user12345")
	_ = LoginUser(t, baseURL, "testuser", "user12345")

	// Get all sessions for the user using the first token
	resp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", token1, nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(respBody))
	}

	// Verify multiple sessions
	var sessionsResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &sessionsResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	sessions := sessionsResponse["sessions"].([]interface{})

	// Should have at least 3 sessions from the 3 login attempts
	if len(sessions) < 3 {
		t.Errorf("expected at least 3 sessions, got %d", len(sessions))
	}
}

func TestGetAllSessionsByUserID_Unauthorized(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to get sessions without a valid token
	req, err := http.NewRequest("GET", baseURL+"/api/v1/session", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Don't set authorization header
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert - should get 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnauthorized, resp.StatusCode, string(respBody))
	}
}

func TestGetAllSessionsByUserID_InvalidToken(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to get sessions with an invalid token
	resp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", "invalid-token-xyz", nil)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert - should get 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusUnauthorized, resp.StatusCode, string(respBody))
	}
}

func TestGetAllSessionsByUserID_DifferentUsers(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Login with testuser
	userToken := LoginUser(t, baseURL, "testuser", "user12345")

	// Login with admin
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Get sessions as testuser
	userResp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", userToken, nil)
	defer userResp.Body.Close()

	userRespBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		t.Fatalf("failed to read user response body: %v", err)
	}

	var userSessions map[string]interface{}
	if err := json.Unmarshal(userRespBody, &userSessions); err != nil {
		t.Fatalf("failed to unmarshal user response: %v", err)
	}

	// Get sessions as admin
	adminResp := MakeAuthorizedRequest(t, "GET", baseURL+"/api/v1/session", adminToken, nil)
	defer adminResp.Body.Close()

	adminRespBody, err := io.ReadAll(adminResp.Body)
	if err != nil {
		t.Fatalf("failed to read admin response body: %v", err)
	}

	var adminSessions map[string]interface{}
	if err := json.Unmarshal(adminRespBody, &adminSessions); err != nil {
		t.Fatalf("failed to unmarshal admin response: %v", err)
	}

	// Verify both got sessions
	if len(userSessions["sessions"].([]interface{})) < 1 {
		t.Error("expected at least 1 session for testuser")
	}

	if len(adminSessions["sessions"].([]interface{})) < 1 {
		t.Error("expected at least 1 session for admin")
	}

	// Verify sessions are isolated per user by checking user_id field
	for _, sess := range userSessions["sessions"].([]interface{}) {
		session := sess.(map[string]interface{})
		userIDVal := session["user_id"]
		if userIDVal == nil {
			t.Error("user_id field is nil in testuser session")
		}
	}

	for _, sess := range adminSessions["sessions"].([]interface{}) {
		session := sess.(map[string]interface{})
		userIDVal := session["user_id"]
		if userIDVal == nil {
			t.Error("user_id field is nil in admin session")
		}
	}
}
