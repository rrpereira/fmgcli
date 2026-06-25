//go:build e2e

package fmgcli

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestE2E_LoginLogout(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	user := strings.TrimSpace(os.Getenv("FMG_E2E_USER"))
	password := os.Getenv("FMG_E2E_PASSWORD")

	if host == "" || user == "" || password == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_USER and FMG_E2E_PASSWORD to run e2e login/logout test")
	}

	client := NewUserClient(host, user, password)

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if client.Session == "" {
		t.Fatal("login succeeded but session is empty")
	}

	loggedIn := true
	defer func() {
		if !loggedIn {
			return
		}

		if err := client.Logout(); err != nil {
			t.Errorf("deferred logout cleanup failed: %v", err)
		}

		if client.Session != "" {
			t.Errorf("expected empty session after logout, got %q", client.Session)
		}
	}()

	if err := client.Logout(); err != nil {
		t.Fatalf("logout failed: %v", err)
	}
	loggedIn = false

	if client.Session != "" {
		t.Fatalf("expected empty session after logout, got %q", client.Session)
	}
}

func TestE2E_LockUnlock(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))

	if host == "" || token == "" || adom == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN and FMG_E2E_ADOM to run e2e lock/unlock test")
	}

	client := NewAPIClient(host, token)

	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}

	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
}

func TestE2E_GetAddressByName(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	addressName := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAME"))

	if host == "" || token == "" || adom == "" || addressName == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM and FMG_E2E_ADDRESS_NAME to run e2e get-address test")
	}

	client := NewAPIClient(host, token)

	address, err := client.GetAddressByName(adom, addressName)
	if err != nil {
		t.Fatalf("get address failed for %q in ADOM %q: %v", addressName, adom, err)
	}

	if address == nil {
		t.Fatalf("expected address %q, got nil", addressName)
	}

	if address.Name != addressName {
		t.Fatalf("expected address name %q, got %q", addressName, address.Name)
	}
}

func TestE2E_GetPolicyByID(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	pkg := strings.TrimSpace(os.Getenv("FMG_E2E_PKG"))
	policyIDStr := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_ID"))

	if host == "" || token == "" || adom == "" || pkg == "" || policyIDStr == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_PKG and FMG_E2E_POLICY_ID to run e2e get-policy test")
	}

	policyID, err := strconv.Atoi(policyIDStr)
	if err != nil {
		t.Fatalf("invalid FMG_E2E_POLICY_ID %q: %v", policyIDStr, err)
	}

	client := NewAPIClient(host, token)

	policies, err := client.GetPoliciesByID(pkg, adom, []int{policyID})
	if err != nil {
		t.Fatalf("get policy failed for policy ID %d in ADOM %q pkg %q: %v", policyID, adom, pkg, err)
	}

	if len(policies) != 1 {
		t.Fatalf("expected exactly 1 policy, got %d", len(policies))
	}

	if policies[0].PolicyID != policyID {
		t.Fatalf("expected policy ID %d, got %d", policyID, policies[0].PolicyID)
	}
}

func TestE2E_GetServiceByNamePortAndProtocol(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	serviceName := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_NAME"))
	protocol := strings.ToLower(strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_PROTOCOL")))
	minPortStr := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_MIN_PORT"))
	maxPortStr := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_MAX_PORT"))

	if host == "" || token == "" || adom == "" || serviceName == "" || protocol == "" || minPortStr == "" || maxPortStr == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_SERVICE_NAME, FMG_E2E_SERVICE_PROTOCOL, FMG_E2E_SERVICE_MIN_PORT and FMG_E2E_SERVICE_MAX_PORT to run e2e get-service test")
	}

	minPort, err := strconv.Atoi(minPortStr)
	if err != nil {
		t.Fatalf("invalid FMG_E2E_SERVICE_MIN_PORT %q: %v", minPortStr, err)
	}

	maxPort, err := strconv.Atoi(maxPortStr)
	if err != nil {
		t.Fatalf("invalid FMG_E2E_SERVICE_MAX_PORT %q: %v", maxPortStr, err)
	}

	client := NewAPIClient(host, token)

	service, err := client.GetServiceByNamePortAndProtocol(adom, serviceName, protocol, minPort, maxPort)
	if err != nil {
		t.Fatalf("get service failed for %q in ADOM %q: %v", serviceName, adom, err)
	}

	if service == nil {
		t.Fatalf("expected service %q, got nil", serviceName)
	}

	if service.Name != serviceName {
		t.Fatalf("expected service name %q, got %q", serviceName, service.Name)
	}
}
