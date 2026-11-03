package retailshift_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retailshift"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "12121212-1212-1212-1212-121212121212"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *retailshift.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retailshift.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	retailStore  = ref("entity/retailstore/22222222-2222-2222-2222-222222222222", "retailstore")
	openShift    = retailshift.CreateRequest{
		Organization: organization, RetailStore: retailStore,
		Name:          model.Value("0001"),
		Moment:        model.Value(documenttest.DateTime("2026-11-02 09:00:00.000")),
		Store:         model.Value(ref("entity/store/33333333-3333-3333-3333-333333333333", "store")),
		BankComission: model.Value(0.0), BankPercent: model.Value(1.5),
		Cheque: model.Value(retailshift.ChequeRequest{Start: model.Value(retailshift.ChequeStartRequest{FNNumber: model.Value("1234"),
			KKTRegNumber: model.Value("5678"), ShiftNumber: model.Value("17"), Time: model.Value(documenttest.DateTime("2026-11-02 09:00:00.000"))})}),
	}
	closeShift = retailshift.UpdateRequest{
		CloseDate: model.Value(documenttest.DateTime("2026-11-02 21:00:00.000")),
		QRAcquire: model.Null[model.Reference](),
		Cheque: model.Value(retailshift.ChequeRequest{End: model.Value(retailshift.ChequeEndRequest{
			ChequeStartRequest: retailshift.ChequeStartRequest{FNNumber: model.Value("1234"), ShiftNumber: model.Value("17"),
				Time: model.Value(documenttest.DateTime("2026-11-02 21:00:00.000"))},
			ChequesTotal: model.Value(0), FiscalDocsTotal: model.Value(44)})}),
	}
)

func TestRetailShiftFixture(t *testing.T) {
	data := read(t, "upstream/retail_shift.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retailshift.RetailShift
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID == "" || got.Organization == nil || got.RetailStore == nil || got.Store == nil || got.CloseDate.IsZero() ||
				got.Cheque == nil || got.Cheque.Start == nil || got.Cheque.End == nil || got.Cheque.Start.FNNumber != "1234" ||
				len(got.Operations) == 0 || len(got.PaymentOperations) == 0 || got.PaymentOperations[0].Meta.Type != "cashin" {
				t.Fatalf("shift = %+v", got)
			}
			// Открытая смена: closeDate и cheque.end отсутствуют.
			var active retailshift.RetailShift
			if err := codec.Unmarshal([]byte(`{"id":"x","moment":"2026-11-02 09:00:00.000","cheque":{"start":{"shiftNumber":"1"}},"future":1}`), &active); err != nil ||
				!active.CloseDate.IsZero() || active.Cheque == nil || active.Cheque.End != nil || active.Operations != nil {
				t.Fatalf("active = %+v, err = %v", active, err)
			}
		})
	}
}

func TestRetailShiftRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"open", openShift, read(t, "custom/retailshift/create-full.golden.json")},
			{"open minimal", retailshift.CreateRequest{Organization: organization, RetailStore: retailStore},
				[]byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"retailStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailstore/22222222-2222-2222-2222-222222222222","type":"retailstore"}}}`)},
			{"close", closeShift, read(t, "custom/retailshift/update.golden.json")},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[retailshift.CreateRequest](), reflect.TypeFor[retailshift.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Operations",
			"PaymentOperations", "ProceedsCash", "ProceedsNoCash", "ReceivedCash", "ReceivedNoCash", "Files"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestRetailShiftHTTP(t *testing.T) {
	const response = `{"id":"server-id","name":"server-name","proceedsCash":150000}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retailshift.RetailShift, error)
	}{
		{"get", "GET", "entity/retailshift/" + testID, nil, func(ctx context.Context, c *transport.Client) (*retailshift.RetailShift, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/retailshift", read(t, "custom/retailshift/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*retailshift.RetailShift, error) {
			return service(c).Create(ctx, openShift)
		}},
		{"update", "PUT", "entity/retailshift/" + testID, read(t, "custom/retailshift/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retailshift.RetailShift, error) {
			return service(c).Update(ctx, testID, closeShift)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retailshift.RetailShift) {
						if got.ID != "server-id" || got.Name != "server-name" || got.ProceedsCash != 150000 {
							t.Fatalf("shift = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestRetailShiftValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	for name, err := range map[string]error{
		"no organization": func() error {
			_, err := s.Create(t.Context(), retailshift.CreateRequest{RetailStore: retailStore})
			return err
		}(),
		"no retail store": func() error {
			_, err := s.Create(t.Context(), retailshift.CreateRequest{Organization: organization})
			return err
		}(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, retailshift.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", closeShift); return err }(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retailshift") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := retailshift.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestRetailShiftList(t *testing.T) {
	listtest.Run(t, "retailshift", read(t, "upstream/retail_shift.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[retailshift.RetailShift] {
		return service(c).List
	}, func(t *testing.T, got retailshift.RetailShift) {
		if got.ID == "" || got.RetailStore == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
