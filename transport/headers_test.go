package transport_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestTypedRequestHeaders(t *testing.T) {
	var got http.Header
	client := newOptionsClient(t, http.Header{"X-Lognex-Accept-Region": {"ru"}}, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	})
	now := time.Date(2017, 8, 16, 23, 7, 1, 0, time.FixedZone("", 7*60*60))

	_, err := client.Send(t.Context(), http.MethodGet, "report/turnover/all", nil, nil, nil,
		transport.WithAcceptRegion(transport.RegionKZ),
		transport.WithAcceptLanguage("en"),
		transport.WithAcceptTimezone(now),
		transport.WithBetaFeature("first"),
		transport.WithBetaFeature("assortmentWithoutStock"),
		transport.WithWebHookDisableByPrefix("https://a.example", "https://b.example/hook"),
		transport.WithDownloadExpiration(time.Hour+500*time.Millisecond),
		transport.WithHeader(transport.HeaderPrettyPrintJSON, "true"),
	)
	if err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{
		transport.HeaderAcceptRegion:           "kz",
		transport.HeaderAcceptLanguage:         "en",
		transport.HeaderAcceptTimezone:         "Wed, 16 Aug 2017 23:07:01 +0700",
		transport.HeaderBetaFeature:            "assortmentWithoutStock",
		transport.HeaderWebHookDisableByPrefix: "https://a.example|https://b.example/hook",
		transport.HeaderDownloadExpiration:     "3600",
		transport.HeaderPrettyPrintJSON:        "true",
		"X-Lognex-WebHook-Disable":             "",
	} {
		if v := strings.Join(got.Values(key), "|"); v != want {
			t.Errorf("%s = %q, want %q", key, v, want)
		}
	}

	if _, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil, transport.WithWebHookDisableByPrefix()); err != nil {
		t.Fatal(err)
	}
	if _, ok := got[http.CanonicalHeaderKey(transport.HeaderWebHookDisableByPrefix)]; ok || got.Get(transport.HeaderAcceptRegion) != "ru" {
		t.Fatalf("headers = %v", got)
	}
}

func TestResponseHeaderValues(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		value string
		want  int64
		ok    bool
	}{
		{"", 0, false},
		{"45", 45, true},
		{" 7 ", 7, true},
		{"0", 0, true},
		{"abc", 0, false},
		{"1.5", 0, false},
		{"-1", 0, false},
		{"99999999999999999999", 0, false},
	} {
		info := transport.ResponseInfo{Header: http.Header{}}
		if tc.value != "" {
			for _, key := range []string{transport.HeaderRateLimitLimit, transport.HeaderRateLimitRemaining,
				transport.HeaderRetryTimeInterval, transport.HeaderReset, transport.HeaderLognexRetryAfter} {
				info.Header.Set(key, tc.value)
			}
		}
		for name, got := range map[string]func() (int64, bool){
			"limit":     info.RateLimitLimit,
			"remaining": info.RateLimitRemaining,
			"interval":  func() (int64, bool) { d, ok := info.RetryTimeInterval(); return d.Milliseconds(), ok },
			"reset":     func() (int64, bool) { d, ok := info.ResetAfter(); return d.Milliseconds(), ok },
			"retry":     func() (int64, bool) { d, ok := info.LognexRetryAfter(); return d.Milliseconds(), ok },
		} {
			if v, ok := got(); v != tc.want && tc.ok || ok != tc.ok {
				t.Errorf("%s(%q) = %d, %v", name, tc.value, v, ok)
			}
		}
		if info.Header.Get(transport.HeaderReset) != tc.value {
			t.Errorf("raw value lost: %q", info.Header.Get(transport.HeaderReset))
		}
	}

	// Миллисекунды, которые не помещаются в time.Duration.
	info := transport.ResponseInfo{Header: http.Header{transport.HeaderReset: {"9223372036854775807"}}}
	if _, ok := info.ResetAfter(); ok {
		t.Error("overflowing duration accepted")
	}

	for _, tc := range []struct {
		value string
		want  time.Duration
		ok    bool
	}{
		{"120", 2 * time.Minute, true},
		{"Sun, 01 Nov 2026 12:00:30 GMT", 30 * time.Second, true},
		{"Sun, 01 Nov 2026 11:00:00 GMT", 0, true},
		{"1500ms", 0, false},
		{"-5", 0, false},
		{"9223372036854775807", 0, false},
		{"", 0, false},
	} {
		info := transport.ResponseInfo{Header: http.Header{transport.HeaderRetryAfter: {tc.value}}}
		if got, ok := info.RetryAfter(now); got != tc.want || ok != tc.ok {
			t.Errorf("RetryAfter(%q) = %v, %v", tc.value, got, ok)
		}
	}

	info = transport.ResponseInfo{Header: http.Header{transport.HeaderContentTimezone: {"Wed, 16 Aug 2017 23:07:01 +0700"}}}
	if date, ok := info.ContentTimezone(); !ok || date.Format(time.RFC3339) != "2017-08-16T23:07:01+07:00" {
		t.Errorf("ContentTimezone = %v, %v", date, ok)
	}
	info.Header.Set(transport.HeaderContentTimezone, "Asia/Novosibirsk")
	if _, ok := info.ContentTimezone(); ok {
		t.Error("timezone ID accepted as date")
	}
}

// Заголовок зоны отчёта не меняет разбор DateTime: даты API остаются московскими.
func TestContentTimezoneDoesNotChangeDateTime(t *testing.T) {
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(transport.HeaderContentTimezone, "Wed, 16 Aug 2017 23:07:01 +0700")
		w.Write([]byte(`{"updated":"2017-08-16 12:00:00.000"}`))
	})
	var first, second struct {
		Updated model.DateTime `json:"updated"`
	}
	for _, dst := range []any{&first, &second} {
		if _, err := client.Send(t.Context(), http.MethodGet, "report/turnover/all", nil, nil, dst,
			transport.WithAcceptTimezone(time.Date(2017, 8, 16, 23, 7, 1, 0, time.FixedZone("", 7*60*60)))); err != nil {
			t.Fatal(err)
		}
	}
	for _, got := range []model.DateTime{first.Updated, second.Updated} {
		if _, offset := time.Time(got).Zone(); offset != 3*60*60 {
			t.Fatalf("offset = %d", offset)
		}
	}
}

func TestSendCanceledContext(t *testing.T) {
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		t.Error("request sent with canceled context")
	})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	info := transport.ResponseInfo{StatusCode: 200}
	_, err := client.Send(ctx, http.MethodGet, "entity/uom", nil, nil, nil, transport.WithResponseInfo(&info))
	if err == nil || !errors.Is(err, context.Canceled) || info.StatusCode != 0 {
		t.Fatalf("err = %v, info = %+v", err, info)
	}
}

func TestSendCanceledInFlight(t *testing.T) {
	received := make(chan struct{})
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		close(received)
		<-r.Context().Done()
	})
	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		<-received
		cancel()
	}()
	info := transport.ResponseInfo{StatusCode: 200}
	_, err := client.Send(ctx, http.MethodGet, "entity/uom", nil, nil, nil,
		transport.WithAcceptLanguage("en"), transport.WithResponseInfo(&info))
	if !errors.Is(err, context.Canceled) || info.StatusCode != 0 || info.Header != nil {
		t.Fatalf("err = %v, info = %+v", err, info)
	}
}
