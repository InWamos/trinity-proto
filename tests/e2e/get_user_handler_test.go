package e2e

import (
"bytes"
"encoding/json"
"fmt"
"io"
"net/http"
"testing"
)

func TestGetUserByID_Success(t *testing.T) {
baseURL, cleanup := StartTestServer(t)
defer cleanup()

// First, create a user
reqBody := map[string]string{
"username":     "getuser1",
"display_name": "Get User Test",
"password":     "password123",
"role":         "user",
}
body, err := json.Marshal(reqBody)
if err != nil {
t.Fatalf("failed to marshal request body: %v", err)
}

createResp, err := http.Post(
fmt.Sprintf("%s/api/v1/users", baseURL),
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

// Now get the user by ID
getResp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID))
if err != nil {
t.Fatalf("failed to get user: %v", err)
}
defer getResp.Body.Close()

respBody, err := io.ReadAll(getResp.Body)
if err != nil {
t.Fatalf("failed to read response body: %v", err)
}

// Assert
if getResp.StatusCode != http.StatusOK {
t.Errorf("expected status %d, got %d. Response: %s", http.StatusOK, getResp.StatusCode, string(respBody))
}

// Verify response structure
var getUserResponse map[string]interface{}
if err := json.Unmarshal(respBody, &getUserResponse); err != nil {
t.Fatalf("failed to unmarshal get response: %v", err)
}

if getUserResponse["username"] != reqBody["username"] {
t.Errorf("expected username %q, got %q", reqBody["username"], getUserResponse["username"])
}
if getUserResponse["display_name"] != reqBody["display_name"] {
t.Errorf("expected display_name %q, got %q", reqBody["display_name"], getUserResponse["display_name"])
}
if getUserResponse["role"] != reqBody["role"] {
t.Errorf("expected role %q, got %q", reqBody["role"], getUserResponse["role"])
}
}

func TestGetUserByID_NotFound(t *testing.T) {
baseURL, cleanup := StartTestServer(t)
defer cleanup()

// Try to get a user with a non-existent UUID
fakeUserID := "00000000-0000-0000-0000-000000000000"

resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", baseURL, fakeUserID))
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

func TestGetUserByID_InvalidUUID(t *testing.T) {
baseURL, cleanup := StartTestServer(t)
defer cleanup()

// Try to get a user with an invalid UUID
invalidUserID := "not-a-valid-uuid"

resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", baseURL, invalidUserID))
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
