package bonustransaction_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/bonustransaction"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "7c6ecd51-b738-11e8-727d-307300000028"
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
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *bonustransaction.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := bonustransaction.NewService(c)
	return s
}

var (
	agent   = ref("entity/counterparty/11111111-1111-1111-1111-111111111111", "counterparty")
	program = ref("entity/bonusprogram/22222222-2222-2222-2222-222222222222", "bonusprogram")
	earning = bonustransaction.CreateRequest{
		Agent: agent, BonusProgram: program, TransactionType: "EARNING",
		BonusValue:     model.Value(0),
		ExecutionDate:  model.Value(documenttest.DateTime("2026-11-10 00:00:00.000")),
		Applicable:     model.Value(false),
		ParentDocument: model.Value(ref("entity/retaildemand/33333333-3333-3333-3333-333333333333", "retaildemand")),
	}
	update = bonustransaction.UpdateRequest{BonusValue: model.Value(25), Description: model.Value("")}
)

func TestBonusTransactionFixture(t *testing.T) {
	data := read(t, "upstream/bonus_transaction.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got bonustransaction.BonusTransaction
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Agent == nil || got.BonusProgram == nil || got.ParentDocument == nil || got.BonusValue != 15 ||
				got.TransactionType != "EARNING" || got.TransactionStatus != "COMPLETED" || got.CategoryType != "REGULAR" ||
				got.UpdatedBy != "admin@admin" || got.ExecutionDate.IsZero() {
				t.Fatalf("transaction = %+v", got)
			}
			// Ожидающее списание и неизвестные значения перечислений читаются как есть.
			for _, status := range []string{"WAIT_PROCESSING", "CANCELED", "FUTURE"} {
				var sparse bonustransaction.BonusTransaction
				data := []byte(`{"id":"x","transactionType":"SPENDING","transactionStatus":"` + status + `","categoryType":"WELCOME","bonusValue":5}`)
				if err := codec.Unmarshal(data, &sparse); err != nil || sparse.TransactionStatus != status ||
					sparse.TransactionType != "SPENDING" || sparse.CategoryType != "WELCOME" || sparse.ParentDocument != nil {
					t.Fatalf("sparse = %+v, err = %v", sparse, err)
				}
			}
		})
	}
}

func TestBonusTransactionRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create earning", earning, read(t, "custom/bonustransaction/create-full.golden.json")},
			{"create spending minimal", bonustransaction.CreateRequest{Agent: agent, BonusProgram: program, TransactionType: "SPENDING"},
				[]byte(`{"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/11111111-1111-1111-1111-111111111111","type":"counterparty"}},"bonusProgram":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/bonusprogram/22222222-2222-2222-2222-222222222222","type":"bonusprogram"}},"transactionType":"SPENDING"}`)},
			{"update", update, []byte(`{"bonusValue":25,"description":""}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[bonustransaction.CreateRequest](), reflect.TypeFor[bonustransaction.UpdateRequest]()} {
		names := []string{"ID", "Meta", "AccountID", "Created", "Updated", "UpdatedBy", "TransactionStatus", "CategoryType"}
		if typ.Name() == "UpdateRequest" {
			names = append(names, "TransactionType")
		}
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains %s", typ.Name(), name)
			}
		}
	}
}

func TestBonusTransactionHTTP(t *testing.T) {
	const response = `{"id":"server-id","transactionType":"EARNING","transactionStatus":"WAIT_PROCESSING","bonusValue":10}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*bonustransaction.BonusTransaction, error)
	}{
		{"get", "GET", "entity/bonustransaction/" + testID, nil, func(ctx context.Context, c *transport.Client) (*bonustransaction.BonusTransaction, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/bonustransaction", read(t, "custom/bonustransaction/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*bonustransaction.BonusTransaction, error) {
			return service(c).Create(ctx, earning)
		}},
		{"update", "PUT", "entity/bonustransaction/" + testID, []byte(`{"bonusValue":25,"description":""}`), func(ctx context.Context, c *transport.Client) (*bonustransaction.BonusTransaction, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *bonustransaction.BonusTransaction) {
						if got.ID != "server-id" || got.TransactionStatus != "WAIT_PROCESSING" || got.BonusValue != 10 {
							t.Fatalf("transaction = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Отказ токену решения (WithAuth) возвращается как APIError после одной попытки с этим токеном;
// SDK не повторяет запрос с токеном клиента.
func TestBonusTransactionSolutionTokenDenied(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get("Authorization") != "Bearer solution-token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"error":"Доступ запрещён","code":1016}]}`))
	}))
	defer server.Close()
	personal, err := transport.BearerToken("personal-token")
	if err != nil {
		t.Fatal(err)
	}
	solution, err := transport.BearerToken("solution-token")
	if err != nil {
		t.Fatal(err)
	}
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Auth: personal})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service(c).Create(t.Context(), earning, transport.WithAuth(solution))
	if sdkerrors.StatusCode(err) != http.StatusForbidden || !sdkerrors.HasCode(err, 1016) || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}

func TestBonusTransactionValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	for name, err := range map[string]error{
		"no agent": func() error {
			_, err := s.Create(t.Context(), bonustransaction.CreateRequest{BonusProgram: program, TransactionType: "EARNING"})
			return err
		}(),
		"no program": func() error {
			_, err := s.Create(t.Context(), bonustransaction.CreateRequest{Agent: agent, TransactionType: "EARNING"})
			return err
		}(),
		"no type": func() error {
			_, err := s.Create(t.Context(), bonustransaction.CreateRequest{Agent: agent, BonusProgram: program})
			return err
		}(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, bonustransaction.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "bonustransaction") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := bonustransaction.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestBonusTransactionList(t *testing.T) {
	listtest.Run(t, "bonustransaction", read(t, "upstream/bonus_transaction.json"), "", "", func(c *transport.Client) listtest.ListFunc[bonustransaction.BonusTransaction] {
		return service(c).List
	}, func(t *testing.T, got bonustransaction.BonusTransaction) {
		if got.ID != fixtureID {
			t.Fatalf("row = %+v", got)
		}
	})
}
