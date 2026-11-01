package organization_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/organization.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "organization", row, "", "archived", func(c *transport.Client) listtest.ListFunc[organization.Organization] {
		s, err := organization.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got organization.Organization) {
		if got.ID != "4b9d5bec-0575-11e6-9464-e4de00000008" || got.Name != "ОАО СветПром" {
			t.Fatalf("row = %+v", got)
		}
	})
}
