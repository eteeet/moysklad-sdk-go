package country_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/country.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "country", row, "", "", func(c *transport.Client) listtest.ListFunc[country.Country] {
		s, err := country.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got country.Country) {
		if got.ID != "44444444-4444-4444-4444-444444444444" || got.Name != "Россия" {
			t.Fatalf("row = %+v", got)
		}
	})
}
