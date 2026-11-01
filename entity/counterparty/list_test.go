package counterparty_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/counterparty.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "counterparty", row, "", "archived", func(c *transport.Client) listtest.ListFunc[counterparty.Counterparty] {
		s, err := counterparty.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got counterparty.Counterparty) {
		if got.ID != "12a8b923-692c-11e6-8a84-bae500000056" || got.Name != "ООО Тестовый контрагент" {
			t.Fatalf("row = %+v", got)
		}
	})
}
