//go:build e2e

package fmgcli

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func cleanupUnlockIfLocked(t *testing.T, client *Client, adom string, locked *bool) {
	t.Helper()

	t.Cleanup(func() {
		if locked == nil || !*locked {
			return
		}

		if err := client.Unlock(adom); err != nil {
			t.Errorf("deferred unlock cleanup failed for ADOM %q: %v", adom, err)
		}
	})
}

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
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}

func TestE2E_Commit(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))

	if host == "" || token == "" || adom == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN and FMG_E2E_ADOM to run e2e commit test")
	}

	client := NewAPIClient(host, token)
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed for ADOM %q: %v", adom, err)
	}

	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
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

func TestE2E_GetPoliciesByID(t *testing.T) {
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

func TestE2E_GetPolicyByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	pkg := strings.TrimSpace(os.Getenv("FMG_E2E_PKG"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_KEY"))
	metafieldValue := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_VALUE"))
	policyIDStr := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_ID"))

	if host == "" || token == "" || adom == "" || pkg == "" || metafieldKey == "" || metafieldValue == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_PKG, FMG_E2E_POLICY_METAFIELD_KEY and FMG_E2E_POLICY_METAFIELD_VALUE to run e2e get-policy-by-metafield test")
	}

	client := NewAPIClient(host, token)

	policy, err := client.GetPolicyByMetafield(adom, pkg, metafieldKey, metafieldValue)
	if err != nil {
		t.Fatalf("get policy by metafield failed for %q=%q in ADOM %q pkg %q: %v", metafieldKey, metafieldValue, adom, pkg, err)
	}

	if policy == nil {
		t.Fatal("expected policy, got nil")
	}

	value, ok := policy.Metafields[metafieldKey]
	if !ok {
		t.Fatalf("expected metafield key %q in returned policy", metafieldKey)
	}

	if !reflect.DeepEqual(value, metafieldValue) {
		t.Fatalf("expected metafield %q to be %q, got %v", metafieldKey, metafieldValue, value)
	}

	if policyIDStr != "" {
		expectedPolicyID, err := strconv.Atoi(policyIDStr)
		if err != nil {
			t.Fatalf("invalid FMG_E2E_POLICY_ID %q: %v", policyIDStr, err)
		}

		if policy.PolicyID != expectedPolicyID {
			t.Fatalf("expected policy ID %d, got %d", expectedPolicyID, policy.PolicyID)
		}
	}
}

func TestE2E_GetPoliciesByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	pkg := strings.TrimSpace(os.Getenv("FMG_E2E_PKG"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_KEY"))
	metafieldValuesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_VALUES"))
	policyIDsCSV := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_IDS"))

	if host == "" || token == "" || adom == "" || pkg == "" || metafieldKey == "" || metafieldValuesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_PKG, FMG_E2E_POLICY_METAFIELD_KEY and FMG_E2E_POLICY_METAFIELD_VALUES to run e2e get-policies-by-metafield test")
	}

	parts := strings.Split(metafieldValuesCSV, ",")
	values := make([]interface{}, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		t.Fatalf("FMG_E2E_POLICY_METAFIELD_VALUES %q produced no values", metafieldValuesCSV)
	}

	client := NewAPIClient(host, token)

	policies, err := client.GetPoliciesByMetafield(adom, pkg, metafieldKey, values)
	if err != nil {
		t.Fatalf("get policies by metafield failed for key %q values %v in ADOM %q pkg %q: %v", metafieldKey, values, adom, pkg, err)
	}

	if len(policies) != len(values) {
		t.Fatalf("expected %d policies, got %d", len(values), len(policies))
	}

	for i, policy := range policies {
		value, ok := policy.Metafields[metafieldKey]
		if !ok {
			t.Fatalf("policy index %d missing metafield key %q", i, metafieldKey)
		}

		if !reflect.DeepEqual(value, values[i]) {
			t.Fatalf("policy index %d expected metafield %q=%v, got %v", i, metafieldKey, values[i], value)
		}
	}

	if policyIDsCSV != "" {
		idParts := strings.Split(policyIDsCSV, ",")
		expectedIDs := make([]int, 0, len(idParts))
		for _, part := range idParts {
			idStr := strings.TrimSpace(part)
			if idStr == "" {
				continue
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				t.Fatalf("invalid policy ID %q in FMG_E2E_POLICY_IDS: %v", idStr, err)
			}

			expectedIDs = append(expectedIDs, id)
		}

		if len(expectedIDs) != len(policies) {
			t.Fatalf("FMG_E2E_POLICY_IDS count %d does not match returned policies count %d", len(expectedIDs), len(policies))
		}

		for i, policy := range policies {
			if policy.PolicyID != expectedIDs[i] {
				t.Fatalf("policy index %d expected policy ID %d, got %d", i, expectedIDs[i], policy.PolicyID)
			}
		}
	}
}
