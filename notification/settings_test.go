package notification_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/notification"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

var groups = []string{"customerOrder", "dataExchange", "invoice", "retail", "scripts", "stock", "task", "mentions", "onlineStores", "followedEvents"}

func checkSettingsFixture(t *testing.T, got notification.Settings) {
	t.Helper()
	if len(got) != len(groups) {
		t.Fatalf("groups = %d", len(got))
	}
	for _, name := range groups {
		g, ok := got[name]
		want := map[string]bool{"email": name != "followedEvents", "interface": true, "push": true}
		if !ok || !g.Enabled || len(g.ChannelsEnabled) != 3 {
			t.Fatalf("%s = %+v", name, g)
		}
		for channel, v := range want {
			if got, ok := g.ChannelsEnabled[channel]; !ok || got != v {
				t.Fatalf("%s.%s = %v, %v", name, channel, got, ok)
			}
		}
	}
}

func TestSettingsFixture(t *testing.T) {
	data := readFile(t, "../testdata/upstream/notification_settings.json")
	for _, codec := range documenttest.Codecs {
		var got notification.Settings
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		checkSettingsFixture(t, got)
	}
}

// Группы и каналы — открытые ключи; enabled=false и отсутствие канала различимы.
func TestSettingsOpenKeys(t *testing.T) {
	data := []byte(`{"futureGroup":{"enabled":false,"channelsEnabled":{"sms":true,"email":false}},"stock":{"enabled":true}}`)
	for _, codec := range documenttest.Codecs {
		var got notification.Settings
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		future, stock := got["futureGroup"], got["stock"]
		email, emailOK := future.ChannelsEnabled["email"]
		_, pushOK := future.ChannelsEnabled["push"]
		if len(got) != 2 || future.Enabled || !future.ChannelsEnabled["sms"] || email || !emailOK || pushOK ||
			!stock.Enabled || stock.ChannelsEnabled != nil {
			t.Fatalf("%s: settings = %+v", codec.Name, got)
		}
	}
}

func updateRequest() notification.SettingsRequest {
	return notification.SettingsRequest{
		"customerOrder":  {Enabled: model.Value(false)},
		"task":           {ChannelsEnabled: map[string]bool{"email": false, "push": true}},
		"followedEvents": {Enabled: model.Value(true), ChannelsEnabled: map[string]bool{"interface": false}},
		"futureGroup":    {ChannelsEnabled: map[string]bool{"sms": true}},
	}
}

func TestSettingsRequestJSON(t *testing.T) {
	golden := readFile(t, "../testdata/custom/notification/settings-update.golden.json")
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input notification.SettingsRequest
			want  []byte
		}{
			{"golden", updateRequest(), golden},
			{"enabled only false", notification.SettingsRequest{"stock": {Enabled: model.Value(false)}}, []byte(`{"stock":{"enabled":false}}`)},
			{"empty channels kept", notification.SettingsRequest{"task": {Enabled: model.Value(true), ChannelsEnabled: map[string]bool{}}},
				[]byte(`{"task":{"enabled":true,"channelsEnabled":{}}}`)},
			{"null enabled", notification.SettingsRequest{"task": {Enabled: model.Null[bool]()}}, []byte(`{"task":{"enabled":null}}`)},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			t.Run(codec.Name+" "+tc.name, func(t *testing.T) { documenttest.JSONEqual(t, got, tc.want) })
		}
	}
}

func TestSettingsGet(t *testing.T) {
	body := string(readFile(t, "../testdata/upstream/notification_settings.json"))
	cases := []response{
		{"plain", 200, body, false, false, true},
		{"gzip", 200, body, true, false, true},
		{"empty object", 200, `{}`, false, false, true},
		{"null", 200, `null`, false, false, false},
		{"array", 200, `[]`, false, false, false},
		{"empty body", 200, ``, false, false, false},
		{"malformed", 200, `{"task":`, false, false, false},
		{"string enabled", 200, `{"task":{"enabled":"yes"}}`, false, false, false},
		{"no content", 204, ``, false, false, false},
	}
	for _, r := range append(cases, errorResponses()...) {
		t.Run(r.name, func(t *testing.T) {
			s, calls := serve(t, "GET", "notification/settings", "", "", r)
			var info transport.ResponseInfo
			got, err := s.Settings(t.Context(), transport.WithResponseInfo(&info))
			checkResult(t, r, "get notification settings", calls, info, err)
			switch {
			case !r.ok && got != nil:
				t.Fatalf("got = %+v", got)
			case r.ok && r.body == body:
				checkSettingsFixture(t, got)
			case r.ok && (got == nil || len(got) != 0):
				t.Fatalf("got = %+v", got)
			}
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*notification.Settings, error) {
		got, err := newService(t, c).Settings(ctx)
		if err != nil {
			return nil, err
		}
		return &got, nil
	})
}

func TestSettingsUpdate(t *testing.T) {
	golden := string(readFile(t, "../testdata/custom/notification/settings-update.golden.json"))
	cases := []response{
		{"empty 200", 200, ``, false, false, true},
		{"settings body ignored", 200, `{"task":{"enabled":true}}`, true, false, true},
		{"no content", 204, ``, false, false, true},
	}
	for _, r := range append(cases, errorResponses()...) {
		t.Run(r.name, func(t *testing.T) {
			s, calls := serve(t, "PUT", "notification/settings", "", golden, r)
			var info transport.ResponseInfo
			err := s.UpdateSettings(t.Context(), updateRequest(), transport.WithResponseInfo(&info))
			checkResult(t, r, "update notification settings", calls, info, err)
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
		if err := newService(t, c).UpdateSettings(ctx, updateRequest()); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})

	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{}})
	if err != nil {
		t.Fatal(err)
	}
	s := newService(t, client)
	for name, input := range map[string]notification.SettingsRequest{
		"nil":         nil,
		"empty":       {},
		"empty group": {"task": {}},
		"empty chans": {"task": {ChannelsEnabled: map[string]bool{}}},
		"one of many": {"stock": {Enabled: model.Value(false)}, "task": {}},
	} {
		info := transport.ResponseInfo{StatusCode: 503, Attempts: 1}
		if err := s.UpdateSettings(t.Context(), input, transport.WithResponseInfo(&info)); err == nil ||
			strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "update notification settings") || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("local errors sent %d requests", calls.Load())
	}
	// PUT не повторяется даже при RetryPolicy клиента.
	if err := s.UpdateSettings(t.Context(), updateRequest()); err == nil || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}
