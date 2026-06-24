//go:build e2e

package fmgcli

import (
	"os"
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
	if adom == "" {
		adom = "root"
	}

	if host == "" || token == "" {
		t.Fatalf("set FMG_E2E_HOST and FMG_E2E_TOKEN to run e2e lock/unlock test")
	}

	client := NewAPIClient(host, token)

	locked := false
	defer func() {
		if locked {
			if err := client.Unlock(adom); err != nil {
				t.Errorf("deferred unlock cleanup failed for ADOM %q: %v", adom, err)
			}
		}
	}()

	if err := client.Lock(adom); err != nil {
		t.Fatalf("lock failed for ADOM %q: %v", adom, err)
	}
	locked = true

	if err := client.Unlock(adom); err != nil {
		t.Fatalf("unlock failed for ADOM %q: %v", adom, err)
	}
	locked = false
}
