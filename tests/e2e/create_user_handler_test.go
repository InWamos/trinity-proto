package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request
	reqBody := map[string]string{
		"username":     "testuser",
		"display_name": "Test User",
		"password":     "password123",
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
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

	// Verify response message
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "The user has been created. You can login now"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestCreateUser_InvalidUsername_TooShort(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request with invalid username (too short)
	reqBody := map[string]string{
		"username":     "a", // min is 2
		"display_name": "Test User",
		"password":     "password123",
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestCreateUser_InvalidPassword_TooShort(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request with invalid password (too short)
	reqBody := map[string]string{
		"username":     "testuser2",
		"display_name": "Test User",
		"password":     "short", // min is 8
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request with invalid role
	reqBody := map[string]string{
		"username":     "testuser3",
		"display_name": "Test User",
		"password":     "password123",
		"role":         "superadmin", // only "user" or "admin" allowed
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request with missing required fields
	reqBody := map[string]string{
		"username": "testuser4",
		// missing display_name, password, role
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Send invalid JSON
	body := []byte(`{"username": "test", invalid json}`)

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"expected status %d, got %d. Response: %s",
			http.StatusBadRequest,
			resp.StatusCode,
			string(respBody),
		)
	}
}

func TestCreateUser_AdminRole(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Prepare request with admin role
	reqBody := map[string]string{
		"username":     "adminuser",
		"display_name": "Admin User",
		"password":     "adminpass123",
		"role":         "admin",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
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
}
