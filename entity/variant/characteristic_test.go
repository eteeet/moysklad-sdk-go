package variant_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const characteristicID = "d8f5dcfe-4951-4180-bdb5-504dfc9aec8d"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *variant.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := variant.NewService(c)
	return s
}

func TestCharacteristicFixture(t *testing.T) {
	data := read(t, "upstream/variantcharacteristic.json")
	for _, codec := range documenttest.Codecs {
		var got variant.CharacteristicDefinition
		if err := codec.Unmarshal(data, &got); err != nil || got.ID != characteristicID || got.Name != "Размер" ||
			got.Type != "string" || got.Required || got.Meta.Type != "attributemetadata" {
			t.Fatalf("%s: definition = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestCharacteristicHTTP(t *testing.T) {
	definition := string(read(t, "upstream/variantcharacteristic.json"))
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*variant.CharacteristicDefinition, error)
	}{
		{"get", "GET", "entity/variant/metadata/characteristics/" + characteristicID, nil, func(ctx context.Context, c *transport.Client) (*variant.CharacteristicDefinition, error) {
			return service(c).GetCharacteristic(ctx, characteristicID)
		}},
		{"create", "POST", "entity/variant/metadata/characteristics", []byte(`{"name":"Размер"}`), func(ctx context.Context, c *transport.Client) (*variant.CharacteristicDefinition, error) {
			return service(c).CreateCharacteristic(ctx, "Размер")
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(definition) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *variant.CharacteristicDefinition) {
						if got.ID != characteristicID {
							t.Fatalf("definition = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func list(ctx context.Context, c *transport.Client) (*[]variant.CharacteristicDefinition, error) {
	got, err := service(c).ListCharacteristics(ctx)
	if err != nil {
		return nil, err
	}
	return &got, nil
}

// Описания читаются из метаданных variant; пустой список допустим, ответ без meta
// или без characteristics — нет.
func TestListCharacteristics(t *testing.T) {
	metadata := read(t, "custom/variant/metadata.json")
	for _, codec := range documenttest.Codecs {
		var got struct {
			Characteristics []variant.CharacteristicDefinition `json:"characteristics"`
		}
		if err := codec.Unmarshal(metadata, &got); err != nil || len(got.Characteristics) != 2 || got.Characteristics[1].Name != "Цвет" {
			t.Fatalf("%s: metadata = %+v, err = %v", codec.Name, got, err)
		}
	}
	for _, tc := range documenttest.ResponseCases(string(metadata)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/variant/metadata", nil, list, func(got *[]variant.CharacteristicDefinition) {
				if len(*got) != 2 || (*got)[0].Name != "Размер" {
					t.Fatalf("definitions = %+v", *got)
				}
			})
		})
	}
	for _, tc := range []struct {
		name, body string
		n          int
		wantErr    bool
	}{
		{"none", `{"meta":{"href":"h"},"characteristics":[]}`, 0, false},
		{"unknown type and fields", `{"meta":{"href":"h"},"characteristics":[{"id":"x","name":"Цвет","type":"FUTURE","future":1}],"attributes":[]}`, 1, false},
		{"no meta", `{"characteristics":[]}`, 0, true},
		{"no characteristics", `{"meta":{"href":"h"}}`, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tc.body))
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			got, err := service(c).ListCharacteristics(t.Context())
			if tc.wantErr != (err != nil) || len(got) != tc.n || tc.n == 1 && got[0].Type != "FUTURE" {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
		})
	}
	documenttest.LocalErrors(t, list)
}

// Значение характеристики в модификации ссылается на описание; поля читаются обоими codecs.
func TestCharacteristicValue(t *testing.T) {
	data := []byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/variant/metadata/characteristics/` + characteristicID + `","type":"attributemetadata"},"id":"` + characteristicID + `","name":"Размер","value":"XL"}`)
	for _, codec := range documenttest.Codecs {
		var got variant.Characteristic
		if err := codec.Unmarshal(data, &got); err != nil || got.ID != characteristicID || got.Name != "Размер" || got.Value != "XL" || got.Meta.Href == "" {
			t.Fatalf("%s: value = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestCharacteristicValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	if _, err := s.CreateCharacteristic(t.Context(), ""); err == nil {
		t.Error("empty name accepted")
	}
	if _, err := s.GetCharacteristic(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if _, err := variant.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
