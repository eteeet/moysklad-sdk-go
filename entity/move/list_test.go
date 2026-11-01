package move_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/move.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "move", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[move.Move] {
		s, err := move.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got move.Move) {
		if got.ID != "5ac52c9a-4f67-11e6-8a84-bae500000084" || got.Name != "705507" {
			t.Fatalf("row = %+v", got)
		}
	})
}
