package counterparty_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	agentID  = "12a8b923-692c-11e6-8a84-bae500000056"
	itemID   = "40702810-1234-1234-1234-123456789012"
	apiPath  = "entity/counterparty/" + agentID
	itemJSON = `{"meta":{"type":"x"},"id":"` + itemID + `","accountNumber":"1","name":"n","description":"d"}`
)

func readNested(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestNestedFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var a counterparty.Account
			if err := codec.Unmarshal(readNested(t, "upstream/counterparty_account.json"), &a); err != nil {
				t.Fatal(err)
			}
			if a.ID != itemID || a.Meta.Type != "account" || !a.IsDefault || a.AccountNumber != "40702810123456789012" ||
				a.BankName != "ПАО Сбербанк" || a.BankLocation == "" || a.CorrespondentAccount != "30101810400000000225" ||
				a.BIC != "044525225" || a.Updated.IsZero() {
				t.Fatalf("account = %+v", a)
			}
			var p counterparty.ContactPerson
			if err := codec.Unmarshal(readNested(t, "upstream/counterparty_contactperson.json"), &p); err != nil {
				t.Fatal(err)
			}
			if p.Name != "Иванов Иван Иванович" || p.Description != "Основной контакт" || p.Email != "ivanov@example.com" ||
				p.Phone == "" || p.Position != "Директор" || p.ExternalCode == "" || p.Agent == nil || p.Agent.Meta.Type != "counterparty" || p.Updated.IsZero() {
				t.Fatalf("contact person = %+v", p)
			}
			var n counterparty.Note
			if err := codec.Unmarshal(readNested(t, "upstream/counterparty_note.json"), &n); err != nil {
				t.Fatal(err)
			}
			if n.Description != "Текст" || n.Agent == nil || n.Author == nil || n.Author.Meta.Type != "employee" ||
				n.AuthorApplication != nil || n.Created.IsZero() {
				t.Fatalf("note = %+v", n)
			}
		})
	}
}

func TestNestedRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"account create", counterparty.AccountCreateRequest{AccountNumber: "86686868768768757656876876"}, `{"accountNumber":"86686868768768757656876876"}`},
			{"account update", counterparty.AccountUpdateRequest{AccountNumber: model.Value("8668"), IsDefault: model.Value(false),
				BankLocation: model.Value("г Москва"), BankName: model.Value("ВТБ"), BIC: model.Value("1005002good"), CorrespondentAccount: model.Value("")},
				`{"accountNumber":"8668","isDefault":false,"bankLocation":"г Москва","bankName":"ВТБ","bic":"1005002good","correspondentAccount":""}`},
			{"contact create", counterparty.ContactPersonCreateRequest{Name: "Петров", Email: model.Value(""), Position: model.Value("Бухгалтер")},
				`{"name":"Петров","email":"","position":"Бухгалтер"}`},
			{"contact update null", counterparty.ContactPersonUpdateRequest{Phone: model.Null[string]()}, `{"phone":null}`},
			{"note", counterparty.NoteRequest{Description: "текст"}, `{"description":"текст"}`},
			// Пример создания контрагента Узбекистана из источника; null у kz отправляется как есть.
			{"requisites", counterparty.UpdateRequest{RequisitesUZ: model.Value(model.RequisitesUZRequest{INN: model.Value("125152124"),
				OKED: model.Value("2013557"), VATPayerRegCode: model.Value("901283746528"), PINFL: model.Value("")}), RequisitesKZ: model.Null[model.RequisitesKZRequest]()},
				`{"mod__requisites__uz":{"inn":"125152124","oked":"2013557","vatPayerRegCode":"901283746528","pinfl":""},"mod__requisites__kz":null}`},
			{"requisites kz", counterparty.CreateRequest{Name: "ТОО", RequisitesKZ: model.Value(model.RequisitesKZRequest{BIN: model.Value("123"),
				VATCertificateDate: model.Value(model.DateTime(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)))})},
				`{"name":"ТОО","mod__requisites__kz":{"bin":"123","vatCertificateDate":"2024-01-02 03:00:00.000"}}`},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, []byte(tc.want))
			})
		}
	}
}

func svc(c *transport.Client) *counterparty.Service {
	s, _ := counterparty.NewService(c)
	return s
}

// Get/Create/Update каждой вложенной коллекции: общий HTTP-контракт (метод, путь, тело, ошибки,
// одна отправка) и ответ-массив из одного объекта у Create, как в примерах источника.
func TestNestedHTTP(t *testing.T) {
	type op struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}
	id := func(v string) *string { return &v }
	wrap := func(get func() (string, error)) (*string, error) {
		v, err := get()
		if err != nil {
			return nil, err
		}
		return id(v), nil
	}
	ops := []op{
		{"get account", "GET", apiPath + "/accounts/" + itemID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				a, err := svc(c).GetAccount(ctx, agentID, itemID)
				return nonNil(a, err, func() string { return a.ID })
			})
		}},
		{"create account", "POST", apiPath + "/accounts", `{"accountNumber":"1"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				a, err := svc(c).CreateAccount(ctx, agentID, counterparty.AccountCreateRequest{AccountNumber: "1"})
				return nonNil(a, err, func() string { return a.ID })
			})
		}},
		{"update account", "PUT", apiPath + "/accounts/" + itemID, `{"isDefault":true}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				a, err := svc(c).UpdateAccount(ctx, agentID, itemID, counterparty.AccountUpdateRequest{IsDefault: model.Value(true)})
				return nonNil(a, err, func() string { return a.ID })
			})
		}},
		{"get contact", "GET", apiPath + "/contactpersons/" + itemID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				p, err := svc(c).GetContactPerson(ctx, agentID, itemID)
				return nonNil(p, err, func() string { return p.ID })
			})
		}},
		{"create contact", "POST", apiPath + "/contactpersons", `{"name":"n"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				p, err := svc(c).CreateContactPerson(ctx, agentID, counterparty.ContactPersonCreateRequest{Name: "n"})
				return nonNil(p, err, func() string { return p.ID })
			})
		}},
		{"update contact", "PUT", apiPath + "/contactpersons/" + itemID, `{"name":"m"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				p, err := svc(c).UpdateContactPerson(ctx, agentID, itemID, counterparty.ContactPersonUpdateRequest{Name: model.Value("m")})
				return nonNil(p, err, func() string { return p.ID })
			})
		}},
		{"get note", "GET", apiPath + "/notes/" + itemID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				n, err := svc(c).GetNote(ctx, agentID, itemID)
				return nonNil(n, err, func() string { return n.ID })
			})
		}},
		{"create note", "POST", apiPath + "/notes", `{"description":"d"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				n, err := svc(c).CreateNote(ctx, agentID, counterparty.NoteRequest{Description: "d"})
				return nonNil(n, err, func() string { return n.ID })
			})
		}},
		{"update note", "PUT", apiPath + "/notes/" + itemID, `{"description":"e"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return wrap(func() (string, error) {
				n, err := svc(c).UpdateNote(ctx, agentID, itemID, counterparty.NoteRequest{Description: "e"})
				return nonNil(n, err, func() string { return n.ID })
			})
		}},
	}
	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			cases := documenttest.ResponseCases(itemJSON)
			if o.method != "GET" {
				cases = append(cases, documenttest.ResponseCase{Name: "array of one", Status: 200, Body: "[" + itemJSON + "]", Success: true},
					documenttest.ResponseCase{Name: "array of two", Status: 200, Body: "[" + itemJSON + "," + itemJSON + "]"},
					documenttest.ResponseCase{Name: "empty array", Status: 200, Body: "[]"})
			}
			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if *got != itemID {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

func nonNil[T any](v *T, err error, id func() string) (string, error) {
	if err != nil {
		return "", err
	}
	return id(), nil
}

// Списки, удаления и массовое удаление счетов: путь, query и тело.
func TestNestedCollections(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[`+itemJSON+`]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := svc(c), t.Context()
	if a, err := s.ListAccounts(ctx, agentID, query.PageOptions{Limit: 5}); err != nil || len(a.Rows) != 1 || a.Rows[0].AccountNumber != "1" {
		t.Fatalf("accounts = %+v, err = %v", a, err)
	}
	if p, err := s.ListContactPersons(ctx, agentID, query.PageOptions{Offset: 10}); err != nil || p.Rows[0].Name != "n" {
		t.Fatalf("contacts = %+v, err = %v", p, err)
	}
	if n, err := s.ListNotes(ctx, agentID, query.PageOptions{}); err != nil || n.Rows[0].Description != "d" {
		t.Fatalf("notes = %+v, err = %v", n, err)
	}
	if err := s.DeleteAccount(ctx, agentID, itemID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteNote(ctx, agentID, itemID); err != nil {
		t.Fatal(err)
	}
	if res, err := s.DeleteManyAccounts(ctx, agentID, []string{itemID}); err != nil || len(res) != 1 || res[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", res, err)
	}
	meta := `{"href":"` + server.URL + "/" + apiPath + `/accounts/` + itemID + `","type":"account","mediaType":"application/json"}`
	want := []string{
		"GET /" + apiPath + "/accounts?limit=5&offset=0 ",
		"GET /" + apiPath + "/contactpersons?limit=1000&offset=10 ",
		"GET /" + apiPath + "/notes?limit=1000&offset=0 ",
		"DELETE /" + apiPath + "/accounts/" + itemID + " ",
		"DELETE /" + apiPath + "/notes/" + itemID + " ",
		"POST /" + apiPath + "/accounts/delete " + `[{"meta":` + meta + `}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// Ошибки входа — до HTTP; ResponseInfo сброшен.
func TestNestedValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := svc(c), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list bad agent": func(o ...transport.RequestOption) error {
			_, err := s.ListAccounts(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"list fields": func(o ...transport.RequestOption) error {
			_, err := s.ListNotes(ctx, agentID, query.PageOptions{Fields: "x"}, o...)
			return err
		},
		"get bad item": func(o ...transport.RequestOption) error {
			_, err := s.GetAccount(ctx, agentID, "bad", o...)
			return err
		},
		"create no number": func(o ...transport.RequestOption) error {
			_, err := s.CreateAccount(ctx, agentID, counterparty.AccountCreateRequest{}, o...)
			return err
		},
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.CreateContactPerson(ctx, agentID, counterparty.ContactPersonCreateRequest{}, o...)
			return err
		},
		"create no text": func(o ...transport.RequestOption) error {
			_, err := s.CreateNote(ctx, agentID, counterparty.NoteRequest{}, o...)
			return err
		},
		"create bad agent": func(o ...transport.RequestOption) error {
			_, err := s.CreateNote(ctx, "../x", counterparty.NoteRequest{Description: "d"}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateAccount(ctx, agentID, itemID, counterparty.AccountUpdateRequest{}, o...)
			return err
		},
		"update empty note": func(o ...transport.RequestOption) error {
			_, err := s.UpdateNote(ctx, agentID, itemID, counterparty.NoteRequest{}, o...)
			return err
		},
		"update bad contact": func(o ...transport.RequestOption) error {
			_, err := s.UpdateContactPerson(ctx, agentID, "x", counterparty.ContactPersonUpdateRequest{Name: model.Value("n")}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.DeleteNote(ctx, agentID, itemID+"/x", o...) },
		"delete many bad": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyAccounts(ctx, agentID, []string{"bad"}, o...)
			return err
		},
		"delete many empty": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyAccounts(ctx, agentID, nil, o...)
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
