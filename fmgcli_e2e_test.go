//go:build e2e

package fmgcli

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

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

func parseCSVList(csv string) []string {
	parts := strings.Split(csv, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}

	return values
}

func parseCSVAnyList(csv string) []interface{} {
	values := parseCSVList(csv)
	anyValues := make([]interface{}, 0, len(values))
	for _, value := range values {
		anyValues = append(anyValues, value)
	}

	return anyValues
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

func TestE2E_GetAddressesByName(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	addressNamesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAMES"))

	if host == "" || token == "" || adom == "" || addressNamesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM and FMG_E2E_ADDRESS_NAMES to run e2e get-addresses-by-name test")
	}

	addressNames := parseCSVList(addressNamesCSV)

	if len(addressNames) == 0 {
		t.Fatalf("FMG_E2E_ADDRESS_NAMES %q produced no values", addressNamesCSV)
	}

	client := NewAPIClient(host, token)

	addresses, err := client.GetAddressesByName(addressNames, adom)
	if err != nil {
		t.Fatalf("get addresses by names failed for %v in ADOM %q: %v", addressNames, adom, err)
	}

	if len(addresses) != len(addressNames) {
		t.Fatalf("expected %d addresses, got %d", len(addressNames), len(addresses))
	}

	expected := make(map[string]int, len(addressNames))
	for _, name := range addressNames {
		expected[name]++
	}

	for i, address := range addresses {
		if address.Name == "" {
			t.Fatalf("address index %d has empty name", i)
		}

		if expected[address.Name] == 0 {
			t.Fatalf("address index %d returned unexpected name %q", i, address.Name)
		}

		expected[address.Name]--
	}

	for name, remaining := range expected {
		if remaining != 0 {
			t.Fatalf("expected address name %q to be returned %d more time(s)", name, remaining)
		}
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

func TestE2E_GetServiceByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_METAFIELD_KEY"))
	metafieldValue := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_METAFIELD_VALUE"))

	if host == "" || token == "" || adom == "" || metafieldKey == "" || metafieldValue == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_SERVICE_METAFIELD_KEY and FMG_E2E_SERVICE_METAFIELD_VALUE to run e2e get-service-by-metafield test")
	}

	client := NewAPIClient(host, token)

	service, err := client.GetServiceByMetafield(adom, metafieldKey, metafieldValue)
	if err != nil {
		t.Fatalf("get service by metafield failed for %q=%q in ADOM %q: %v", metafieldKey, metafieldValue, adom, err)
	}

	if service == nil {
		t.Fatal("expected service, got nil")
	}

	value, ok := service.Metafields[metafieldKey]
	if !ok {
		t.Fatalf("expected metafield key %q in returned service", metafieldKey)
	}

	if !reflect.DeepEqual(value, metafieldValue) {
		t.Fatalf("expected metafield %q to be %q, got %v", metafieldKey, metafieldValue, value)
	}
}

func TestE2E_GetServicesByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_METAFIELD_KEY"))
	metafieldValuesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_METAFIELD_VALUES"))

	if host == "" || token == "" || adom == "" || metafieldKey == "" || metafieldValuesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_SERVICE_METAFIELD_KEY and FMG_E2E_SERVICE_METAFIELD_VALUES to run e2e get-services-by-metafield test")
	}

	values := parseCSVAnyList(metafieldValuesCSV)

	if len(values) == 0 {
		t.Fatalf("FMG_E2E_SERVICE_METAFIELD_VALUES %q produced no values", metafieldValuesCSV)
	}

	client := NewAPIClient(host, token)

	services, err := client.GetServicesByMetafield(adom, metafieldKey, values)
	if err != nil {
		t.Fatalf("get services by metafield failed for key %q values %v in ADOM %q: %v", metafieldKey, values, adom, err)
	}

	if len(services) != len(values) {
		t.Fatalf("expected %d services, got %d", len(values), len(services))
	}

	for i, service := range services {
		value, ok := service.Metafields[metafieldKey]
		if !ok {
			t.Fatalf("service index %d missing metafield key %q", i, metafieldKey)
		}

		if !reflect.DeepEqual(value, values[i]) {
			t.Fatalf("service index %d expected metafield %q=%v, got %v", i, metafieldKey, values[i], value)
		}
	}
}

func TestE2E_GetServicesByName(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	serviceNamesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_NAMES"))

	if host == "" || token == "" || adom == "" || serviceNamesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM and FMG_E2E_SERVICE_NAMES to run e2e get-services-by-name test")
	}

	serviceNames := parseCSVList(serviceNamesCSV)

	if len(serviceNames) == 0 {
		t.Fatalf("FMG_E2E_SERVICE_NAMES %q produced no values", serviceNamesCSV)
	}

	client := NewAPIClient(host, token)

	services, err := client.GetServicesByName(serviceNames, adom)
	if err != nil {
		t.Fatalf("get services by names failed for %v in ADOM %q: %v", serviceNames, adom, err)
	}

	if len(services) != len(serviceNames) {
		t.Fatalf("expected %d services, got %d", len(serviceNames), len(services))
	}

	expected := make(map[string]int, len(serviceNames))
	for _, name := range serviceNames {
		expected[name]++
	}

	for i, service := range services {
		if service.Name == "" {
			t.Fatalf("service index %d has empty name", i)
		}

		if expected[service.Name] == 0 {
			t.Fatalf("service index %d returned unexpected name %q", i, service.Name)
		}

		expected[service.Name]--
	}

	for name, remaining := range expected {
		if remaining != 0 {
			t.Fatalf("expected service name %q to be returned %d more time(s)", name, remaining)
		}
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
}

func TestE2E_GetPoliciesByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	pkg := strings.TrimSpace(os.Getenv("FMG_E2E_PKG"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_KEY"))
	metafieldValuesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_POLICY_METAFIELD_VALUES"))

	if host == "" || token == "" || adom == "" || pkg == "" || metafieldKey == "" || metafieldValuesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_PKG, FMG_E2E_POLICY_METAFIELD_KEY and FMG_E2E_POLICY_METAFIELD_VALUES to run e2e get-policies-by-metafield test")
	}

	values := parseCSVAnyList(metafieldValuesCSV)

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
}

func TestE2E_GetAddressByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_METAFIELD_KEY"))
	metafieldValue := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_METAFIELD_VALUE"))

	if host == "" || token == "" || adom == "" || metafieldKey == "" || metafieldValue == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_ADDRESS_METAFIELD_KEY and FMG_E2E_ADDRESS_METAFIELD_VALUE to run e2e get-address-by-metafield test")
	}

	client := NewAPIClient(host, token)

	address, err := client.GetAddressByMetafield(adom, metafieldKey, metafieldValue)
	if err != nil {
		t.Fatalf("get address by metafield failed for %q=%q in ADOM %q: %v", metafieldKey, metafieldValue, adom, err)
	}

	if address == nil {
		t.Fatal("expected address, got nil")
	}

	value, ok := address.Metafields[metafieldKey]
	if !ok {
		t.Fatalf("expected metafield key %q in returned address", metafieldKey)
	}

	if !reflect.DeepEqual(value, metafieldValue) {
		t.Fatalf("expected metafield %q to be %q, got %v", metafieldKey, metafieldValue, value)
	}
}

func TestE2E_GetAddressesByMetafield(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_METAFIELD_KEY"))
	metafieldValuesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_METAFIELD_VALUES"))

	if host == "" || token == "" || adom == "" || metafieldKey == "" || metafieldValuesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_ADDRESS_METAFIELD_KEY and FMG_E2E_ADDRESS_METAFIELD_VALUES to run e2e get-addresses-by-metafield test")
	}

	values := parseCSVAnyList(metafieldValuesCSV)

	if len(values) == 0 {
		t.Fatalf("FMG_E2E_ADDRESS_METAFIELD_VALUES %q produced no values", metafieldValuesCSV)
	}

	client := NewAPIClient(host, token)

	addresses, err := client.GetAddressesByMetafield(adom, metafieldKey, values)
	if err != nil {
		t.Fatalf("get addresses by metafield failed for key %q values %v in ADOM %q: %v", metafieldKey, values, adom, err)
	}

	if len(addresses) != len(values) {
		t.Fatalf("expected %d addresses, got %d", len(values), len(addresses))
	}

	for i, address := range addresses {
		value, ok := address.Metafields[metafieldKey]
		if !ok {
			t.Fatalf("address index %d missing metafield key %q", i, metafieldKey)
		}

		if !reflect.DeepEqual(value, values[i]) {
			t.Fatalf("address index %d expected metafield %q=%v, got %v", i, metafieldKey, values[i], value)
		}
	}
}

func TestE2E_GetAddressByNameIPAndNetmask(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	addressName := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAME_IP_NETMASK_NAME"))
	ip := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAME_IP_NETMASK_IP"))
	netmask := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAME_IP_NETMASK_NETMASK"))

	if host == "" || token == "" || adom == "" || addressName == "" || ip == "" || netmask == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_ADDRESS_NAME_IP_NETMASK_NAME, FMG_E2E_ADDRESS_NAME_IP_NETMASK_IP and FMG_E2E_ADDRESS_NAME_IP_NETMASK_NETMASK to run e2e get-address-by-name-ip-netmask test")
	}

	client := NewAPIClient(host, token)

	address, err := client.GetAddressByNameIPAndNetmask(adom, addressName, ip, netmask)
	if err != nil {
		t.Fatalf("get address by name, IP and netmask failed for %q (IP: %s, Netmask: %s) in ADOM %q: %v", addressName, ip, netmask, adom, err)
	}

	if address == nil {
		t.Fatalf("expected address %q, got nil", addressName)
	}

	if address.Name != addressName {
		t.Fatalf("expected address name %q, got %q", addressName, address.Name)
	}

	if len(address.Subnet) < 2 {
		t.Fatalf("expected address to have subnet with at least 2 elements, got %d", len(address.Subnet))
	}

	if address.Subnet[0] != ip {
		t.Fatalf("expected IP %q, got %q", ip, address.Subnet[0])
	}

	if address.Subnet[1] != netmask {
		t.Fatalf("expected netmask %q, got %q", netmask, address.Subnet[1])
	}
}

func TestE2E_CreateDisableDeletePolicy(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	pkg := strings.TrimSpace(os.Getenv("FMG_E2E_PKG"))
	addressNamesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAMES"))
	serviceNamesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_NAMES"))

	if host == "" || token == "" || adom == "" || pkg == "" || addressNamesCSV == "" || serviceNamesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_PKG, FMG_E2E_ADDRESS_NAMES and FMG_E2E_SERVICE_NAMES to run e2e create-disable-delete-policy test")
	}

	srcs := parseCSVList(addressNamesCSV)
	dsts := parseCSVList(addressNamesCSV)
	services := parseCSVList(serviceNamesCSV)

	if len(srcs) == 0 {
		t.Fatalf("FMG_E2E_ADDRESS_NAMES %q produced no values", addressNamesCSV)
	}

	if len(services) == 0 {
		t.Fatalf("FMG_E2E_SERVICE_NAMES %q produced no values", serviceNamesCSV)
	}

	client := NewAPIClient(host, token)
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	// Lock the ADOM for modifications
	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	// Create a policy
	policyID, err := client.CreatePolicy(pkg, adom, "Test policy for E2E test", srcs, dsts, services)
	if err != nil {
		t.Fatalf("create policy failed for ADOM %q pkg %q: %v", adom, pkg, err)
	}

	if policyID == 0 {
		t.Fatalf("expected non-zero policy ID, got %d", policyID)
	}

	t.Logf("Successfully created policy with ID: %d", policyID)

	// Commit the policy creation
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after creating policy for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed policy creation")

	// Verify the policy was created and is enabled
	policies, err := client.GetPoliciesByID(pkg, adom, []int{policyID})
	if err != nil {
		t.Fatalf("get policy failed for policy ID %d in ADOM %q pkg %q: %v", policyID, adom, pkg, err)
	}

	if len(policies) != 1 {
		t.Fatalf("expected exactly 1 policy, got %d", len(policies))
	}

	if policies[0].Status != "enable" {
		t.Fatalf("expected policy status 'enable', got %q", policies[0].Status)
	}

	t.Logf("Verified policy is enabled")

	// Disable the policy
	if err := client.DisablePolicy(adom, pkg, policyID); err != nil {
		t.Fatalf("disable policy failed for policy ID %d in ADOM %q pkg %q: %v", policyID, adom, pkg, err)
	}

	t.Logf("Successfully disabled policy with ID: %d", policyID)

	// Commit the policy disabling
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after disabling policy for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed policy disabling")

	// Verify the policy was disabled
	policies, err = client.GetPoliciesByID(pkg, adom, []int{policyID})
	if err != nil {
		t.Fatalf("get policy failed for policy ID %d in ADOM %q pkg %q: %v", policyID, adom, pkg, err)
	}

	if len(policies) != 1 {
		t.Fatalf("expected exactly 1 policy, got %d", len(policies))
	}

	if policies[0].Status != "disable" {
		t.Fatalf("expected policy status 'disable', got %q", policies[0].Status)
	}

	t.Logf("Verified policy is disabled")

	// Delete the policy
	if err := client.DeletePolicy(adom, pkg, policyID); err != nil {
		t.Fatalf("delete policy failed for policy ID %d in ADOM %q pkg %q: %v", policyID, adom, pkg, err)
	}

	t.Logf("Successfully deleted policy with ID: %d", policyID)

	// Commit the policy deletion
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after deleting policy for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed policy deletion")

	// Verify the policy was deleted (should return an error since policy no longer exists)
	_, err = client.GetPoliciesByID(pkg, adom, []int{policyID})
	if err == nil {
		t.Fatalf("expected error when fetching deleted policy, but got none")
	}

	t.Logf("Verified policy was deleted (GetPoliciesByID returned expected error)")

	// Unlock the ADOM
	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}

func TestE2E_CreateUpdateDeleteService(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_SERVICE_METAFIELD_KEY"))
	protocol := "udp"
	minPort := 59999
	maxPort := 60002

	if host == "" || token == "" || adom == "" || metafieldKey == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM and FMG_E2E_SERVICE_METAFIELD_KEY to run e2e create-update-delete-service test")
	}

	serviceName := "e2e_test_service_" + strconv.FormatInt(time.Now().Unix(), 10)
	metafieldValue := "e2e-service-update-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	client := NewAPIClient(host, token)
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	// Lock the ADOM for modifications
	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	// Create a service
	err := client.CreateService(adom, serviceName, protocol, minPort, maxPort, "E2E test service")
	if err != nil {
		t.Fatalf("create service failed for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully created service with name: %s", serviceName)

	// Commit the service creation
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after creating service for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed service creation")

	// Verify the service was created
	service, err := client.GetServiceByNamePortAndProtocol(adom, serviceName, protocol, minPort, maxPort)
	if err != nil {
		t.Fatalf("get service failed for service name %q in ADOM %q: %v", serviceName, adom, err)
	}

	if service == nil {
		t.Fatalf("expected service %q, got nil", serviceName)
	}

	if service.Name != serviceName {
		t.Fatalf("expected service name %q, got %q", serviceName, service.Name)
	}

	t.Logf("Verified service was created")

	// Update the service
	err = client.UpdateServiceWithMetafields(adom, serviceName, map[string]interface{}{
		metafieldKey: metafieldValue,
	})
	if err != nil {
		t.Fatalf("update service failed for service name %q in ADOM %q: %v", serviceName, adom, err)
	}

	t.Logf("Successfully updated service with metafield %s=%s", metafieldKey, metafieldValue)

	// Commit the service update
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after updating service for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed service update")

	// Verify the service was updated
	updatedService, err := client.GetServiceByMetafield(adom, metafieldKey, metafieldValue)
	if err != nil {
		t.Fatalf("get updated service by metafield failed for %q=%q in ADOM %q: %v", metafieldKey, metafieldValue, adom, err)
	}

	if updatedService == nil {
		t.Fatalf("expected updated service %q, got nil", serviceName)
	}

	if updatedService.Name != serviceName {
		t.Fatalf("expected updated service name %q, got %q", serviceName, updatedService.Name)
	}

	t.Logf("Verified service was updated")

	// Delete the service
	if err := client.DeleteService(adom, serviceName); err != nil {
		t.Fatalf("delete service failed for service name %q in ADOM %q: %v", serviceName, adom, err)
	}

	t.Logf("Successfully deleted service with name: %s", serviceName)

	// Commit the service deletion
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after deleting service for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed service deletion")

	// Verify the service was deleted (should return an error since service no longer exists)
	_, err = client.GetServiceByNamePortAndProtocol(adom, serviceName, protocol, minPort, maxPort)
	if err == nil {
		t.Fatalf("expected error when fetching deleted service, but got none")
	}

	t.Logf("Verified service was deleted (GetServiceByNamePortAndProtocol returned expected error)")

	// Unlock the ADOM
	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}

func TestE2E_CreateUpdateDeleteAddress(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_METAFIELD_KEY"))
	ip := "8.8.8.8"
	netmask := "255.255.255.255"

	if host == "" || token == "" || adom == "" || metafieldKey == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM and FMG_E2E_ADDRESS_METAFIELD_KEY to run e2e create-update-delete-address test")
	}

	addressName := "e2e_test_address_" + strconv.FormatInt(time.Now().Unix(), 10)
	metafieldValue := "e2e-address-update-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	client := NewAPIClient(host, token)
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	// Lock the ADOM for modifications
	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	// Create an address
	err := client.CreateSubnetAddress(adom, addressName, ip, netmask, "E2E test address")
	if err != nil {
		t.Fatalf("create address failed for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully created address with name: %s", addressName)

	// Commit the address creation
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after creating address for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed address creation")

	// Verify the address was created
	address, err := client.GetAddressByNameIPAndNetmask(adom, addressName, ip, netmask)
	if err != nil {
		t.Fatalf("get address failed for address name %q in ADOM %q: %v", addressName, adom, err)
	}

	if address == nil {
		t.Fatalf("expected address %q, got nil", addressName)
	}

	if address.Name != addressName {
		t.Fatalf("expected address name %q, got %q", addressName, address.Name)
	}

	t.Logf("Verified address was created")

	// Update the address
	err = client.UpdateSubnetAddressWithMetafields(adom, addressName, map[string]interface{}{
		metafieldKey: metafieldValue,
	})
	if err != nil {
		t.Fatalf("update address failed for address name %q in ADOM %q: %v", addressName, adom, err)
	}

	t.Logf("Successfully updated address with metafield %s=%s", metafieldKey, metafieldValue)

	// Commit the address update
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after updating address for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed address update")

	// Verify the address was updated
	updatedAddress, err := client.GetAddressByMetafield(adom, metafieldKey, metafieldValue)
	if err != nil {
		t.Fatalf("get updated address by metafield failed for %q=%q in ADOM %q: %v", metafieldKey, metafieldValue, adom, err)
	}

	if updatedAddress == nil {
		t.Fatalf("expected updated address %q, got nil", addressName)
	}

	if updatedAddress.Name != addressName {
		t.Fatalf("expected updated address name %q, got %q", addressName, updatedAddress.Name)
	}

	t.Logf("Verified address was updated")

	// Delete the address
	if err := client.DeleteAddress(adom, addressName); err != nil {
		t.Fatalf("delete address failed for address name %q in ADOM %q: %v", addressName, adom, err)
	}

	t.Logf("Successfully deleted address with name: %s", addressName)

	// Commit the address deletion
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after deleting address for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed address deletion")

	// Verify the address was deleted (should return an error since address no longer exists)
	_, err = client.GetAddressByNameIPAndNetmask(adom, addressName, ip, netmask)
	if err == nil {
		t.Fatalf("expected error when fetching deleted address, but got none")
	}

	t.Logf("Verified address was deleted (GetAddressByNameIPAndNetmask returned expected error)")

	// Unlock the ADOM
	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}

func TestE2E_CreateUpdateDeleteGroup(t *testing.T) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("FMG_E2E_HOST"))
	token := strings.TrimSpace(os.Getenv("FMG_E2E_TOKEN"))
	adom := strings.TrimSpace(os.Getenv("FMG_E2E_ADOM"))
	metafieldKey := strings.TrimSpace(os.Getenv("FMG_E2E_GROUP_METAFIELD_KEY"))
	addressNamesCSV := strings.TrimSpace(os.Getenv("FMG_E2E_ADDRESS_NAMES"))

	if host == "" || token == "" || adom == "" || metafieldKey == "" || addressNamesCSV == "" {
		t.Fatalf("set FMG_E2E_HOST, FMG_E2E_TOKEN, FMG_E2E_ADOM, FMG_E2E_GROUP_METAFIELD_KEY and FMG_E2E_ADDRESS_NAMES to run e2e create-update-delete-group test")
	}

	groupMembers := parseCSVList(addressNamesCSV)
	if len(groupMembers) == 0 {
		t.Fatalf("FMG_E2E_ADDRESS_NAMES %q produced no values", addressNamesCSV)
	}

	groupName := "e2e_test_group_" + strconv.FormatInt(time.Now().Unix(), 10)
	metafieldValue := "e2e-group-update-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	client := NewAPIClient(host, token)
	locked := false
	cleanupUnlockIfLocked(t, client, adom, &locked)

	// Lock the ADOM for modifications
	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	// Create a group
	err := client.CreateGroup(adom, groupName, groupMembers, "E2E test group")
	if err != nil {
		t.Fatalf("create group failed for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully created group with name: %s", groupName)

	// Commit the group creation
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after creating group for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed group creation")

	// Update the group
	err = client.UpdateGroupWithMetafields(adom, groupName, map[string]interface{}{
		metafieldKey: metafieldValue,
	})
	if err != nil {
		t.Fatalf("update group failed for group name %q in ADOM %q: %v", groupName, adom, err)
	}

	t.Logf("Successfully updated group with metafield %s=%s", metafieldKey, metafieldValue)

	// Commit the group update
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after updating group for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed group update")

	// Delete the group
	if err := client.DeleteGroup(adom, groupName); err != nil {
		t.Fatalf("delete group failed for group name %q in ADOM %q: %v", groupName, adom, err)
	}

	t.Logf("Successfully deleted group with name: %s", groupName)

	// Commit the group deletion
	if err := client.Commit(adom); err != nil {
		t.Fatalf("commit failed after deleting group for ADOM %q: %v", adom, err)
	}

	t.Logf("Successfully committed group deletion")

	// Verify the group was deleted by expecting a second delete to fail.
	err = client.DeleteGroup(adom, groupName)
	if err == nil {
		t.Fatalf("expected error when deleting already deleted group %q, but got none", groupName)
	}

	t.Logf("Verified group was deleted (second DeleteGroup returned expected error)")

	// Unlock the ADOM
	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}
