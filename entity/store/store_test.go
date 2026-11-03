package store_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/store"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID  = "b733c216-2783-11f1-0a83-1f2500000237"
	zoneID  = "8277a566-75f8-11ed-ac1a-000d00000000"
	slotID  = "9177a566-75f8-11ed-ac1a-000d00000000"
	apiBase = "https://api.moysklad.ru/api/remap/1.2/"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: apiBase + path, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	full = store.CreateRequest{
		Name: "Склад SDK", Description: model.Value(""), Code: model.Value("sdk-store"), ExternalCode: model.Value("SDK-EXT"),
		Archived: model.Value(false), Shared: model.Value(true), PathName: model.Value("Основной склад"), Address: model.Value(""),
		AddressFull: model.Value(model.AddressRequest{City: model.Value("Москва"), Street: model.Value("Ленинские горы"),
			Country: model.Value(ref("entity/country/40e6f69a-991c-4fbc-8be9-d0d906cad180", "country"))}),
		Attributes: model.Value([]model.Attribute{}),
		Parent:     model.Value(ref("entity/store/faf3ff5b-2e58-11e6-8a84-bae500000050", "store")),
		Owner:      model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")),
		Group:      model.Null[model.Reference](),
	}
	update = store.UpdateRequest{Name: model.Value("Склад3"), PathName: model.Value("Не Основной склад"), Address: model.Value(""),
		Archived: model.Value(true), Parent: model.Value(ref("entity/store/faf3ff5b-2e58-11e6-8a84-bae500000050", "store"))}
	slot = store.SlotCreateRequest{Name: "slot 1", Barcode: model.Value("0123456789")}
)

func TestStoreFixture(t *testing.T) {
	data := read(t, "upstream/store.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got store.Store
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "poi" || got.Code != "test" || got.Description != "ytr" || got.Archived || got.Shared ||
				got.PathName != "Основной склад" || got.Address != "" || got.AddressFull == nil || got.AddressFull.Comment != "try" ||
				got.Parent == nil || got.Parent.Meta.Type != "store" || got.Owner == nil || got.Group == nil || got.Updated.IsZero() {
				t.Fatalf("store = %+v", got)
			}
			if got.Zones == nil || got.Zones.Meta.Type != "storezone" || got.Zones.Meta.Size != 1 || len(got.Zones.Rows) != 1 ||
				got.Slots == nil || got.Slots.Meta.Href != apiBase+"entity/store/"+testID+"/slots" {
				t.Fatalf("collections = %+v %+v", got.Zones, got.Slots)
			}
			// Строки коллекций с expand — зоны и ячейки той же формы, что у ListZones/ListSlots.
			var z store.Zone
			var s store.Slot
			if err := codec.Unmarshal(got.Zones.Rows[0], &z); err != nil || z.ID != "b733ce21-2783-11f1-0a83-1f2500000239" || z.Name != "lkjl" {
				t.Fatalf("zone = %+v, err = %v", z, err)
			}
			if err := codec.Unmarshal(got.Slots.Rows[0], &s); err != nil || s.ID != "b733cab9-2783-11f1-0a83-1f2500000238" || s.Zone != nil {
				t.Fatalf("slot = %+v, err = %v", s, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/store/create-full.golden.json")},
			{"create minimal", store.CreateRequest{Name: "Склад"}, []byte(`{"name":"Склад"}`)},
			{"update", update, read(t, "custom/store/update.golden.json")},
			{"update null parent", store.UpdateRequest{Parent: model.Null[model.Reference]()}, []byte(`{"parent":null}`)},
			{"zone create", store.ZoneCreateRequest{Name: "Zona 2"}, []byte(`{"name":"Zona 2"}`)},
			{"zone update", store.ZoneUpdateRequest{Name: model.Value("Zone new"), ExternalCode: model.Value("")}, []byte(`{"name":"Zone new","externalCode":""}`)},
			{"slot create", slot, read(t, "custom/store/slot-create.golden.json")},
			{"slot update empty barcode", store.SlotUpdateRequest{Barcode: model.Value("")}, []byte(`{"barcode":""}`)},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[store.CreateRequest](), reflect.TypeFor[store.UpdateRequest](),
		reflect.TypeFor[store.ZoneCreateRequest](), reflect.TypeFor[store.ZoneUpdateRequest](),
		reflect.TypeFor[store.SlotCreateRequest](), reflect.TypeFor[store.SlotUpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Zones", "Slots", "Zone"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *store.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := store.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/` + testID + `"},"id":"` + testID + `","name":"poi"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*store.Store, error)
	}{
		{"get", "GET", "entity/store/" + testID, nil, func(ctx context.Context, c *transport.Client) (*store.Store, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/store", read(t, "custom/store/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*store.Store, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/store/" + testID, read(t, "custom/store/update.golden.json"), func(ctx context.Context, c *transport.Client) (*store.Store, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *store.Store) {
						if got.ID != testID || got.Name != "poi" {
							t.Fatalf("store = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Зоны и ячейки: Get — объект; Create и Update принимают объект и массив из одного объекта
// (форма примеров источника). Пустой массив, массив из двух, null и объект без id — ошибки.
func TestNestedHTTP(t *testing.T) {
	const zone = `{"meta":{"href":"` + apiBase + `entity/store/` + testID + `/zones/` + zoneID + `","type":"storezone"},"id":"` + zoneID + `","name":"Zone new","externalCode":"x"}`
	const slotJSON = `{"meta":{"href":"` + apiBase + `entity/store/` + testID + `/slots/` + slotID + `","type":"slot"},"id":"` + slotID + `","name":"slot 1","barcode":"0123456789","zone":{"meta":{"href":"` + apiBase + `entity/store/` + testID + `/zones/` + zoneID + `","type":"storezone"}}}`
	type call func(*store.Service, context.Context) (string, error)
	zoneCall := func(f func(*store.Service, context.Context) (*store.Zone, error)) call {
		return func(s *store.Service, ctx context.Context) (string, error) {
			z, err := f(s, ctx)
			if z == nil {
				return "", err
			}
			return z.ID + "/" + z.Name, err
		}
	}
	slotCall := func(f func(*store.Service, context.Context) (*store.Slot, error)) call {
		return func(s *store.Service, ctx context.Context) (string, error) {
			v, err := f(s, ctx)
			if v == nil {
				return "", err
			}
			if v.Zone == nil || v.Zone.Meta.Type != "storezone" {
				return "", nil
			}
			return v.ID + "/" + v.Barcode, err
		}
	}
	zoneOK, slotOK := zoneID+"/Zone new", slotID+"/0123456789"
	for _, tc := range []struct {
		name, method, path, body, response, want string
		invoke                                   call
	}{
		{"get zone", "GET", "entity/store/" + testID + "/zones/" + zoneID, "", zone, zoneOK,
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.GetZone(ctx, testID, zoneID)
			})},
		{"get zone array", "GET", "entity/store/" + testID + "/zones/" + zoneID, "", "[" + zone + "]", "",
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.GetZone(ctx, testID, zoneID)
			})},
		{"create zone array", "POST", "entity/store/" + testID + "/zones", `{"name":"Zone new"}`, "[" + zone + "]", zoneOK,
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.CreateZone(ctx, testID, store.ZoneCreateRequest{Name: "Zone new"})
			})},
		{"create zone object", "POST", "entity/store/" + testID + "/zones", `{"name":"Zone new"}`, zone, zoneOK,
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.CreateZone(ctx, testID, store.ZoneCreateRequest{Name: "Zone new"})
			})},
		{"create zone two", "POST", "entity/store/" + testID + "/zones", `{"name":"Zone new"}`, "[" + zone + "," + zone + "]", "",
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.CreateZone(ctx, testID, store.ZoneCreateRequest{Name: "Zone new"})
			})},
		{"create zone empty array", "POST", "entity/store/" + testID + "/zones", `{"name":"Zone new"}`, "[]", "",
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.CreateZone(ctx, testID, store.ZoneCreateRequest{Name: "Zone new"})
			})},
		{"update zone null", "PUT", "entity/store/" + testID + "/zones/" + zoneID, `{"name":"Zone new"}`, "null", "",
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.UpdateZone(ctx, testID, zoneID, store.ZoneUpdateRequest{Name: model.Value("Zone new")})
			})},
		{"update zone no id", "PUT", "entity/store/" + testID + "/zones/" + zoneID, `{"name":"Zone new"}`, `[{"name":"x"}]`, "",
			zoneCall(func(s *store.Service, ctx context.Context) (*store.Zone, error) {
				return s.UpdateZone(ctx, testID, zoneID, store.ZoneUpdateRequest{Name: model.Value("Zone new")})
			})},
		{"get slot", "GET", "entity/store/" + testID + "/slots/" + slotID, "", slotJSON, slotOK,
			slotCall(func(s *store.Service, ctx context.Context) (*store.Slot, error) {
				return s.GetSlot(ctx, testID, slotID)
			})},
		{"create slot", "POST", "entity/store/" + testID + "/slots", string(read(t, "custom/store/slot-create.golden.json")), "[" + slotJSON + "]", slotOK,
			slotCall(func(s *store.Service, ctx context.Context) (*store.Slot, error) {
				return s.CreateSlot(ctx, testID, slot)
			})},
		{"update slot", "PUT", "entity/store/" + testID + "/slots/" + slotID, `{"barcode":"0123456789"}`, "[" + slotJSON + "]", slotOK,
			slotCall(func(s *store.Service, ctx context.Context) (*store.Slot, error) {
				return s.UpdateSlot(ctx, testID, slotID, store.SlotUpdateRequest{Barcode: model.Value("0123456789")})
			})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(r.Body)
				if r.Method != tc.method || r.URL.Path != "/"+tc.path || r.URL.RawQuery != "" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if tc.body == "" && len(body) != 0 {
					t.Errorf("body = %s", body)
				} else if tc.body != "" {
					documenttest.JSONEqual(t, body, []byte(tc.body))
				}
				io.WriteString(w, tc.response)
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			got, err := tc.invoke(service(c), t.Context())
			if tc.want == "" {
				if err == nil || !strings.Contains(err.Error(), "store "+testID) || strings.Count(err.Error(), "moysklad:") != 1 {
					t.Fatalf("got = %q, err = %v", got, err)
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("got = %q, err = %v", got, err)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
		})
	}
}

// Вложенные списки — одна страница с limit/offset; Delete и массовые операции идут на
// коллекцию склада, meta элементов — с типом storezone или slot.
func TestNestedCollections(t *testing.T) {
	type request struct {
		method, uri string
		body        []byte
	}
	var got []request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, request{r.Method, r.URL.RequestURI(), body})
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"`+apiBase+r.URL.Path+`","size":1},"rows":[{"id":"`+zoneID+`","name":"z"}]}`)
		case r.Method == http.MethodDelete:
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		default:
			io.WriteString(w, `[{"meta":{"type":"slot"},"id":"`+slotID+`","name":"n"}]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	zones, err := s.ListZones(ctx, testID, query.PageOptions{Limit: 10, Offset: 20})
	if err != nil || len(zones.Rows) != 1 || zones.Rows[0].ID != zoneID {
		t.Fatalf("zones = %+v, err = %v", zones, err)
	}
	if slots, err := s.ListSlots(ctx, testID, query.PageOptions{}); err != nil || len(slots.Rows) != 1 {
		t.Fatalf("slots = %+v, err = %v", slots, err)
	}
	if err := s.DeleteZone(ctx, testID, zoneID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteSlot(ctx, testID, slotID); err != nil {
		t.Fatal(err)
	}
	if res, err := s.CreateManySlots(ctx, testID, []store.SlotCreateRequest{{Name: "n"}}); err != nil || len(res) != 1 || res[0].Value.ID != slotID {
		t.Fatalf("create many = %+v, err = %v", res, err)
	}
	if _, err := s.UpdateManyZones(ctx, testID, []batch.Update[store.ZoneUpdateRequest]{{ID: zoneID, Request: store.ZoneUpdateRequest{Name: model.Value("n")}}}); err != nil {
		t.Fatal(err)
	}
	if res, err := s.DeleteManySlots(ctx, testID, []string{slotID}); err != nil || len(res) != 1 || res[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", res, err)
	}
	if _, err := s.CreateManyZones(ctx, testID, []store.ZoneCreateRequest{{Name: "z"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpdateManySlots(ctx, testID, []batch.Update[store.SlotUpdateRequest]{{ID: slotID, Request: store.SlotUpdateRequest{Barcode: model.Value("1")}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.DeleteManyZones(ctx, testID, []string{zoneID}); err != nil {
		t.Fatal(err)
	}
	base := "/entity/store/" + testID
	zoneMeta := `{"href":"` + server.URL + base + `/zones/` + zoneID + `","type":"storezone","mediaType":"application/json"}`
	slotMeta := `{"href":"` + server.URL + base + `/slots/` + slotID + `","type":"slot","mediaType":"application/json"}`
	want := []request{
		{"GET", base + "/zones?limit=10&offset=20", nil},
		{"GET", base + "/slots?limit=1000&offset=0", nil},
		{"DELETE", base + "/zones/" + zoneID, nil},
		{"DELETE", base + "/slots/" + slotID, nil},
		{"POST", base + "/slots", []byte(`[{"name":"n"}]`)},
		{"POST", base + "/zones", []byte(`[{"meta":` + zoneMeta + `,"name":"n"}]`)},
		{"POST", base + "/slots/delete", []byte(`[{"meta":` + slotMeta + `}]`)},
		{"POST", base + "/zones", []byte(`[{"name":"z"}]`)},
		{"POST", base + "/slots", []byte(`[{"meta":` + slotMeta + `,"barcode":"1"}]`)},
		{"POST", base + "/zones/delete", []byte(`[{"meta":` + zoneMeta + `}]`)},
	}
	if len(got) != len(want) {
		t.Fatalf("requests = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].method != want[i].method || got[i].uri != want[i].uri {
			t.Errorf("request %d = %s %s, want %s %s", i, got[i].method, got[i].uri, want[i].method, want[i].uri)
		}
		if want[i].body == nil {
			if len(got[i].body) != 0 {
				t.Errorf("request %d body = %s", i, got[i].body)
			}
		} else {
			documenttest.JSONEqual(t, got[i].body, want[i].body)
		}
	}
}

// Ошибки входа вложенных методов — до HTTP, ResponseInfo сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, store.CreateRequest{}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, store.UpdateRequest{}, o...)
			return err
		},
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"list zones bad store": func(o ...transport.RequestOption) error {
			_, err := s.ListZones(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"list zones fields": func(o ...transport.RequestOption) error {
			_, err := s.ListZones(ctx, testID, query.PageOptions{Fields: "x"}, o...)
			return err
		},
		"list slots limit": func(o ...transport.RequestOption) error {
			_, err := s.ListSlots(ctx, testID, query.PageOptions{Limit: 1001}, o...)
			return err
		},
		"get zone bad id": func(o ...transport.RequestOption) error {
			_, err := s.GetZone(ctx, testID, zoneID+"/x", o...)
			return err
		},
		"get slot bad store": func(o ...transport.RequestOption) error { _, err := s.GetSlot(ctx, "../x", slotID, o...); return err },
		"create zone no name": func(o ...transport.RequestOption) error {
			_, err := s.CreateZone(ctx, testID, store.ZoneCreateRequest{}, o...)
			return err
		},
		"create slot bad store": func(o ...transport.RequestOption) error {
			_, err := s.CreateSlot(ctx, "bad", store.SlotCreateRequest{Name: "n"}, o...)
			return err
		},
		"update zone empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateZone(ctx, testID, zoneID, store.ZoneUpdateRequest{}, o...)
			return err
		},
		"update slot empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateSlot(ctx, testID, slotID, store.SlotUpdateRequest{}, o...)
			return err
		},
		"delete zone bad id":    func(o ...transport.RequestOption) error { return s.DeleteZone(ctx, testID, "bad", o...) },
		"delete slot bad store": func(o ...transport.RequestOption) error { return s.DeleteSlot(ctx, "bad", slotID, o...) },
		"create many zones bad store": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyZones(ctx, "bad", []store.ZoneCreateRequest{{Name: "z"}}, o...)
			return err
		},
		"create many slots empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManySlots(ctx, testID, nil, o...)
			return err
		},
		"update many zones bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyZones(ctx, testID, []batch.Update[store.ZoneUpdateRequest]{{ID: "bad", Request: store.ZoneUpdateRequest{Name: model.Value("n")}}}, o...)
			return err
		},
		"update many slots empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManySlots(ctx, testID, []batch.Update[store.SlotUpdateRequest]{{ID: slotID}}, o...)
			return err
		},
		"delete many zones bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyZones(ctx, testID, []string{"bad"}, o...)
			return err
		},
		"delete many slots bad store": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManySlots(ctx, "bad", []string{slotID}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "store", read(t, "upstream/store.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[store.Store] {
		return service(c).List
	}, func(t *testing.T, got store.Store) {
		if got.ID != testID || got.Name != "poi" {
			t.Fatalf("row = %+v", got)
		}
	})
}
