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

func TestGetLatestTelegramRecordsByTelegramID_Success(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user
	telegramID := uint64(10101010101)
	telegramUserReq := map[string]interface{}{
		"telegram_id": telegramID,
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

	// Add multiple telegram records for this user
	for i := 0; i < 3; i++ {
		recordReq := map[string]interface{}{
			"message_telegram_id":   uint64(100000000 + i),
			"from_user_telegram_id": userID,
			"in_telegram_chat_id":   int64(200000000 + i),
			"message_text":          fmt.Sprintf("Test message %d", i),
			"posted_at":             time.Now().Add(time.Duration(i) * time.Minute).Format(time.RFC3339),
		}
		recordResp := MakeAuthorizedRequest(
			t,
			"POST",
			fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL),
			token,
			recordReq,
		)
		recordResp.Body.Close()

		if recordResp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to create telegram record %d", i)
		}
	}

	// Now get the latest records
	getReq := map[string]interface{}{
		"telegram_id": telegramID,
	}

	resp := MakeAuthorizedRequest(
		t,
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID),
		token,
		getReq,
	)
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
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	records, ok := response["records"].([]interface{})
	if !ok {
		t.Fatalf("expected records array in response, got: %v", response)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	returnedTelegramID, ok := response["telegram_id"].(float64)
	if !ok || uint64(returnedTelegramID) != telegramID {
		t.Errorf("expected telegram_id %d, got %v", telegramID, response["telegram_id"])
	}
}

func TestGetLatestTelegramRecordsByTelegramID_NoRecords(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Use a telegram ID that doesn't exist
	telegramID := uint64(99999999999)

	getReq := map[string]interface{}{
		"telegram_id": telegramID,
	}

	resp := MakeAuthorizedRequest(
		t,
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID),
		token,
		getReq,
	)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert not found
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusNotFound, resp.StatusCode, string(respBody))
	}
}

func TestGetLatestTelegramRecordsByTelegramID_UserExistsNoRecords(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add a telegram user but no records
	telegramID := uint64(20202020202)
	telegramUserReq := map[string]interface{}{
		"telegram_id": telegramID,
	}
	userResp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		telegramUserReq,
	)
	userResp.Body.Close()

	if userResp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create telegram user")
	}

	// Try to get records (should fail - no records exist yet)
	getReq := map[string]interface{}{
		"telegram_id": telegramID,
	}

	resp := MakeAuthorizedRequest(
		t,
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID),
		token,
		getReq,
	)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert not found (no records for this telegram ID)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusNotFound, resp.StatusCode, string(respBody))
	}
}

func TestGetLatestTelegramRecordsByTelegramID_Unauthorized(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	// Don't provide a token
	telegramID := uint64(30303030303)

	reqBody := map[string]interface{}{
		"telegram_id": telegramID,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID),
		bytes.NewReader(bodyJSON),
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert unauthorized
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d or %d, got %d. Response: %s",
			http.StatusUnauthorized, http.StatusForbidden, resp.StatusCode, string(respBody))
	}
}

func TestGetLatestTelegramRecordsByTelegramID_InvalidFormat(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Send invalid request body
	getReq := map[string]interface{}{
		"invalid_field": "value",
	}

	resp := MakeAuthorizedRequest(t, "GET", fmt.Sprintf("%s/api/v1/record/telegram/0/records", baseURL), token, getReq)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Assert bad request or not found
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d. Response: %s",
			http.StatusBadRequest, http.StatusNotFound, resp.StatusCode, string(respBody))
	}
}

func TestGetLatestTelegramRecordsByTelegramID_MultipleUsers(t *testing.T) {
	baseURL, cleanup := StartTestServer(t)
	defer cleanup()

	token := LoginUser(t, baseURL, "admin", "admin123")

	// Add first telegram user and records
	telegramID1 := uint64(40404040404)
	user1Req := map[string]interface{}{
		"telegram_id": telegramID1,
	}
	user1Resp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		user1Req,
	)
	user1RespBody, _ := io.ReadAll(user1Resp.Body)
	user1Resp.Body.Close()

	var user1Response map[string]string
	json.Unmarshal(user1RespBody, &user1Response)
	userID1 := user1Response["record_id"]

	// Add record for user1
	record1Req := map[string]interface{}{
		"message_telegram_id":   uint64(400000000),
		"from_user_telegram_id": userID1,
		"in_telegram_chat_id":   int64(500000000),
		"message_text":          "User 1 message",
		"posted_at":             time.Now().Format(time.RFC3339),
	}
	record1Resp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL),
		token,
		record1Req,
	)
	record1Resp.Body.Close()

	// Add second telegram user and records
	telegramID2 := uint64(50505050505)
	user2Req := map[string]interface{}{
		"telegram_id": telegramID2,
	}
	user2Resp := MakeAuthorizedRequest(
		t,
		"POST",
		fmt.Sprintf("%s/api/v1/record/telegram/user", baseURL),
		token,
		user2Req,
	)
	user2RespBody, _ := io.ReadAll(user2Resp.Body)
	user2Resp.Body.Close()

	var user2Response map[string]string
	json.Unmarshal(user2RespBody, &user2Response)
	userID2 := user2Response["record_id"]

	// Add records for user2
	for i := 0; i < 2; i++ {
		record2Req := map[string]interface{}{
			"message_telegram_id":   uint64(600000000 + i),
			"from_user_telegram_id": userID2,
			"in_telegram_chat_id":   int64(700000000 + i),
			"message_text":          fmt.Sprintf("User 2 message %d", i),
			"posted_at":             time.Now().Add(time.Duration(i) * time.Minute).Format(time.RFC3339),
		}
		record2Resp := MakeAuthorizedRequest(
			t,
			"POST",
			fmt.Sprintf("%s/api/v1/record/telegram/record", baseURL),
			token,
			record2Req,
		)
		record2Resp.Body.Close()
	}

	// Get records for user1 - should only return user1's records
	get1Req := map[string]interface{}{
		"telegram_id": telegramID1,
	}
	resp1 := MakeAuthorizedRequest(
		t,
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID1),
		token,
		get1Req,
	)
	defer resp1.Body.Close()

	resp1Body, _ := io.ReadAll(resp1.Body)
	var response1 map[string]interface{}
	json.Unmarshal(resp1Body, &response1)

	records1, _ := response1["records"].([]interface{})
	if len(records1) != 1 {
		t.Errorf("expected 1 record for user1, got %d", len(records1))
	}

	// Get records for user2 - should only return user2's records
	get2Req := map[string]interface{}{
		"telegram_id": telegramID2,
	}
	resp2 := MakeAuthorizedRequest(
		t,
		"GET",
		fmt.Sprintf("%s/api/v1/record/telegram/%d/records", baseURL, telegramID2),
		token,
		get2Req,
	)
	defer resp2.Body.Close()

	resp2Body, _ := io.ReadAll(resp2.Body)
	var response2 map[string]interface{}
	json.Unmarshal(resp2Body, &response2)

	records2, _ := response2["records"].([]interface{})
	if len(records2) != 2 {
		t.Errorf("expected 2 records for user2, got %d", len(records2))
	}
}
