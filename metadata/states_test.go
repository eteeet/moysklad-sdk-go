package metadata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const stateID = "4dcb3f23-60c4-11e7-6adb-ede500000019"

func color(c int64) *int64 { return &c }

var (
	createState = metadata.StateCreateRequest{Name: "Одобрено", Color: color(10667543), StateType: metadata.StateSuccessful}
	// Нулевой цвет — явное значение, а не «не передано».
	updateState = metadata.StateUpdateRequest{Color: model.Value(int64(0)), StateType: model.Value(metadata.StateUnsuccessful)}
)

func TestStateRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create", createState, readFile(t, "custom/metadata/state-create.golden.json")},
			{"create zero color", metadata.StateCreateRequest{Name: "Новый", Color: color(0), StateType: "Future"}, []byte(`{"name":"Новый","color":0,"stateType":"Future"}`)},
			{"update", updateState, readFile(t, "custom/metadata/state-update.golden.json")},
			{"update name", metadata.StateUpdateRequest{Name: model.Value("Переименован")}, []byte(`{"name":"Переименован"}`)},
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
}

// Цвет ARGB с прозрачностью больше int32 и неизвестный тип статуса читаются.
func TestStateRead(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got metadata.State
			if err := codec.Unmarshal([]byte(`{"id":"x","color":4294967295,"stateType":"Future","entityType":"demand"}`), &got); err != nil ||
				got.Color != 4294967295 || got.StateType != "Future" {
				t.Fatalf("state = %+v, err = %v", got, err)
			}
		})
	}
}

func TestStateHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata/states/` + stateID +
		`"},"id":"` + stateID + `","name":"Одобрено","color":10667543,"stateType":"Successful","entityType":"counterparty"}`
	for _, op := range []struct {
		method, path string
		body         []byte
		invoke       func(context.Context, *transport.Client) (*metadata.State, error)
	}{
		{"POST", "entity/counterparty/metadata/states", readFile(t, "custom/metadata/state-create.golden.json"), func(ctx context.Context, c *transport.Client) (*metadata.State, error) {
			return service(c).CreateState(ctx, "counterparty", createState)
		}},
		{"PUT", "entity/counterparty/metadata/states/" + stateID, readFile(t, "custom/metadata/state-update.golden.json"), func(ctx context.Context, c *transport.Client) (*metadata.State, error) {
			return service(c).UpdateState(ctx, "counterparty", stateID, updateState)
		}},
	} {
		t.Run(op.method, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *metadata.State) {
						if got.ID != stateID || got.Color != 10667543 {
							t.Fatalf("state = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestStateValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(client)
	info := transport.ResponseInfo{StatusCode: 503}
	for name, call := range map[string]func() error{
		"create no name": func() error {
			_, err := s.CreateState(t.Context(), "demand", metadata.StateCreateRequest{Color: color(1), StateType: "Regular"}, transport.WithResponseInfo(&info))
			return err
		},
		"create no color": func() error {
			_, err := s.CreateState(t.Context(), "demand", metadata.StateCreateRequest{Name: "x", StateType: "Regular"})
			return err
		},
		"create no type": func() error {
			_, err := s.CreateState(t.Context(), "demand", metadata.StateCreateRequest{Name: "x", Color: color(1)})
			return err
		},
		"create bad type": func() error { _, err := s.CreateState(t.Context(), "../x", createState); return err },
		"update empty": func() error {
			_, err := s.UpdateState(t.Context(), "demand", stateID, metadata.StateUpdateRequest{})
			return err
		},
		"update bad id":   func() error { _, err := s.UpdateState(t.Context(), "demand", "bad", updateState); return err },
		"update bad type": func() error { _, err := s.UpdateState(t.Context(), "../x", stateID, updateState); return err },
	} {
		if err := call(); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	if calls.Load() != 0 || info.StatusCode != 0 {
		t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
	}
}
