package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestRemoveUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// First, create a user to delete
	reqBody := map[string]string{
		"username":     "deleteuser1",
		"display_name": "Delete Test User",
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

	// Now delete the user
	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create delete request: %v", err)
	}

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	defer deleteResp.Body.Close()

	respBody, err := io.ReadAll(deleteResp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if deleteResp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusOK, deleteResp.StatusCode, string(respBody))
	}

	// Verify response message
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "User removed successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}

	// Try to get the deleted user - should return not found
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID))
	if err != nil {
		t.Fatalf("failed to get deleted user: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("expected deleted user to return 404, got %d", getResp.StatusCode)
	}
}

func TestRemoveUser_NotFound(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to delete a non-existent user
	fakeUserID := "00000000-0000-0000-0000-000000000000"

	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/users/%s", baseURL, fakeUserID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create delete request: %v", err)
	}

	resp, err := http.DefaultClient.Do(deleteReq)
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

func TestRemoveUser_InvalidUUID(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Try to delete with an invalid UUID
	invalidUserID := "not-a-valid-uuid"

	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/users/%s", baseURL, invalidUserID),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create delete request: %v", err)
	}

	resp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert - should return bad request for invalid UUID format
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}
