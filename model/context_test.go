package model_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestResponseContextUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantEmployee bool
		wantError    bool
	}{
		{"wrapped", `{"employee":{"meta":{"href":"https://example.com/employee","type":"employee"}}}`, true, false},
		{"flat", `{"employee":{"href":"https://example.com/employee","type":"employee"}}`, true, false},
		{"missing", `{}`, false, false},
		{"null", `{"employee":null}`, false, false},
		{"invalid", `{"employee":42}`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got model.ResponseContext
			err := json.Unmarshal([]byte(tt.body), &got)
			if tt.wantError {
				if err == nil {
					t.Fatal("ожидалась ошибка декодирования")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if (got.Employee != nil) != tt.wantEmployee {
				t.Fatalf("Employee = %v, ожидалось наличие: %v", got.Employee, tt.wantEmployee)
			}
			if got.Employee == nil {
				return
			}
			if got.Employee.Meta.Href != "https://example.com/employee" {
				t.Fatalf("неверный href: %q", got.Employee.Meta.Href)
			}
			if got.Employee.Meta.Type != "employee" {
				t.Fatalf("неверный type: %q", got.Employee.Meta.Type)
			}
		})
	}
}
