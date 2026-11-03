package variant_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Пример источника: создание характеристик «Размер» и «Цвет» одним POST.
func TestCreateManyCharacteristics(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		io.WriteString(w, `[{"id":"b55d2ddf-60c3-11e7-6adb-ede500000010","name":"Размер","type":"string"},{"id":"b55d2ddf-60c3-11e7-6adb-ede500000011","name":"Цвет","type":"string"}]`)
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := variant.NewService(c)
	r, err := s.CreateManyCharacteristics(t.Context(), []string{"Размер", "Цвет"})
	if err != nil || len(r) != 2 || r[1].Value.Name != "Цвет" {
		t.Fatalf("result = %+v, err = %v", r, err)
	}
	if want := `POST /entity/variant/metadata/characteristics [{"name":"Размер"},{"name":"Цвет"}]`; strings.Join(got, "\n") != want {
		t.Fatalf("requests = %v", got)
	}
	for _, names := range [][]string{nil, {"a", ""}, make([]string, 1001)} {
		info := transport.ResponseInfo{StatusCode: 503}
		if _, err := s.CreateManyCharacteristics(t.Context(), names, transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Errorf("names %d: err = %v", len(names), err)
		}
	}
	if len(got) != 1 {
		t.Fatalf("calls = %d", len(got))
	}
}
