package main

import (
	"strings"
	"testing"
)

func TestRunRequiresWritePermission(t *testing.T) {
	for _, flag := range []string{"", "0", "true"} {
		t.Setenv("MOYSKLAD_ALLOW_WRITE", flag)
		if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_ALLOW_WRITE=1") {
			t.Fatalf("flag %q: %v", flag, err)
		}
	}
	t.Setenv("MOYSKLAD_ALLOW_WRITE", "1")
	t.Setenv("MOYSKLAD_TOKEN", "")
	if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_TOKEN") {
		t.Fatalf("missing token: %v", err)
	}
}
