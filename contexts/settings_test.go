package contexts_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/contexts"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/upstream/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func settingsService(c *transport.Client) *contexts.Service {
	s, _ := contexts.NewService(c)
	return s
}

func TestSettingsFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var company contexts.CompanySettings
		if err := codec.Unmarshal(fixture(t, "company_settings.json"), &company); err != nil {
			t.Fatal(err)
		}
		if company.Meta.Href == "" || company.Currency == nil || len(company.PriceTypes) == 0 || company.PriceTypes[0].ID == "" ||
			company.DiscountStrategy != "bySum" || !company.GlobalOperationNumbering || !company.CheckShippingStock || !company.CheckMinPrice ||
			!company.UseRecycleBin || !company.UseCompanyAddress || company.CompanyAddress != "MyCompany@moysklad.ru" || company.AccountCountry != "RU" {
			t.Fatalf("%s: company = %+v", codec.Name, company)
		}
		var user contexts.UserSettings
		if err := codec.Unmarshal(fixture(t, "user_settings.json"), &user); err != nil {
			t.Fatal(err)
		}
		if user.Meta.Type != "usersettings" || user.DefaultCompany == nil || user.DefaultCompany.Meta.Href != "" ||
			user.DefaultCustomerCounterparty == nil || user.DefaultPlace == nil || user.DefaultProject == nil || user.DefaultPurchaseCounterparty == nil ||
			user.Locale != "ru_RU" || user.MailFooter != "подпись в письме" || user.FieldsPerRow != 3 || user.DefaultScreen != "importcustom" ||
			user.PrintFormat != "pdf" || user.AutoShowReports {
			t.Fatalf("%s: user = %+v", codec.Name, user)
		}
		var sub contexts.Subscription
		if err := codec.Unmarshal(fixture(t, "subscription.json"), &sub); err != nil {
			t.Fatal(err)
		}
		if sub != (contexts.Subscription{Role: "admin", Tariff: "CORPORATE", IsSubscriptionChangeAvailable: true, SubscriptionEndDate: 1642520241922}) {
			t.Fatalf("%s: subscription = %+v", codec.Name, sub)
		}
	}
}

func TestSettingsRequestJSON(t *testing.T) {
	store := model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/store/5f13ac1b-a419-11ea-ac12-000a00000072", Type: "store", MediaType: "application/json"}}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"company", contexts.CompanySettingsRequest{DiscountStrategy: model.Value("byPriority"), UseRecycleBin: model.Value(false), CompanyAddress: model.Value("")},
				`{"discountStrategy":"byPriority","useRecycleBin":false,"companyAddress":""}`},
			{"user", contexts.UserSettingsRequest{DefaultPlace: model.Value(store), DefaultProject: model.Null[model.Reference](), FieldsPerRow: model.Value(0),
				Locale: model.Value("en_US"), AutoShowReports: model.Value(false)},
				`{"defaultPlace":{"meta":{"href":"` + store.Meta.Href + `","type":"store","mediaType":"application/json"}},"defaultProject":null,"fieldsPerRow":0,"locale":"en_US","autoShowReports":false}`},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, []byte(tc.want))
		}
	}
}

func TestSettingsHTTP(t *testing.T) {
	company, user := string(fixture(t, "company_settings.json")), string(fixture(t, "user_settings.json"))
	type op struct {
		name, method, path, body, response string
		invoke                             func(context.Context, *transport.Client) (*string, error)
	}
	href := func(m model.Meta, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &m.Href, nil
	}
	for _, o := range []op{
		{"company get", "GET", "context/companysettings", "", company, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := settingsService(c).CompanySettings(ctx)
			if r == nil {
				return href(model.Meta{}, err)
			}
			return href(r.Meta, err)
		}},
		{"company update", "PUT", "context/companysettings", `{"checkMinPrice":false}`, company, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := settingsService(c).UpdateCompanySettings(ctx, contexts.CompanySettingsRequest{CheckMinPrice: model.Value(false)})
			if r == nil {
				return href(model.Meta{}, err)
			}
			return href(r.Meta, err)
		}},
		{"user get", "GET", "context/usersettings", "", user, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := settingsService(c).UserSettings(ctx)
			if r == nil {
				return href(model.Meta{}, err)
			}
			return href(r.Meta, err)
		}},
		{"user update", "PUT", "context/usersettings", `{"printFormat":"individual"}`, user, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := settingsService(c).UpdateUserSettings(ctx, contexts.UserSettingsRequest{PrintFormat: model.Value("individual")})
			if r == nil {
				return href(model.Meta{}, err)
			}
			return href(r.Meta, err)
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			for _, tc := range documenttest.ResponseCases(o.response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if !strings.HasSuffix(*got, o.path) {
							t.Fatalf("href = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
	c, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"company": func(o ...transport.RequestOption) error {
			_, err := settingsService(c).UpdateCompanySettings(t.Context(), contexts.CompanySettingsRequest{}, o...)
			return err
		},
		"user": func(o ...transport.RequestOption) error {
			_, err := settingsService(c).UpdateUserSettings(t.Context(), contexts.UserSettingsRequest{}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Errorf("%s empty update: err = %v, info = %+v", name, err, info)
		}
	}
}

// Подписка: ответ без meta и id; ошибка — только не-объект и не-200.
func TestSubscriptionHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		ok         bool
	}{
		{"fixture", `{"role":"admin","tariff":"CORPORATE","isSubscriptionChangeAvailable":true,"subscriptionEndDate":1642520241922}`, 200, true},
		{"empty object", `{}`, 200, true},
		{"null", `null`, 200, false},
		{"empty body", ``, 200, false},
		{"array", `[]`, 200, false},
		{"no content", ``, 204, false},
		{"api error", `{"errors":[{"error":"forbidden","code":1016}]}`, 403, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != http.MethodGet || r.URL.Path != "/accountSettings/subscription" || r.URL.RawQuery != "" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			got, err := settingsService(c).Subscription(t.Context())
			if tc.ok != (err == nil) || tc.ok != (got != nil) || calls != 1 {
				t.Fatalf("got = %+v, err = %v, calls = %d", got, err, calls)
			}
			if !tc.ok && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get subscription")) {
				t.Errorf("error = %v", err)
			}
			if tc.name == "fixture" && (got.Tariff != "CORPORATE" || got.SubscriptionEndDate != 1642520241922) {
				t.Errorf("subscription = %+v", got)
			}
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*contexts.Subscription, error) {
		return settingsService(c).Subscription(ctx)
	})
}
