package fmgcli

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestLogin_Success_UserLogin(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{ // Change to []interface{} to match decoding behavior
			map[string]interface{}{
				"data": map[string]interface{}{
					"user":     "fake-user",
					"password": "fake-password",
				},
				"url": "/sys/login/user",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/sys/login/user",
			},
		},
		"session": "fake-session",
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")

	// Call the Login method
	err := client.Login()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the session
	if client.Session != "fake-session" {
		t.Errorf("Expected session 'fake-session', got %s", client.Session)
	}
}

func TestLogin_Failure_InvalidUserCredentials(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{ // Change to []interface{} to match decoding behavior
			map[string]interface{}{
				"data": map[string]interface{}{
					"user":     "fake-invalid-user",
					"password": "fake-password",
				},
				"url": "/sys/login/user",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -22,
					"message": "Login fail",
				},
				"url": "/sys/login/user",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-invalid-user", "fake-password")

	// Call the Login method
	err := client.Login()
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	// Validate the session
	if client.Session != "" {
		t.Errorf("Expected session to be empty, got %s", client.Session)
	}

	// Validate the error message
	expectedErrorMessage := "login failed"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestLogin_Success_NoSecretsInLogsDuringUserLogin(t *testing.T) {
	t.Skip("temporarily disabled")
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"user":     "fake-user",
					"password": "fake-password",
				},
				"url": "/sys/login/user",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/sys/login/user",
			},
		},
		"session": "fake-session",
	}

	// Var to capture logs
	var logBuffer bytes.Buffer

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")

	// Capture the logs
	//client.infoLog.SetOutput(&logBuffer) //this worked last time with this code
	client.log = slog.New(slog.NewTextHandler(&logBuffer, nil))

	// Call the Login method
	err := client.Login()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the session
	if client.Session != "fake-session" {
		t.Errorf("Expected session 'fake-session', got %s", client.Session)
	}

	// Check logs for secrets
	logContent := logBuffer.String()
	if strings.Contains(logContent, "fake-password") || strings.Contains(logContent, "fake-session") {
		t.Errorf("Logs contain sensitive information: %s", logContent)
	}
}

func TestLogin_Failure_EmptyUserName(t *testing.T) {
	// Create a new FortiManager client with empty user info
	client := NewUserClient("http://fake-host", "", "fake-password") // Empty user

	// Call the Login method
	err := client.Login()

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "username or password is empty"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestLogin_Failure_EmptyPw(t *testing.T) {
	// Create a new FortiManager client with empty user info
	client := NewUserClient("http://fake-host", "fake-user", "") // Empty password

	// Call the Login method
	err := client.Login()

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "username or password is empty"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestLogout_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/sys/logout",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/sys/logout",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Logout method
	err := client.Logout()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the session is cleared
	if client.Session != "" {
		t.Errorf("Expected session to be empty, got %s", client.Session)
	}
}

func TestLogout_Failure_AlreadyLoggedOutWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/sys/logout",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/sys/logout",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = ""

	// Call the Logout method
	err := client.Logout()
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "logout failed"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeleteFromPolicy_SuccessRemovingSrcsWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"object1", "object2"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/srcaddr",
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"object1", "object2"},
		"srcaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteFromPolicy_SuccessRemovingSrcsWithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"object1", "object2"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/srcaddr",
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"object1", "object2"},
		"srcaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteFromPolicy_SuccessRemovingDstsWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"object1", "object2"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/dstaddr",
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"object1", "object2"},
		"dstaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteFromPolicy_SuccessRemovingDstsWithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"object1", "object2"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/dstaddr",
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"object1", "object2"},
		"dstaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteFromPolicy_Failure_InvalidRoleWithUserClient(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("The server should not be called for an invalid role")
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteFromPolicy method with an invalid role
	err := client.DeleteFromPolicy(
		12345,
		[]string{"object1", "object2"},
		"invalid-role", // Invalid role
		"test-vdom",
		"test-device",
		"test-adom",
	)

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "invalid role"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeleteFromPolicy_Failure_LastSrcObjectWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"last-object"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/srcaddr",
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure due to last object removal
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -1,
					"message": "firewall/policy/12345/ : runtime error -2: srcaddr in Policy \"12345\" Package \"test-vdom\" cannot be empty if dstaddr is set",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/srcaddr",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"last-object"},
		"srcaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to delete objects: firewall/policy/12345/ : runtime error -2: srcaddr in Policy \"12345\" Package \"test-vdom\" cannot be empty if dstaddr is set"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeleteFromPolicy_Failure_LastDstObjectWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": map[string]interface{}{
			"data": []interface{}{"last-object"},
			"url":  "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/dstaddr",
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure due to last object removal
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -1,
					"message": "firewall/policy/12345/ : runtime error -2: dstaddr in Policy \"12345\" Package \"test-vdom\" cannot be empty if srcaddr is set",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345/dstaddr",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteFromPolicy method
	err := client.DeleteFromPolicy(
		12345,
		[]string{"last-object"},
		"dstaddr",
		"test-vdom",
		"test-device",
		"test-adom",
	)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to delete objects: firewall/policy/12345/ : runtime error -2: dstaddr in Policy \"12345\" Package \"test-vdom\" cannot be empty if srcaddr is set"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDisablePolicies_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			{
				"data": map[string]interface{}{
					"policyid": 12346,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DisablePolicy method
	err := client.DisablePolicies([]int{12345, 12346}, "test-vdom", "test-device", "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDisablePolicies_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			{
				"data": map[string]interface{}{
					"policyid": 12346,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DisablePolicy method
	err := client.DisablePolicies([]int{12345, 12346}, "test-vdom", "test-device", "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDisablePolicies_Failure_InexistentPolicyWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			{
				"status": map[string]interface{}{
					"code":    -9998,
					"message": "firewall/policy/12345/ : srcintf in Policy \"12345\" Package \"v-test\" cannot be empty",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-policy-package/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DisablePolicy method
	err := client.DisablePolicies([]int{12345, 12346}, "test-vdom", "test-device", "test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to disable the following policies: [12346]"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDisablePolicy_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DisablePolicy method
	err := client.DisablePolicy("test-adom", "test-device/test-vdom", 12345)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDisablePolicy_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DisablePolicy method
	err := client.DisablePolicy("test-adom", "test-device/test-vdom", 12345)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDisablePolicy_Failure_InexistentPolicyWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"status": "disable",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -9998,
					"message": "firewall/policy/12345/ : srcintf in Policy \"12345\" Package \"v-test\" cannot be empty",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-policy-package/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DisablePolicy method
	err := client.DisablePolicy("test-adom", "test-device/test-vdom", 12345)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to disable policy: firewall/policy/12345/ : srcintf in Policy \"12345\" Package \"v-test\" cannot be empty"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeletePolicy_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeletePolicy method
	err := client.DeletePolicy("test-adom", "test-device/test-vdom", 85)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeletePolicy_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DeletePolicy method
	err := client.DeletePolicy("test-adom", "test-device/test-vdom", 85)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeletePolicy_Failure_InexistentPolicyWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "entry not found",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/85",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeletePolicy method
	err := client.DeletePolicy("test-adom", "test-device/test-vdom", 85)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to delete policy: entry not found"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCreatePolicy_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"action":     "accept",
					"srcintf":    []interface{}{"any"},
					"dstintf":    []interface{}{"any"},
					"srcaddr":    []interface{}{"testobj1"},
					"dstaddr":    []interface{}{"testobj2", "testobj3"},
					"service":    []interface{}{"tcp-443"},
					"status":     "enable",
					"schedule":   "always",
					"logtraffic": "all",
					"comments":   "tHis Is a rANdoM cOmmENt",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the CreatePolicy method
	policyID, err := client.CreatePolicy(
		"test-device/test-vdom",
		"test-adom",
		"tHis Is a rANdoM cOmmENt",
		[]string{"testobj1"},
		[]string{"testobj2", "testobj3"},
		[]string{"tcp-443"},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policyID != 12345 {
		t.Fatalf("Expected policy ID 12345, got %d", policyID)
	}
}

func TestCreatePolicy_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"action":     "accept",
					"srcintf":    []interface{}{"any"},
					"dstintf":    []interface{}{"any"},
					"srcaddr":    []interface{}{"testobj1"},
					"dstaddr":    []interface{}{"testobj2", "testobj3"},
					"service":    []interface{}{"tcp-443"},
					"status":     "enable",
					"schedule":   "always",
					"logtraffic": "all",
					"comments":   "tHis Is a rANdoM cOmmENt",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreatePolicy method
	policyID, err := client.CreatePolicy(
		"test-device/test-vdom",
		"test-adom",
		"tHis Is a rANdoM cOmmENt",
		[]string{"testobj1"},
		[]string{"testobj2", "testobj3"},
		[]string{"tcp-443"},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policyID != 12345 {
		t.Fatalf("Expected policy ID 12345, got %d", policyID)
	}
}

func TestCreatePolicy_Success_WithAPIClientAndMetafields(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"action":     "accept",
					"srcintf":    []interface{}{"any"},
					"dstintf":    []interface{}{"any"},
					"srcaddr":    []interface{}{"testobj1"},
					"dstaddr":    []interface{}{"testobj2", "testobj3"},
					"service":    []interface{}{"tcp-443"},
					"status":     "enable",
					"schedule":   "always",
					"logtraffic": "all",
					"meta fields": map[string]interface{}{
						"some_uuid":    "821c2a23-3c4f-4b1d-8e5f-4d8e175af6f9",
						"some_randid":  "abc",
						"another_uuid": "4f8b2c1e-3d5f-4b1d-8e5f-4d8e175af6f9",
						"third_uuid":   "ABCD-1234",
					},
					"comments": "tHis Is a rANdoM cOmmENt",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreatePolicy method
	policyID, err := client.CreatePolicy(
		"test-device/test-vdom",
		"test-adom",
		"tHis Is a rANdoM cOmmENt",
		[]string{"testobj1"},
		[]string{"testobj2", "testobj3"},
		[]string{"tcp-443"},
		WithMetafields(map[string]interface{}{
			"some_uuid":    "821c2a23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid":  "abc",
			"another_uuid": "4f8b2c1e-3d5f-4b1d-8e5f-4d8e175af6f9",
			"third_uuid":   "ABCD-1234",
		}),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policyID != 12345 {
		t.Fatalf("Expected policy ID 12345, got %d", policyID)
	}
}

func TestCreatePolicy_Failure_InexistentSrcAddrWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"action":     "accept",
					"srcintf":    []interface{}{"any"},
					"dstintf":    []interface{}{"any"},
					"srcaddr":    []interface{}{"testobj1"},
					"dstaddr":    []interface{}{"testobj2", "testobj3"},
					"service":    []interface{}{"tcp-443"},
					"status":     "enable",
					"schedule":   "always",
					"logtraffic": "all",
					"comments":   "tHis Is a rANdoM cOmmENt",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    -10131,
					"message": "datasrc invalid. object: policy package \"test-vdom\" - firewall policy.12345:srcaddr. detail: testobj1. solution: data not exist",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the CreatePolicy method
	_, err := client.CreatePolicy(
		"test-device/test-vdom",
		"test-adom",
		"tHis Is a rANdoM cOmmENt",
		[]string{"testobj1"},
		[]string{"testobj2", "testobj3"},
		[]string{"tcp-443"},
	)

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to create policy: datasrc invalid. object: policy package \"test-vdom\" - firewall policy.12345:srcaddr. detail: testobj1. solution: data not exist"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCreatePolicy_Failure_InexistentSrcAddrWithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": map[string]interface{}{
					"action":     "accept",
					"srcintf":    []interface{}{"any"},
					"dstintf":    []interface{}{"any"},
					"srcaddr":    []interface{}{"testobj1"},
					"dstaddr":    []interface{}{"testobj2", "testobj3"},
					"service":    []interface{}{"tcp-443"},
					"status":     "enable",
					"schedule":   "always",
					"logtraffic": "all",
					"comments":   "tHis Is a rANdoM cOmmENt",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"policyid": 12345,
				},
				"status": map[string]interface{}{
					"code":    -10131,
					"message": "datasrc invalid. object: policy package \"test-vdom\" - firewall policy.12345:srcaddr. detail: testobj1. solution: data not exist",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreatePolicy method
	_, err := client.CreatePolicy(
		"test-device/test-vdom",
		"test-adom",
		"tHis Is a rANdoM cOmmENt",
		[]string{"testobj1"},
		[]string{"testobj2", "testobj3"},
		[]string{"tcp-443"},
	)

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to create policy: datasrc invalid. object: policy package \"test-vdom\" - firewall policy.12345:srcaddr. detail: testobj1. solution: data not exist"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCommit_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/commit",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Commit method
	err := client.Commit("test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommit_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/commit",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the Commit method
	err := client.Commit("test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommit_Failure_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/commit",
			},
		},
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -1,
					"message": "Commit failed due to invalid configuration",
				},
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the Commit method
	err := client.Commit("test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to save ADOM: Commit failed due to invalid configuration"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByName_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj2",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []string{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj2",
					"oid":             6292,
					"subnet":          []string{"123.123.123.124", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj2",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetAddresses method
	addresses, err := client.GetAddressesByName([]string{"testobj1", "testobj2"}, "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(addresses) != 2 {
		t.Fatalf("Expected 2 addresses, got %d", len(addresses))
	}

	if addresses[0].Name != "testobj1" || addresses[1].Name != "testobj2" {
		t.Errorf("Unexpected address names: %v", addresses)
	}
}

func TestGetAddressesByName_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj2",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []string{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj2",
					"oid":             6292,
					"subnet":          []string{"123.123.123.124", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj2",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetAddresses method
	addresses, err := client.GetAddressesByName([]string{"testobj1", "testobj2"}, "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(addresses) != 2 {
		t.Fatalf("Expected 2 addresses, got %d", len(addresses))
	}

	if addresses[0].Name != "testobj1" || addresses[1].Name != "testobj2" {
		t.Errorf("Unexpected address names: %v", addresses)
	}
}

func TestGetAddressesByName_Failure_OneExistsAnotherDoesntWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/inexistent",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testObj",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/inexistent",
			},
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testObj",
					"oid":             6292,
					"subnet":          []string{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testObj",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetAddresses method
	addresses, err := client.GetAddressesByName([]string{"inexistent", "testObj"}, "test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the response
	if len(addresses) != 0 {
		t.Fatalf("Expected 0 addresses, got %d", len(addresses))
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following objects"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByName_Failure_OneExistsTwoDontWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/inexistent1",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/inexistent2",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testObj",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/inexistent1",
			},
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/inexistent2",
			},
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testObj",
					"oid":             6292,
					"subnet":          []string{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testObj",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetAddresses method
	addresses, err := client.GetAddressesByName([]string{"inexistent1", "inexistent2", "testObj"}, "test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the response
	if len(addresses) != 0 {
		t.Fatalf("Expected 0 addresses, got %d", len(addresses))
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following objects"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByMetafield_Success_StringMetafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
						},
					},
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "other-uuid",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByMetafield("test-adom", "some_uuid", "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatalf("Expected service, got nil")
	}

	if service.Name != "tcp-443" {
		t.Errorf("Expected service name 'tcp-443', got '%s'", service.Name)
	}

	if !reflect.DeepEqual(service.TCPPortRange, []interface{}{"443"}) &&
		!reflect.DeepEqual(service.TCPPortRange, []string{"443"}) {
		t.Errorf("Expected TCP port range [443], got %v", service.TCPPortRange)
	}

	if service.Metafields["some_uuid"] != "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9" {
		t.Errorf("Expected meta field 'some_uuid' to be '821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9', got %v", service.Metafields["some_uuid"])
	}
}

func TestGetServiceByMetafield_Success_Float64Metafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": float64(123456789),
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": float64(1234),
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByMetafield("test-adom", "some_uuid", float64(123456789))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatalf("Expected service, got nil")
	}

	if service.Name != "tcp-8080" {
		t.Errorf("Expected service name 'tcp-8080', got '%s'", service.Name)
	}

	if !reflect.DeepEqual(service.TCPPortRange, []interface{}{"8080"}) &&
		!reflect.DeepEqual(service.TCPPortRange, []string{"8080"}) {
		t.Errorf("Expected TCP port range [8080], got %v", service.TCPPortRange)
	}

	if service.Metafields["some_uuid"] != float64(123456789) {
		t.Errorf("Expected meta field 'some_uuid' to be 123456789, got %v", service.Metafields["some_uuid"])
	}
}

func TestGetServiceByMetafield_Failure_MultipleMatches(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"duplicate-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "multiple services found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByMetafield_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"no-permission-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "error fetching services: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByMetafield_Failure_NoServices(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"no-match-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "found 0 services corresponding the following 'some_uuid' metafield values: '[no-match-value]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByMetafield_Failure_NoMatchingService(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body: no service matches the metafield value
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "other-uuid",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "another-uuid",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	// Try to get a service by a metafield value that does not exist
	service, err := client.GetServiceByMetafield("test-adom", "some_uuid", "not-matched-uuid")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "no service found with metafield 'some_uuid'='not-matched-uuid' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByMetafield_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "2b1c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"1b1c2d3e-5678-90ab-cdef-1234567890ab", "2b1c2d3e-5678-90ab-cdef-1234567890ab"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "tcp-8080" || services[1].Name != "tcp-443" {
		t.Errorf("Unexpected service names: %v", services)
	}
}

func TestGetServicesByMetafield_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "2b1c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"1b1c2d3e-5678-90ab-cdef-1234567890ab", "2b1c2d3e-5678-90ab-cdef-1234567890ab"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "tcp-8080" || services[1].Name != "tcp-443" {
		t.Errorf("Unexpected service names: %v", services)
	}
}

func TestGetServicesByMetafield_Success_Float64Metafield(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": float64(123456789),
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": float64(1234),
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{float64(123456789)})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(services) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(services))
	}

	if services[0].Name != "tcp-8080" {
		t.Errorf("Expected service name 'tcp-8080', got '%s'", services[0].Name)
	}

	if services[0].Metafields["some_uuid"] != float64(123456789) {
		t.Errorf("Expected meta field 'some_uuid' to be 123456789, got %v", services[0].Metafields["some_uuid"])
	}
}

func TestGetServicesByMetafield_Failure_MultipleMatches(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"duplicate-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "multiple services found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByMetafield_Failure_NoPermission(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"no-permission-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "error fetching services: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByMetafield_Failure_NoServices(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"no-match-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "found 0 services corresponding the following 'some_uuid' metafield values: '[no-match-value]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByMetafield_Failure_NoMatchingServices(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"random-not-matched-uuid"})
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "found 0 services corresponding the following 'some_uuid' metafield values: '[random-not-matched-uuid]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByMetafield_Failure_NotAllMatching(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"name":          "tcp-8080",
						"obj seq":       92,
						"oid":           5450,
						"tcp-portrange": []interface{}{"8080"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"name":          "tcp-443",
						"obj seq":       101,
						"oid":           5519,
						"tcp-portrange": []interface{}{"443"},
						"udp-portrange": []interface{}{},
						"unset attrs":   []interface{}{"icmptype", "icmpcode"},
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByMetafield("test-adom", "some_uuid", []interface{}{"no-permission-value", "random-uuid-1"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if services != nil {
		t.Fatalf("Expected services to be nil, got %v", services)
	}

	expectedErrorMessage := "found 1 services corresponding the following 'some_uuid' metafield values: '[no-permission-value random-uuid-1]', but expected 2 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByName_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"tcp-portrange": []interface{}{"8080"},
					"udp-portrange": []interface{}{},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
			{
				"data": map[string]interface{}{
					"name":          "tcp-443",
					"obj seq":       101,
					"oid":           5519,
					"tcp-portrange": []interface{}{"443"},
					"udp-portrange": []interface{}{},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetServices method
	services, err := client.GetServicesByName([]string{"tcp-8080", "tcp-443"}, "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "tcp-8080" || services[1].Name != "tcp-443" {
		t.Errorf("Unexpected service names: %v", services)
	}
}

func TestGetServicesByName_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"tcp-portrange": []interface{}{"8080"},
					"udp-portrange": []interface{}{},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
			{
				"data": map[string]interface{}{
					"name":          "tcp-443",
					"obj seq":       101,
					"oid":           5519,
					"tcp-portrange": []interface{}{"443"},
					"udp-portrange": []interface{}{},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetServices method
	services, err := client.GetServicesByName([]string{"tcp-8080", "tcp-443"}, "test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "tcp-8080" || services[1].Name != "tcp-443" {
		t.Errorf("Unexpected service names: %v", services)
	}
}

func TestGetServicesByName_Failure_OneExistsAnotherDoesntWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent",
			},
			{
				"data": map[string]interface{}{
					"name":    "tcp-443",
					"obj seq": 101,
					"oid":     5519,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetServices method
	services, err := client.GetServicesByName([]string{"inexistent", "tcp-443"}, "test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the response
	if len(services) != 0 {
		t.Fatalf("Expected 0 services, got %d", len(services))
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following services"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServicesByName_Failure_OneExistsTwoDontWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent1",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent2",
			},
			map[string]interface{}{
				"fields": []interface{}{"name", "tcp-portrange", "udp-portrange"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent1",
			},
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent2",
			},
			{
				"data": map[string]interface{}{
					"name":    "tcp-443",
					"obj seq": 101,
					"oid":     5519,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-443",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetServices method
	services, err := client.GetServicesByName([]string{"inexistent1", "inexistent2", "tcp-443"}, "test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the response
	if len(services) != 0 {
		t.Fatalf("Expected 0 services, got %d", len(services))
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following services"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestLock_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Lock method
	err := client.Lock("test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestLock_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the Lock method
	err := client.Lock("test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestLock_Failure_WithUserClientWorkspaceLockedByOtherUser(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -20055,
					"message": "Workspace is locked by other user",
				},
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Lock method
	err := client.Lock("test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to lock ADOM: Workspace is locked by other user"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestLock_Failure_NoPermissionWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/dvmdb/adom/test-adom/workspace/lock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Lock method
	err := client.Lock("test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to lock ADOM: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestUnlock_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/unlock",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/dvmdb/adom/test-adom/workspace/unlock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Unlock method
	err := client.Unlock("test-adom")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUnlock_Failure_WithUserClientNoPermission(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "exec",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/dvmdb/adom/test-adom/workspace/unlock",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/dvmdb/adom/test-adom/workspace/unlock",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the Unlock method
	err := client.Unlock("test-adom")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to unlock ADOM: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPolicyByMetafield_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "1b9c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetPolicyByMetafield method
	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", "1b9c2d3e-5678-90ab-cdef-1234567890ab")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if policy == nil {
		t.Fatalf("Expected a policy, got nil")
	}

	if policy.PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policy.PolicyID)
	}
}

func TestGetPolicyByMetafield_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "1b9c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPolicyByMetafield method
	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", "1b9c2d3e-5678-90ab-cdef-1234567890ab")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if policy == nil {
		t.Fatalf("Expected a policy, got nil")
	}

	if policy.PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policy.PolicyID)
	}
}

func TestGetPolicyByMetafield_Success_Float64Metafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": float64(123456789),
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", float64(123456789))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policy == nil {
		t.Fatalf("Expected a policy, got nil")
	}

	if policy.PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policy.PolicyID)
	}

	if policy.Metafields["some_uuid"] != float64(123456789) {
		t.Errorf("Expected meta field 'some_uuid' to be 123456789, got %v", policy.Metafields["some_uuid"])
	}
}

func TestGetPolicyByMetafield_Failure_MultipleMatches(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", "duplicate-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if policy != nil {
		t.Fatalf("Expected policy to be nil, got %v", policy)
	}
	expectedErrorMessage := "multiple policies found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPolicyByMetafield_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", "no-permission-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if policy != nil {
		t.Fatalf("Expected policy to be nil, got %v", policy)
	}
	expectedErrorMessage := "error fetching policies: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPolicyByMetafield_Failure_NoPolicies(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data":   []map[string]interface{}{},
				"status": map[string]interface{}{"code": float64(0), "message": "OK"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	policy, err := client.GetPolicyByMetafield("test-adom", "test-device/test-vdom", "some_uuid", "no-match-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if policy != nil {
		t.Fatalf("Expected policy to be nil, got %v", policy)
	}
	expectedErrorMessage := "no policy found with metafield 'some_uuid'='no-match-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Failure_NoMatchingPolicies(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method with a value that does not match any policy
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"random-not-matched-uuid"})
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "found 0 policies corresponding the following 'some_uuid' metafield values: '[random-not-matched-uuid]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByID_Success(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"_created timestamp":  1745786504,
					"_created-by":         "fake-user",
					"_last-modified-by":   "fake-user",
					"_last_hit":           0,
					"_modified timestamp": 1745796504,
					"action":              "accept",
					"dstaddr":             []interface{}{"duplicate2.vt.ch"},
					"meta fields": map[string]interface{}{
						"some_uuid":    "821c2a23-3c4f-4b1d-8e5f-4d8e175af6f9",
						"some_randid":  "abc",
						"another_uuid": "4f8b2c1e-3d5f-4b1d-8e5f-4d8e175af6f9",
						"third_uuid":   "ABCD-1234",
					},
					"obj seq":      60,
					"obj ver":      4,
					"oid":          6611,
					"policyid":     12345,
					"schedule":     []interface{}{"always"},
					"service":      []interface{}{"tcp-443"},
					"srcaddr":      []interface{}{"duplicate1.vt.ch"},
					"status":       "enable",
					"vpn_dst_node": nil,
					"vpn_src_node": nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPolicies method
	policies, err := client.GetPoliciesByID("test-device/test-vdom", "test-adom", []int{12345})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(policies) != 1 {
		t.Fatalf("Expected 1 policy, got %d", len(policies))
	}

	if policies[0].PolicyID != 12345 {
		t.Errorf("Unexpected policy ID: %v", policies[0].PolicyID)
	}
}

func TestGetPoliciesByID_Failure_PolicyDoesntExist(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPolicies method
	_, err := client.GetPoliciesByID("test-device/test-vdom", "test-adom", []int{12345})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following policies: [12345]"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByID_Failure_OnePolicyDoesntExistTheOtherDoes(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			map[string]interface{}{
				"fields": []interface{}{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"_created timestamp":  1745786504,
					"_created-by":         "fake-user",
					"_last-modified-by":   "fake-user",
					"_last_hit":           0,
					"_modified timestamp": 1745796504,
					"action":              "accept",
					"dstaddr":             []interface{}{"duplicate2.vt.ch"},
					"meta fields": map[string]interface{}{
						"some_uuid":    "821c2a23-3c4f-4b1d-8e5f-4d8e175af6f9",
						"some_randid":  "abc",
						"another_uuid": "4f8b2c1e-3d5f-4b1d-8e5f-4d8e175af6f9",
						"third_uuid":   "ABCD-1234",
					},
					"obj seq":      60,
					"obj ver":      4,
					"oid":          6611,
					"policyid":     12345,
					"schedule":     []interface{}{"always"},
					"service":      []interface{}{"tcp-443"},
					"srcaddr":      []interface{}{"duplicate1.vt.ch"},
					"status":       "enable",
					"vpn_dst_node": nil,
					"vpn_src_node": nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12345",
			},
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy/12346",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPolicies method
	_, err := client.GetPoliciesByID("test-device/test-vdom", "test-adom", []int{12345, 12346})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "there were errors when fetching the following policies: [12346]"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"session": "fake-session",
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "1b9c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"1b9c2d3e-5678-90ab-cdef-1234567890ab"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(policies) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(policies))
	}

	if policies[0].PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policies[0].PolicyID)
	}
}

func TestGetPoliciesByMetafield_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "1b9c2d3e-5678-90ab-cdef-1234567890ab",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"1b9c2d3e-5678-90ab-cdef-1234567890ab"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(policies) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(policies))
	}

	if policies[0].PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policies[0].PolicyID)
	}
}

func TestGetPoliciesByMetafield_Success_Float64Metafield(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": float64(123456789),
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": float64(0),
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{float64(123456789)})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(policies) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(policies))
	}

	if policies[0].PolicyID != 32 {
		t.Errorf("Unexpected policy ID: %v", policies[0].PolicyID)
	}
}

func TestGetPoliciesByMetafield_Failure_MultipleMatches(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"duplicate-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "multiple policies found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Failure_NoPermission(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"no-permission-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "error fetching policies: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Failure_NoServices(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"no-match-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "found 0 policies corresponding the following 'some_uuid' metafield values: '[no-match-value]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Failure_NoMatchingServices(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"random-not-matched-uuid"})
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "found 0 policies corresponding the following 'some_uuid' metafield values: '[random-not-matched-uuid]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetPoliciesByMetafield_Failure_NotAllMatching(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{
					"obj seq",
					"status",
					"policyid",
					"srcaddr",
					"dstaddr",
					"service",
					"action",
					"schedule",
					"extra info",
					"_last_hit",
				},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
		"verbose": float64(1),
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"_created timestamp":  float64(1710346227),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716275354),
						"obj seq":             float64(1),
						"obj ver":             float64(1),
						"oid":                 float64(6051),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(32),
						"srcaddr":             []interface{}{"testobj1"},
						"dstaddr":             []interface{}{"testobj2"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"udp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"_created timestamp":  float64(1716538811),
						"_last-modified-by":   "v006951",
						"_modified timestamp": float64(1716538811),
						"obj seq":             float64(2),
						"obj ver":             float64(1),
						"oid":                 float64(6183),
						"vpn_dst_node":        nil,
						"vpn_src_node":        nil,
						"policyid":            float64(49),
						"srcaddr":             []interface{}{"testobj3"},
						"dstaddr":             []interface{}{"testobj4"},
						"action":              "accept",
						"status":              "enable",
						"schedule":            []interface{}{"always"},
						"service":             []interface{}{"tcp-443"},
						"_last_hit":           float64(0),
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    float64(0),
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/pkg/test-device/test-vdom/firewall/policy",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the GetPoliciesByMetafield method
	policies, err := client.GetPoliciesByMetafield("test-adom", "test-device/test-vdom", "some_uuid", []interface{}{"no-permission-value", "random-uuid-1"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if policies != nil {
		t.Fatalf("Expected policies to be nil, got %v", policies)
	}

	expectedErrorMessage := "found 1 policies corresponding the following 'some_uuid' metafield values: '[no-permission-value random-uuid-1]', but expected 2 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCreateService_Success_AddingTCPService(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":          "tcp-8080",
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080-8080"},
						"comment":       "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "tcp-8080",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateService method
	err := client.CreateService("test-adom", "tcp-8080", "tcp", 8080, 8080, "test comment")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateService_Success_AddingUDPService(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":          "udp-8080",
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"udp-portrange": []interface{}{"8080-8080"},
						"comment":       "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "udp-8080",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateService method
	err := client.CreateService("test-adom", "udp-8080", "udp", 8080, 8080, "test comment")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateAddressGroup_Success_AddingGroup(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":        "test-group",
						"member":      []interface{}{"addr-1", "addr-2"},
						"comment":     "test comment",
						"meta fields": map[string]interface{}{"some_uuid": "random-uuid"},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "test-group",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	err := client.CreateAddressGroup("test-adom", "test-group", []string{"addr-1", "addr-2"}, "test comment", WithAddressGroupMetafields(map[string]interface{}{"some_uuid": "random-uuid"}))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateAddressGroupWithMetafields_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":        "test-group",
						"meta fields": map[string]interface{}{"some_uuid": "random-uuid"},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "test-group",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp/test-group",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	err := client.UpdateAddressGroupWithMetafields("test-adom", "test-group", map[string]interface{}{"some_uuid": "random-uuid"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteAddressGroup_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp/test-group",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/addrgrp/test-group",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	err := client.DeleteAddressGroup("test-adom", "test-group")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateService_Failure_InvalidPortRange(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":          "tcp-8080",
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080-8080"},
						"comment":       "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "tcp-8080",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateService method
	err := client.CreateService("test-adom", "tcp-8080", "tcp", 8080, 8079, "test comment")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "invalid port range for service 'tcp-8080': minPort (8080) cannot be greater than maxPort (8079)"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCreateService_Success_AddingTCPServiceWithMetafields(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":          "tcp-8080",
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080-8080"},
						"comment":       "test comment",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "tcp-8080",
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateService method
	err := client.CreateService(
		"test-adom",
		"tcp-8080",
		"tcp",
		8080,
		8080,
		"test comment",
		WithServiceMetafields(map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		}),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateService_Failure_ObjectAlreadyExists(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":          "tcp-8080",
						"protocol":      "TCP/UDP/UDP-Lite/SCTP",
						"tcp-portrange": []interface{}{"8080-8080"},
						"comment":       "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -2,
					"message": "Object already exists",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateService method
	err := client.CreateService("test-adom", "tcp-8080", "tcp", 8080, 8080, "test comment")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to create service: Object already exists"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestUpdateServiceWithMetafields_Success(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name": "tcp-443",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name": "tcp-443",
					"meta fields": map[string]interface{}{
						"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
						"some_randid": "abc",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the UpdateServiceWithMetafields method
	err := client.UpdateServiceWithMetafields(
		"test-adom",
		"tcp-443",
		map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateServiceWithMetafields_Failure_NoPermission(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name": "tcp-443",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the UpdateServiceWithMetafields method
	err := client.UpdateServiceWithMetafields(
		"test-adom",
		"tcp-443",
		map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		},
	)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to update service with metafields: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

//todo do more tests for TestUpdateSubnetAddressWithMetafields

// todo this test is not correct, in the response the name of the address is returned
// todo check the other tests for the same issue
func TestCreateSubnetAddress_Success_AddingSingleIP(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":    "test-subnet",
						"type":    "ipmask",
						"subnet":  "8.8.8.8/255.255.255.255",
						"comment": "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateSubnetAddress method
	err := client.CreateSubnetAddress("test-adom", "test-subnet", "8.8.8.8", "255.255.255.255", "test comment")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateSubnetAddress_Success_AddingSubnet(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":    "test-subnet",
						"type":    "ipmask",
						"subnet":  "8.8.8.0/255.255.255.0",
						"comment": "test comment",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateSubnetAddress method
	err := client.CreateSubnetAddress("test-adom", "test-subnet", "8.8.8.0", "255.255.255.0", "test comment")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateSubnetAddress_Success_AddingSubnetWithMetafields(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":    "test-subnet",
						"type":    "ipmask",
						"subnet":  "8.8.8.0/255.255.255.0",
						"comment": "test comment",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateSubnetAddress method
	err := client.CreateSubnetAddress(
		"test-adom",
		"test-subnet",
		"8.8.8.0",
		"255.255.255.0",
		"test comment",
		WithAddressMetafields(map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		}),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateSubnetAddress_Failure_ObjectAlreadyExists(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "add",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name":    "existing-subnet",
						"type":    "ipmask",
						"subnet":  "8.8.8.0/255.255.255.0",
						"comment": "already exists",
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Mocked response body indicating object already exists
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -2,
					"message": "Object already exists",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the CreateSubnetAddress method
	err := client.CreateSubnetAddress("test-adom", "existing-subnet", "8.8.8.0", "255.255.255.0", "already exists")

	// Validate the error
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to create subnet address: Object already exists"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestUpdateSubnetAddressWithMetafields_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name": "test-subnet",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	err := client.UpdateSubnetAddressWithMetafields(
		"test-adom",
		"test-subnet",
		map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateSubnetAddressWithMetafields_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "set",
		"params": []interface{}{
			map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"name": "test-subnet",
						"meta fields": map[string]interface{}{
							"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
							"some_randid": "abc",
						},
					},
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/sys/logout",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	err := client.UpdateSubnetAddressWithMetafields(
		"test-adom",
		"test-subnet",
		map[string]interface{}{
			"some_uuid":   "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
			"some_randid": "abc",
		},
	)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := "failed to update address with metafields: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

//todo do more tests for TestUpdateSubnetAddressWithMetafields

func TestGetAddressByName_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByName("test-adom", "testobj1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == nil {
		t.Fatalf("Expected address, got nil")
	}

	if address.Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", address.Name)
	}

	if !reflect.DeepEqual(address.Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(address.Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", address.Subnet)
	}
}

func TestGetAddressByName_Failure_ObjectDoesNotExist(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/inexistent",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/inexistent",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByName("test-adom", "inexistent")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "address 'inexistent' not found in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByMetafield_Success_StringMetafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "other-uuid",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == nil {
		t.Fatalf("Expected address, got nil")
	}

	if address.Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", address.Name)
	}

	if !reflect.DeepEqual(address.Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(address.Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", address.Subnet)
	}
}

func TestGetAddressByMetafield_Success_Float64Metafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": 123456789,
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": 1234,
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", float64(123456789))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == nil {
		t.Fatalf("Expected address, got nil")
	}

	if address.Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", address.Name)
	}

	if !reflect.DeepEqual(address.Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(address.Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", address.Subnet)
	}
}

func TestGetAddressByMetafield_Failure_MultipleMatches(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", "duplicate-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "multiple addresses found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByMetafield_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", "no-permission-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "error fetching addresses: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByMetafield_Failure_NoAddresses(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data":   []map[string]interface{}{},
				"status": map[string]interface{}{"code": 0, "message": "OK"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", "no-match-value")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "no address found with metafield 'some_uuid'='no-match-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByMetafield_Failure_NoMatchingAddress(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{"code": 0, "message": "OK"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByMetafield("test-adom", "some_uuid", "random-not-matched-uuid")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "no address found with metafield 'some_uuid'='random-not-matched-uuid' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByMetafield_Success_StringMetafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "other-uuid",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(address) != 1 {
		t.Fatalf("Expected 1 address, got %d", len(address))
	}

	if address[0].Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", address[0].Name)
	}

	if !reflect.DeepEqual(address[0].Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(address[0].Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", address[0].Subnet)
	}
}

func TestGetAddressesByMetafield_Success_TwoAddressesWithStringMetafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "other-uuid",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj3",
						"subnet":          []interface{}{"123.123.123.125", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f8",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"821c2x23-3c4f-4b1d-8e5f-4d8e175af6f8", "821c2x23-3c4f-4b1d-8e5f-4d8e175af6f9"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(addresses) != 2 {
		t.Fatalf("Expected 2 addresses, got %d", len(addresses))
	}

	if addresses[0].Name != "testobj3" || addresses[1].Name != "testobj1" {
		t.Errorf("Expected address names 'testobj1' or 'testobj3', got '%s' and '%s'", addresses[0].Name, addresses[1].Name)
	}

	if !reflect.DeepEqual(addresses[0].Subnet, []interface{}{"123.123.123.125", "255.255.255.255"}) &&
		!reflect.DeepEqual(addresses[0].Subnet, []string{"123.123.123.125", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.125 255.255.255.255], got %v", addresses[0].Subnet)
	}

	if !reflect.DeepEqual(addresses[1].Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(addresses[1].Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", addresses[1].Subnet)
	}
}

func TestGetAddressesByMetafield_Success_Float64Metafield(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": 123456789,
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": 1234,
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{float64(123456789)})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(addresses) != 1 {
		t.Fatalf("Expected 1 address, got %d", len(addresses))
	}

	if addresses[0].Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", addresses[0].Name)
	}

	if !reflect.DeepEqual(addresses[0].Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(addresses[0].Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", addresses[0].Subnet)
	}
}

func TestGetAddressesByMetafield_Failure_MultipleMatches(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "duplicate-value",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj3",
						"subnet":          []interface{}{"123.123.123.125", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "normal-value",
						},
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"duplicate-value", "normal-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if addresses != nil {
		t.Fatalf("Expected address to be nil, got %v", addresses)
	}

	expectedErrorMessage := "multiple addresses found with metafield 'some_uuid'='duplicate-value' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByMetafield_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"no-permission-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if addresses != nil {
		t.Fatalf("Expected address to be nil, got %v", addresses)
	}

	expectedErrorMessage := "error fetching addresses: No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByMetafield_Failure_NoAddresses(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data":   []map[string]interface{}{},
				"status": map[string]interface{}{"code": 0, "message": "OK"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"no-match-value"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if addresses != nil {
		t.Fatalf("Expected address to be nil, got %v", addresses)
	}

	expectedErrorMessage := "found 0 addresses corresponding the following 'some_uuid' metafield values: '[no-match-value]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByMetafield_Failure_NoMatchingAddress(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{"code": 0, "message": "OK"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"random-not-matched-uuid"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if addresses != nil {
		t.Fatalf("Expected address to be nil, got %v", addresses)
	}

	expectedErrorMessage := "found 0 addresses corresponding the following 'some_uuid' metafield values: '[random-not-matched-uuid]', but expected 1 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressesByMetafield_Failure_NotAllMatching(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet", "type"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": []map[string]interface{}{
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6291,
						"tagging":         nil,
						"name":            "testobj1",
						"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-1",
						},
					},
					{
						"dynamic_mapping": nil,
						"list":            nil,
						"oid":             6292,
						"tagging":         nil,
						"name":            "testobj2",
						"subnet":          []interface{}{"123.123.123.124", "255.255.255.255"},
						"type":            "ipmask",
						"meta fields": map[string]interface{}{
							"some_uuid": "random-uuid-2",
						},
					},
				},
				"status": map[string]interface{}{"code": 0, "message": "OK"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	addresses, err := client.GetAddressesByMetafield("test-adom", "some_uuid", []interface{}{"no-permission-value", "random-uuid-1"})
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if addresses != nil {
		t.Fatalf("Expected address to be nil, got %v", addresses)
	}

	expectedErrorMessage := "found 1 addresses corresponding the following 'some_uuid' metafield values: '[no-permission-value random-uuid-1]', but expected 2 in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByNameIPAndNetmask_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []interface{}{"123.123.123.123", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByNameIPAndNetmask("test-adom", "testobj1", "123.123.123.123", "255.255.255.255")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == nil {
		t.Fatalf("Expected address, got nil")
	}

	if address.Name != "testobj1" {
		t.Errorf("Expected address name 'testobj1', got '%s'", address.Name)
	}

	if !reflect.DeepEqual(address.Subnet, []interface{}{"123.123.123.123", "255.255.255.255"}) &&
		!reflect.DeepEqual(address.Subnet, []string{"123.123.123.123", "255.255.255.255"}) {
		t.Errorf("Expected subnet [123.123.123.123 255.255.255.255], got %v", address.Subnet)
	}
}

func TestGetAddressByNameIPAndNetmask_Failure_SameNameAndNetmaskButDifferentIP(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []interface{}{"111.111.111.111", "255.255.255.255"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByNameIPAndNetmask("test-adom", "testobj1", "123.123.123.123", "255.255.255.255")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "address 'testobj1' does not match IP '123.123.123.123' and netmask '255.255.255.255' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByNameIPAndNetmask_Failure_SameNameAndIPButDifferentNetmask(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"dynamic_mapping": nil,
					"list":            nil,
					"name":            "testobj1",
					"oid":             6291,
					"subnet":          []interface{}{"123.123.123.123", "255.255.255.0"},
					"tagging":         nil,
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByNameIPAndNetmask("test-adom", "testobj1", "123.123.123.123", "255.255.255.255")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "address 'testobj1' does not match IP '123.123.123.123' and netmask '255.255.255.255' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetAddressByNameIPAndNetmask_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "subnet"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/testobj1",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	address, err := client.GetAddressByNameIPAndNetmask("test-adom", "testobj1", "123.123.123.123", "255.255.255.255")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if address != nil {
		t.Fatalf("Expected address to be nil, got %v", address)
	}

	expectedErrorMessage := "error fetching address 'testobj1': No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndTCPProtocol_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{"8080-8080"},
					"udp-portrange": []interface{}{},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "tcp-8080", "tcp", 8080, 8080)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatalf("Expected service, got nil")
	}

	if service.Name != "tcp-8080" {
		t.Errorf("Expected service name 'tcp-8080', got '%s'", service.Name)
	}

	if service.Protocol != "TCP/UDP/UDP-Lite/SCTP" {
		t.Errorf("Expected service protocol 'TCP/UDP/UDP-Lite/SCTP', got '%s'", service.Protocol)
	}

	if service.TCPPortRange[0] != "8080-8080" {
		t.Errorf("Expected TCP port range '8080-8080', got '%v'", service.TCPPortRange)
	}
}

func TestGetServiceByNamePortRangeAndUDPProtocol_Success(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "udp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{},
					"udp-portrange": []interface{}{"8080-8080"},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "udp-8080", "udp", 8080, 8080)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatalf("Expected service, got nil")
	}

	if service.Name != "udp-8080" {
		t.Errorf("Expected service name 'udp-8080', got '%s'", service.Name)
	}

	if service.Protocol != "TCP/UDP/UDP-Lite/SCTP" {
		t.Errorf("Expected service protocol 'TCP/UDP/UDP-Lite/SCTP', got '%s'", service.Protocol)
	}

	if service.UDPPortRange[0] != "8080-8080" {
		t.Errorf("Expected UDP port range '8080-8080', got '%v'", service.UDPPortRange)
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_ObjectDoesNotExist(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "Object does not exist",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/inexistent",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "inexistent", "tcp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'inexistent' not found in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndPortRangeButDifferentProtocolRangeThanTCP(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{},
					"udp-portrange": []interface{}{"8080-8080"},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "tcp-8080", "tcp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'tcp-8080' does not match TCP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndPortRangeButDifferentProtocolThanTCP(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "ALL",
					"tcp-portrange": []interface{}{"8080-8080"},
					"udp-portrange": []interface{}{},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "tcp-8080", "tcp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'tcp-8080' does not match TCP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndTCPProtocolButDifferentPortRange(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "tcp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{"8080-8081"},
					"udp-portrange": []interface{}{},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "tcp-8080", "tcp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'tcp-8080' does not match TCP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndPortRangeButDifferentProtocolRangeThanUDP(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "udp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{"8080-8080"},
					"udp-portrange": []interface{}{},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "udp-8080", "udp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'udp-8080' does not match UDP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndPortRangeButDifferentProtocolThanUDP(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "udp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "ALL",
					"tcp-portrange": []interface{}{},
					"udp-portrange": []interface{}{"8080-8080"},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "udp-8080", "udp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'udp-8080' does not match UDP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_SameNameAndUDPProtocolButDifferentPortRange(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"name":          "udp-8080",
					"obj seq":       92,
					"oid":           5450,
					"protocol":      "TCP/UDP/UDP-Lite/SCTP",
					"tcp-portrange": []interface{}{},
					"udp-portrange": []interface{}{"8080-8081"},
					"meta fields": map[string]interface{}{
						"some_uuid": "1b1c2d3e-5678-90ab-cdef-1234567890ab",
					},
				},
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/udp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "udp-8080", "udp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "service 'udp-8080' does not match UDP port range '8080'-'8080' in ADOM 'test-adom'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestGetServiceByNamePortRangeAndProtocol_Failure_NoPermission(t *testing.T) {
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	expectedRequestBody := map[string]interface{}{
		"method": "get",
		"params": []interface{}{
			map[string]interface{}{
				"fields": []interface{}{"name", "protocol", "tcp-portrange", "udp-portrange"},
				"option": []interface{}{"get meta"},
				"url":    "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
		"verbose": float64(1),
	}

	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -11,
					"message": "No permission for the resource",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/tcp-8080",
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}
		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}
		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}
		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewAPIClient(mockServer.URL, "fake-key")

	service, err := client.GetServiceByNamePortAndProtocol("test-adom", "tcp-8080", "tcp", 8080, 8080)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if service != nil {
		t.Fatalf("Expected service to be nil, got %v", service)
	}

	expectedErrorMessage := "error fetching service 'tcp-8080': No permission for the resource"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeleteService_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/test-service",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/test-service",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteService method
	err := client.DeleteService("test-adom", "test-service")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteService_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/test-service",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/test-service",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DeleteService method
	err := client.DeleteService("test-adom", "test-service")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteService_Failure_InexistentServiceWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/nonexistent-service",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "entry not found",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/service/custom/nonexistent-service",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteService method
	err := client.DeleteService("test-adom", "nonexistent-service")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to delete service: entry not found"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDeleteAddress_Success_WithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/address/test-address",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/test-address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteAddress method
	err := client.DeleteAddress("test-adom", "test-address")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteAddress_Success_WithAPIClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer fake-key",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/address/test-address",
			},
		},
	}

	// Mocked response body
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    0,
					"message": "OK",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/test-address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewAPIClient(mockServer.URL, "fake-key")

	// Call the DeleteAddress method
	err := client.DeleteAddress("test-adom", "test-address")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteAddress_Failure_InexistentAddressWithUserClient(t *testing.T) {
	// Expected request headers
	expectedRequestHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	// Expected request body
	expectedRequestBody := map[string]interface{}{
		"method": "delete",
		"params": []interface{}{
			map[string]interface{}{
				"url": "/pm/config/adom/test-adom/obj/firewall/address/nonexistent-address",
			},
		},
		"session": "fake-session",
	}

	// Mocked response body indicating failure
	mockResponse := map[string]interface{}{
		"result": []map[string]interface{}{
			{
				"status": map[string]interface{}{
					"code":    -3,
					"message": "entry not found",
				},
				"url": "/pm/config/adom/test-adom/obj/firewall/address/nonexistent-address",
			},
		},
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the request URL
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("Expected URL path /jsonrpc, got %s", r.URL.Path)
		}

		// Validate the request headers
		for key, expectedValue := range expectedRequestHeaders {
			actualValue := r.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Validate the request body
		var actualRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualRequestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Errorf("Expected request body %v, got %v", expectedRequestBody, actualRequestBody)
		}

		// Write the mocked response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// Create a new FortiManager client
	client := NewUserClient(mockServer.URL, "fake-user", "fake-password")
	client.Session = "fake-session" // Set a fake session

	// Call the DeleteAddress method
	err := client.DeleteAddress("test-adom", "nonexistent-address")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Validate the error message
	expectedErrorMessage := "failed to delete address: entry not found"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}
