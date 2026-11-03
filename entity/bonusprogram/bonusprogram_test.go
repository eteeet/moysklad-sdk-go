package bonusprogram_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/bonusprogram"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "1223d051-ba76-11e8-3353-995e0000005a"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *bonusprogram.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := bonusprogram.NewService(c)
	return s
}

var (
	// Обязательные false и 0 должны попасть в JSON.
	full = bonusprogram.CreateRequest{
		Name: "Бонусы", Active: false, AllProducts: true, AllAgents: false,
		EarnRateRoublesToPoint: 0, SpendRatePointsToRouble: 1, MaxPaidRatePercents: 0,
		AgentTags: model.Value([]string{"vip"}), EarnWhileRedeeming: model.Value(false),
		WelcomeBonusesEnabled: model.Value(true), WelcomeBonusesMode: model.Value("FIRST_PURCHASE"),
		WelcomeBonusesValue: model.Value(0), PostponedBonusesDelayDays: model.Value(14),
	}
	update = bonusprogram.UpdateRequest{Active: model.Value(false), AgentTags: model.Value([]string{}), MaxPaidRatePercents: model.Value(0)}
)

func TestBonusProgramFixture(t *testing.T) {
	data := read(t, "upstream/bonus_program.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got bonusprogram.BonusProgram
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || !got.Active || len(got.AgentTags) != 1 || got.MaxPaidRatePercents != 100 ||
				got.WelcomeBonusesMode != "REGISTRATION" || got.WelcomeBonusesValue != 10 || got.PostponedBonusesDelayDays != 14 {
				t.Fatalf("program = %+v", got)
			}
			var sparse bonusprogram.BonusProgram
			if err := codec.Unmarshal([]byte(`{"id":"x","welcomeBonusesMode":"FUTURE","future":1}`), &sparse); err != nil ||
				sparse.WelcomeBonusesMode != "FUTURE" || sparse.AgentTags != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestBonusProgramRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/bonusprogram/create-full.golden.json")},
			{"create minimal zero", bonusprogram.CreateRequest{Name: "Б"},
				[]byte(`{"name":"Б","active":false,"allProducts":false,"allAgents":false,"earnRateRoublesToPoint":0,"spendRatePointsToRouble":0,"maxPaidRatePercents":0}`)},
			{"update", update, []byte(`{"active":false,"agentTags":[],"maxPaidRatePercents":0}`)},
			{"update null mode", bonusprogram.UpdateRequest{WelcomeBonusesMode: model.Null[string]()}, []byte(`{"welcomeBonusesMode":null}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[bonusprogram.CreateRequest](), reflect.TypeFor[bonusprogram.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestBonusProgramHTTP(t *testing.T) {
	const response = `{"id":"server-id","name":"server-name","active":true,"earnRateRoublesToPoint":5}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*bonusprogram.BonusProgram, error)
	}{
		{"get", "GET", "entity/bonusprogram/" + testID, nil, func(ctx context.Context, c *transport.Client) (*bonusprogram.BonusProgram, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/bonusprogram", read(t, "custom/bonusprogram/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*bonusprogram.BonusProgram, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/bonusprogram/" + testID, []byte(`{"active":false,"agentTags":[],"maxPaidRatePercents":0}`), func(ctx context.Context, c *transport.Client) (*bonusprogram.BonusProgram, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *bonusprogram.BonusProgram) {
						if got.ID != "server-id" || !got.Active || got.EarnRateRoublesToPoint != 5 {
							t.Fatalf("program = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Пустое имя отклоняется до запроса; курсы, проценты и приветственные баллы проверяет сервер.
func TestBonusProgramValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"id":"server-id"}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	if _, err := s.Create(t.Context(), bonusprogram.CreateRequest{Name: "Б", MaxPaidRatePercents: 150, WelcomeBonusesValue: model.Value(-1)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"no name":       func() error { _, err := s.Create(t.Context(), bonusprogram.CreateRequest{Active: true}); return err }(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, bonusprogram.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "bonusprogram") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := bonusprogram.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestBonusProgramList(t *testing.T) {
	listtest.Run(t, "bonusprogram", read(t, "upstream/bonus_program.json"), "", "", func(c *transport.Client) listtest.ListFunc[bonusprogram.BonusProgram] {
		return service(c).List
	}, func(t *testing.T, got bonusprogram.BonusProgram) {
		if got.ID != fixtureID {
			t.Fatalf("row = %+v", got)
		}
	})
}
