package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestCreateUser_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// First, create and login an admin user to create another user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request to create a new user with a unique username (alphanumeric only)
	uniqueUsername := fmt.Sprintf("user%d", time.Now().UnixNano()%1000000)
	reqBody := map[string]string{
		"username":     uniqueUsername,
		"display_name": "Test User",
		"password":     "password123",
		"user_role":    "user",
	}
	_, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request with invalid username (too short)
	reqBody := map[string]string{
		"username":     "a", // min is 2
		"display_name": "Test User",
		"password":     "password123",
		"user_role":    "user",
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request with invalid password (too short)
	reqBody := map[string]string{
		"username":     "testuser2",
		"display_name": "Test User",
		"password":     "short", // min is 8
		"user_role":    "user",
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request with invalid role
	reqBody := map[string]string{
		"username":     "testuser3",
		"display_name": "Test User",
		"password":     "password123",
		"user_role":    "superadmin", // only "user" or "admin" allowed
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request with missing required fields
	reqBody := map[string]string{
		"username": "testuser4",
		// missing display_name, password, role
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Send invalid JSON with authorization header
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/users/", baseURL),
		bytes.NewReader([]byte(`{"username": "test", invalid json}`)),
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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

	// First, login an admin user
	adminToken := LoginUser(t, baseURL, "admin", "admin123")

	// Prepare request with admin role
	uniqueUsername := fmt.Sprintf("adminuser%d", time.Now().UnixNano()%1000000)
	reqBody := map[string]string{
		"username":     uniqueUsername,
		"display_name": "Admin User",
		"password":     "adminpass123",
		"user_role":    "admin",
	}

	// Make authorized request
	resp := MakeAuthorizedRequest(t, "POST", fmt.Sprintf("%s/api/v1/users/", baseURL), adminToken, reqBody)
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
