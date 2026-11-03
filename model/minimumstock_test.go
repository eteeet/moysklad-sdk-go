package model_test

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Ключ коллекции: storebalances (таблица) и storeBalances (фикстура); оба читаются обоими кодеками.
func TestMinimumStockKeys(t *testing.T) {
	for _, data := range []string{
		`{"type":"WAREHOUSE_VARIED","storebalances":{"meta":{"href":"https://api.moysklad.ru/x/storebalances","type":"minimumstock","size":2}}}`,
		`{"type":"WAREHOUSE_VARIED","storeBalances":{"meta":{"href":"https://api.moysklad.ru/x/storebalances","type":"minimumstock","size":2}}}`,
	} {
		for name, unmarshal := range map[string]func([]byte, any) error{"v1": jsonv1.Unmarshal, "v2": func(b []byte, v any) error { return jsonv2.Unmarshal(b, v) }} {
			var got model.MinimumStock
			if err := unmarshal([]byte(data), &got); err != nil || got.Type != model.MinimumStockWarehouseVaried ||
				got.StoreBalances == nil || got.StoreBalances.Meta.Size != 2 || got.StoreBalances.Meta.Href != "https://api.moysklad.ru/x/storebalances" {
				t.Fatalf("%s %s: got = %+v, err = %v", name, data, got, err)
			}
		}
	}
	var sum model.MinimumStock
	if err := jsonv2.Unmarshal([]byte(`{"type":"ALL_WAREHOUSE_SUM","quantity":10,"future":1}`), &sum); err != nil || sum.Quantity != 10 || sum.StoreBalances != nil {
		t.Fatalf("sum = %+v, err = %v", sum, err)
	}
	if err := jsonv2.Unmarshal([]byte(`{"type":1}`), &sum); err == nil {
		t.Fatal("number type accepted")
	}
}

// Пример изменения источника: WAREHOUSE_VARIED со строками storebalances.
func TestMinimumStockRequestJSON(t *testing.T) {
	store := model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/store/da1c5548-0ecf-11f0-ac15-001200000125", Type: "store"}}
	req := model.MinimumStockRequest{Type: model.MinimumStockWarehouseVaried, StoreBalances: model.Value([]model.StoreBalanceRequest{
		{Meta: model.Value(model.Meta{Href: "https://api.moysklad.ru/x/storebalances/1", Type: "minimumstock"}), Quantity: model.Value(2.5), Store: model.Value(store)},
		{Store: model.Value(store), Quantity: model.Value(0.0)},
	})}
	want := `{"type":"WAREHOUSE_VARIED","storebalances":[{"meta":{"href":"https://api.moysklad.ru/x/storebalances/1","type":"minimumstock"},"store":{"meta":{"href":"` + store.Meta.Href + `","type":"store"}},"quantity":2.5},` +
		`{"store":{"meta":{"href":"` + store.Meta.Href + `","type":"store"}},"quantity":0}]}`
	for name, marshal := range map[string]func(any) ([]byte, error){"v1": jsonv1.Marshal, "v2": func(v any) ([]byte, error) { return jsonv2.Marshal(v) }} {
		got, err := marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		var a, b any
		jsonv2.Unmarshal(got, &a)
		jsonv2.Unmarshal([]byte(want), &b)
		if ja, _ := jsonv2.Marshal(a, jsonv2.Deterministic(true)); string(ja) != func() string { jb, _ := jsonv2.Marshal(b, jsonv2.Deterministic(true)); return string(jb) }() {
			t.Fatalf("%s: %s", name, got)
		}
	}
}
