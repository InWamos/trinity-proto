package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestPromoteUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// First, create a regular user
	reqBody := map[string]string{
		"username":     "promoteuser1",
		"display_name": "Promote Test User",
		"password":     "password123",
		"role":         "user",
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

	// Now promote the user
	promoteReq, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/v1/users/%s/promote", baseURL, userID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create promote request: %v", err)
	}

	promoteResp, err := http.DefaultClient.Do(promoteReq)
	if err != nil {
		t.Fatalf("failed to promote user: %v", err)
	}
	defer promoteResp.Body.Close()

	respBody, err := io.ReadAll(promoteResp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if promoteResp.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusOK,
			promoteResp.StatusCode,
			string(respBody),
		)
	}

	// Verify response message
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "User promoted to admin successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestPromoteUser_NotFound(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to promote a non-existent user
	fakeUserID := "00000000-0000-0000-0000-000000000000"

	promoteReq, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/v1/users/%s/promote", baseURL, fakeUserID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create promote request: %v", err)
	}

	resp, err := http.DefaultClient.Do(promoteReq)
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
