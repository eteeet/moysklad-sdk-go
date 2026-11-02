package env

import (
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	t.Setenv("MOYSKLAD_TOKEN", "")
	if _, err := Client(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_TOKEN is required") {
		t.Fatalf("missing token: %v", err)
	}
	t.Setenv("MOYSKLAD_TOKEN", "token")
	t.Setenv("MOYSKLAD_BASE_URL", "ftp://example.test")
	if _, err := Client(); err == nil || !strings.Contains(err.Error(), "base URL") {
		t.Fatalf("bad base URL: %v", err)
	}
	t.Setenv("MOYSKLAD_BASE_URL", "")
	if c, err := Client(); err != nil || c == nil {
		t.Fatalf("c = %v, err = %v", c, err)
	}
}
