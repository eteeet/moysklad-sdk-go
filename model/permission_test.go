package model_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestPermissionsRoundTrip(t *testing.T) {
	// Флаг, области с неизвестным действием и scope, большое целое и вложенный объект — без потерь.
	const data = `{"importData":true,"viewAudit":false,"task":{"view":"ALL","done":"ASSIGNEE","archive":"NEW_SCOPE"},` +
		`"future":12345678901234567890,"nested":{"view":{"x":1}},"off":null}`
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got model.Permissions
			if err := codec.unmarshal([]byte(data), &got); err != nil {
				t.Fatal(err)
			}
			if f := got["importData"].Flag; f == nil || !*f {
				t.Errorf("importData = %+v", got["importData"])
			}
			if f := got["viewAudit"].Flag; f == nil || *f {
				t.Errorf("viewAudit = %+v", got["viewAudit"])
			}
			task := got["task"].Actions
			if task["view"] != model.ScopeAll || task["done"] != model.ScopeAssignee || task["archive"] != "NEW_SCOPE" {
				t.Errorf("task = %+v", task)
			}
			if string(got["future"].Raw) != "12345678901234567890" || string(got["nested"].Raw) != `{"view":{"x":1}}` {
				t.Errorf("raw = %s, %s", got["future"].Raw, got["nested"].Raw)
			}
			if p, ok := got["off"]; !ok || p.Flag != nil || p.Actions != nil || p.Raw != nil {
				t.Errorf("null = %+v, %v", p, ok)
			}
			out, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, out, data)
		})
	}
}

// Права индивидуальной роли из upstream-фикстуры: 31 флаг и 69 объектов областей.
func TestPermissionsFixture(t *testing.T) {
	data, err := os.ReadFile("../testdata/upstream/employee_role_individual.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		var role struct {
			Permissions model.Permissions `json:"permissions"`
		}
		if err := codec.unmarshal(data, &role); err != nil {
			t.Fatal(err)
		}
		flags, actions := 0, 0
		for name, p := range role.Permissions {
			switch {
			case p.Flag != nil:
				flags++
			case p.Actions != nil:
				actions++
			default:
				t.Errorf("%s: %s unexpected shape %+v", codec.name, name, p)
			}
		}
		if flags != 31 || actions != 69 {
			t.Fatalf("%s: flags = %d, actions = %d", codec.name, flags, actions)
		}
	}
}
