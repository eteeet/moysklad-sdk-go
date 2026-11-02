package moysklad_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/audit"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const productID = "7944ef04-f831-11e5-7a69-971500188b19"

func newClient(t *testing.T, base, token string) *moysklad.Client {
	t.Helper()
	auth, err := transport.BearerToken(token)
	if err != nil {
		t.Fatal(err)
	}
	c, err := moysklad.NewClient(moysklad.Config{BaseURL: base + "/api/remap/1.2", Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// handles вызывает каждый метод группы без аргументов, кроме перечисленных в skip.
func handles(group any, prefix string, skip ...string) map[string]reflect.Value {
	v := reflect.ValueOf(group)
	out := map[string]reflect.Value{}
	for i := range v.NumMethod() {
		name := v.Type().Method(i).Name
		if !slices.Contains(skip, name) {
			out[prefix+name] = v.Method(i).Call(nil)[0]
		}
	}
	return out
}

// usesTransport: handle хранит именно общий transport.Client.
func usesTransport(handle reflect.Value, shared *transport.Client) bool {
	if handle.Kind() != reflect.Pointer || handle.IsNil() || handle.Elem().Kind() != reflect.Struct {
		return false
	}
	s := handle.Elem()
	for i := range s.NumField() {
		f := s.Field(i)
		if f.Type() == reflect.TypeFor[*transport.Client]() && f.Pointer() == reflect.ValueOf(shared).Pointer() {
			return true
		}
	}
	return false
}

// Каждый пакет entity/* доступен в группе Entity; NewClient и получение handles не делают HTTP;
// все handles держат один transport.Client — общие HTTPClient, Auth, Limiter и Retry.
func TestClientHandles(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.URL.Path)
		mu.Unlock()
	}))
	defer server.Close()
	c := newClient(t, server.URL, "secret-token")

	all := handles(c, "", "Entity", "Report", "Transport")
	for name, h := range handles(c.Entity(), "Entity.") {
		all[name] = h
	}
	for name, h := range handles(c.Report(), "Report.") {
		all[name] = h
	}
	for name, h := range all {
		if !usesTransport(h, c.Transport()) {
			t.Errorf("%s does not use the shared transport", name)
		}
	}

	dirs, err := os.ReadDir("entity")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Context", "Audit", "Security", "Report.Turnover", "Entity.Metadata", "Entity.Files"}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		found := false
		for name := range all {
			if strings.EqualFold(name, "Entity."+dir.Name()) {
				found, want = true, append(want, name)
			}
		}
		if !found {
			t.Errorf("entity/%s has no Entity handle", dir.Name())
		}
	}
	if len(all) != len(want) {
		t.Errorf("handles = %d, want %d: %v", len(all), len(want), want)
	}
	for _, name := range want {
		if _, ok := all[name]; !ok {
			t.Errorf("missing handle %s", name)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if len(requests) != 0 {
		t.Fatalf("hidden requests: %v", requests)
	}
}

// Вызов через фасад и через standalone NewService на c.Transport() идёт с Auth своего клиента;
// у двух клиентов нет общего состояния, токен не попадает в текстовое представление.
func TestClientRequests(t *testing.T) {
	var mu sync.Mutex
	var auths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/product/"+productID {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		mu.Lock()
		auths = append(auths, r.Header.Get("Authorization"))
		mu.Unlock()
		fmt.Fprintf(w, `{"id":%q,"name":"Чай"}`, productID)
	}))
	defer server.Close()
	first, second := newClient(t, server.URL, "first-token"), newClient(t, server.URL, "second-token")

	var info transport.ResponseInfo
	item, err := first.Entity().Product().Get(t.Context(), productID, transport.WithResponseInfo(&info))
	if err != nil || item.Name != "Чай" || info.StatusCode != http.StatusOK {
		t.Fatalf("facade: item = %+v, err = %v, info = %+v", item, err, info)
	}
	standalone, err := product.NewService(first.Transport())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := standalone.Get(t.Context(), productID); err != nil {
		t.Fatalf("standalone: %v", err)
	}
	if _, err := second.Entity().Product().Get(t.Context(), productID); err != nil {
		t.Fatalf("second: %v", err)
	}
	mu.Lock()
	got := strings.Join(auths, ",")
	mu.Unlock()
	if got != "Bearer first-token,Bearer first-token,Bearer second-token" {
		t.Fatalf("authorization = %s", got)
	}
	for _, text := range []string{fmt.Sprintf("%v %+v %#v", first, first, first), fmt.Sprintf("%+v %#v", first.Entity(), first.Report()),
		fmt.Sprintf("%v %+v %#v", first.Transport(), first.Transport(), first.Entity().Product())} {
		if strings.Contains(text, "first-token") {
			t.Fatalf("token leaked: %s", text)
		}
	}
}

func TestNewClientError(t *testing.T) {
	if c, err := moysklad.NewClient(moysklad.Config{BaseURL: "ftp://api.moysklad.ru"}); c != nil || err == nil ||
		!strings.Contains(err.Error(), "moysklad: base URL") {
		t.Fatalf("c = %v, err = %v", c, err)
	}
}

// Каждая группа отправляет запрос через общий transport: Authorization клиента и свой путь.
// Ответ 401 — ошибка вызова, но запрос дошёл с нужным заголовком.
func TestClientGroupsSendThroughTransport(t *testing.T) {
	var mu sync.Mutex
	seen := map[string]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen[r.URL.Path] = r.Header.Get("Authorization")
		mu.Unlock()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"error":"auth","code":1056}]}`))
	}))
	defer server.Close()
	c := newClient(t, server.URL, "group-token")
	ctx := t.Context()
	calls := map[string]func() error{
		"/api/remap/1.2/context/employee": func() error { _, err := c.Context().Employee(ctx); return err },
		"/api/remap/1.2/audit":            func() error { _, err := c.Audit().ListContexts(ctx, audit.ContextsOptions{}); return err },
		"/api/remap/1.2/report/turnover/all": func() error {
			_, err := c.Report().Turnover().All(ctx, turnover.AllOptions{})
			return err
		},
		"/api/remap/1.2/entity/product/" + productID + "/files": func() error {
			_, err := c.Entity().Files().List(ctx, "product", productID, query.PageOptions{})
			return err
		},
		"/api/remap/1.2/entity/customerorder/metadata": func() error {
			_, err := c.Entity().Metadata().Get(ctx, "customerorder")
			return err
		},
	}
	for path, call := range calls {
		if err := call(); sdkerrors.StatusCode(err) != http.StatusUnauthorized {
			t.Errorf("%s: err = %v", path, err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for path := range calls {
		if seen[path] != "Bearer group-token" {
			t.Errorf("%s: authorization = %q, seen = %v", path, seen[path], seen)
		}
	}
}
