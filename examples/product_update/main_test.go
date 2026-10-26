package main

import (
	"strings"
	"testing"
)

func TestRunRequiresExampleTarget(t *testing.T) {
	for _, flag := range []string{"", "0", "true"} {
		t.Setenv("MOYSKLAD_ALLOW_WRITE", flag)
		if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_ALLOW_WRITE=1") {
			t.Fatalf("flag %q: %v", flag, err)
		}
	}
	t.Setenv("MOYSKLAD_ALLOW_WRITE", "1")
	t.Setenv("MOYSKLAD_EXAMPLE_PRODUCT_ID", "")
	if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_EXAMPLE_PRODUCT_ID") {
		t.Fatalf("missing id: %v", err)
	}
	t.Setenv("MOYSKLAD_EXAMPLE_PRODUCT_ID", "12345678-1234-1234-1234-123456789012")
	for _, code := range []string{"", "production-product", "sdk-example-"} {
		t.Setenv("MOYSKLAD_EXAMPLE_EXTERNAL_CODE", code)
		if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_EXAMPLE_EXTERNAL_CODE") {
			t.Fatalf("code %q: %v", code, err)
		}
	}
	t.Setenv("MOYSKLAD_EXAMPLE_EXTERNAL_CODE", "sdk-example-test")
	t.Setenv("MOYSKLAD_TOKEN", "")
	if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_TOKEN") {
		t.Fatalf("missing token: %v", err)
	}
}
