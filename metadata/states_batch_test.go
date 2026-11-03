package metadata_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Пример источника: новый статус и изменение существующего одним POST в states.
func TestStatesBatch(t *testing.T) {
	const stateID = "b56215dc-60c3-11e7-6adb-ede500000013"
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		io.WriteString(w, `[{"meta":{"type":"state"},"id":"`+stateID+`","name":"На подписании","color":34617,"stateType":"Regular"}]`)
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := metadata.NewService(c)
	color := int64(8767198)
	if r, err := s.CreateManyStates(t.Context(), "counterparty", []metadata.StateCreateRequest{{Name: "На рассмотрении", Color: &color, StateType: metadata.StateRegular}}); err != nil || r[0].Value.ID != stateID {
		t.Fatalf("create many = %+v, err = %v", r, err)
	}
	if r, err := s.UpdateManyStates(t.Context(), "counterparty", []batch.Update[metadata.StateUpdateRequest]{{ID: stateID,
		Request: metadata.StateUpdateRequest{Name: model.Value("На подписании"), Color: model.Value(int64(34617))}}}); err != nil || r[0].Value.Color != 34617 {
		t.Fatalf("update many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + `/entity/counterparty/metadata/states/` + stateID + `","type":"state","mediaType":"application/json"}`
	want := []string{
		`POST /entity/counterparty/metadata/states [{"name":"На рассмотрении","color":8767198,"stateType":"Regular"}]`,
		`POST /entity/counterparty/metadata/states [{"meta":` + meta + `,"name":"На подписании","color":34617}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"create bad type": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyStates(t.Context(), "Counter", []metadata.StateCreateRequest{{Name: "n"}}, o...)
			return err
		},
		"create empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyStates(t.Context(), "counterparty", nil, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyStates(t.Context(), "counterparty", []batch.Update[metadata.StateUpdateRequest]{{ID: "x", Request: metadata.StateUpdateRequest{Name: model.Value("n")}}}, o...)
			return err
		},
		"update bad type": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyStates(t.Context(), "../x", []batch.Update[metadata.StateUpdateRequest]{{ID: stateID, Request: metadata.StateUpdateRequest{Name: model.Value("n")}}}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if len(got) != 2 {
		t.Fatalf("requests = %d", len(got))
	}
}
