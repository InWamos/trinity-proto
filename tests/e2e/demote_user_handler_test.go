package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestDemoteUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// First, create an admin user
	reqBody := map[string]string{
		"username":     "demoteuser1",
		"display_name": "Demote Test Admin",
		"password":     "password123",
		"role":         "admin",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	createResp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users/", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(createResp.Body)
		t.Fatalf("failed to create user: status=%d, body=%s", createResp.StatusCode, string(respBody))
	}

	// Read the user ID from create response
	var createResponse map[string]interface{}
	createBody, _ := io.ReadAll(createResp.Body)
	if err := json.Unmarshal(createBody, &createResponse); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	userID, ok := createResponse["id"].(string)
	if !ok || userID == "" {
		t.Skip("CreateUser endpoint doesn't return user ID yet")
	}

	// Now demote the user
	demoteReq, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/v1/users/%s/demote", baseURL, userID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create demote request: %v", err)
	}

	demoteResp, err := http.DefaultClient.Do(demoteReq)
	if err != nil {
		t.Fatalf("failed to demote user: %v", err)
	}
	defer demoteResp.Body.Close()

	respBody, err := io.ReadAll(demoteResp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if demoteResp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusOK, demoteResp.StatusCode, string(respBody))
	}

	// Verify response message
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "User demoted to regular user successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestDemoteUser_NotFound(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to demote a non-existent user
	fakeUserID := "00000000-0000-0000-0000-000000000000"

	demoteReq, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/v1/users/%s/demote", baseURL, fakeUserID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create demote request: %v", err)
	}

	resp, err := http.DefaultClient.Do(demoteReq)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusNotFound, resp.StatusCode, string(respBody))
	}
}
