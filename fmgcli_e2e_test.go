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
		t.Skip("set FMG_E2E_HOST, FMG_E2E_USER and FMG_E2E_PASSWORD to run e2e login/logout test")
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
