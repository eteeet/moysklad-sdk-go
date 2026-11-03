package organization_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	orgID       = "4b9d5bec-0575-11e6-9464-e4de00000008"
	accID       = "4b9d69b7-0575-11e6-9464-e4de00000009"
	accPath     = "entity/organization/" + orgID + "/accounts"
	accResponse = `{"meta":{"type":"account"},"id":"` + accID + `","accountNumber":"40702810100000000001","isDefault":true}`
	currencyURL = "https://api.moysklad.ru/api/remap/1.2/entity/currency/0f5fbe16-ceff-11ee-0a83-00490000009d"
)

func accountService(c *transport.Client) *organization.Service {
	s, _ := organization.NewService(c)
	return s
}

func TestAccountFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/organization_account.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		var a organization.Account
		if err := codec.Unmarshal(data, &a); err != nil {
			t.Fatal(err)
		}
		if a.ID != accID || a.Meta.Type != "account" || !a.IsDefault || a.AccountNumber != "40702810100000000001" || a.BankName != "ОАО Сбербанк" ||
			a.BankLocation != "г Москва" || a.CorrespondentAccount != "30101810400000000225" || a.BIC != "044525225" ||
			a.Currency == nil || a.Currency.Meta.Href != currencyURL || a.Updated.IsZero() {
			t.Fatalf("%s: account = %+v", codec.Name, a)
		}
	}
}

func TestAccountRequestJSON(t *testing.T) {
	currency := model.Reference{Meta: model.Meta{Href: currencyURL, Type: "currency", MediaType: "application/json"}}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"create", organization.AccountCreateRequest{AccountNumber: "123", IsDefault: model.Value(false), Currency: model.Value(currency)},
				`{"accountNumber":"123","isDefault":false,"currency":{"meta":{"href":"` + currencyURL + `","type":"currency","mediaType":"application/json"}}}`},
			{"update", organization.AccountUpdateRequest{BankName: model.Value(""), BIC: model.Null[string]()}, `{"bankName":"","bic":null}`},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, []byte(tc.want))
		}
	}
}

// Get/Create/Update: общий HTTP-контракт; Create принимает массив из одного счёта (пример источника).
func TestAccountHTTP(t *testing.T) {
	id := func(a *organization.Account, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &a.ID, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"get", "GET", accPath + "/" + accID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(accountService(c).GetAccount(ctx, orgID, accID))
		}},
		{"create", "POST", accPath, `{"accountNumber":"1"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(accountService(c).CreateAccount(ctx, orgID, organization.AccountCreateRequest{AccountNumber: "1"}))
		}},
		{"update", "PUT", accPath + "/" + accID, `{"isDefault":false}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(accountService(c).UpdateAccount(ctx, orgID, accID, organization.AccountUpdateRequest{IsDefault: model.Value(false)}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			cases := documenttest.ResponseCases(accResponse)
			if o.method != "GET" {
				cases = append(cases, documenttest.ResponseCase{Name: "array of one", Status: 200, Body: "[" + accResponse + "]", Success: true},
					documenttest.ResponseCase{Name: "array of two", Status: 200, Body: "[" + accResponse + "," + accResponse + "]"})
			}
			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if *got != accID {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

// Список, удаление и массовые операции: «Изменить счета Юрлица» — один POST в accounts;
// новые элементы без meta, изменяемые — с meta по ID.
func TestAccountCollections(t *testing.T) {
	var got []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[`+accResponse+`]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		case r.Method == http.MethodPost:
			io.WriteString(w, `[`+accResponse+`]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := accountService(c), t.Context()
	if l, err := s.ListAccounts(ctx, orgID, query.PageOptions{Limit: 2}); err != nil || len(l.Rows) != 1 || !l.Rows[0].IsDefault {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	if err := s.DeleteAccount(ctx, orgID, accID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.CreateManyAccounts(ctx, orgID, []organization.AccountCreateRequest{{AccountNumber: "2"}}); err != nil || r[0].Value.ID != accID {
		t.Fatalf("create many = %+v, err = %v", r, err)
	}
	if _, err := s.UpdateManyAccounts(ctx, orgID, []batch.Update[organization.AccountUpdateRequest]{{ID: accID, Request: organization.AccountUpdateRequest{IsDefault: model.Value(true)}}}); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteManyAccounts(ctx, orgID, []string{accID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + "/" + accPath + "/" + accID + `","type":"account","mediaType":"application/json"}`
	want := []string{
		"GET /" + accPath + "?limit=2&offset=0 ",
		"DELETE /" + accPath + "/" + accID + " ",
		"POST /" + accPath + ` [{"accountNumber":"2"}]`,
		"POST /" + accPath + ` [{"meta":` + meta + `,"isDefault":true}]`,
		"POST /" + accPath + `/delete [{"meta":` + meta + `}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestAccountValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := accountService(c), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list bad org": func(o ...transport.RequestOption) error {
			_, err := s.ListAccounts(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"get bad account": func(o ...transport.RequestOption) error { _, err := s.GetAccount(ctx, orgID, "x", o...); return err },
		"create no number": func(o ...transport.RequestOption) error {
			_, err := s.CreateAccount(ctx, orgID, organization.AccountCreateRequest{}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateAccount(ctx, orgID, accID, organization.AccountUpdateRequest{}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.DeleteAccount(ctx, "../x", accID, o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyAccounts(ctx, orgID, nil, o...)
			return err
		},
		"update many bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyAccounts(ctx, orgID, []batch.Update[organization.AccountUpdateRequest]{{ID: "bad", Request: organization.AccountUpdateRequest{BIC: model.Value("1")}}}, o...)
			return err
		},
		"delete many bad org": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyAccounts(ctx, "bad", []string{accID}, o...)
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
