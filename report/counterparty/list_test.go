package counterparty_test

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/counterparty"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const host = "https://api.moysklad.ru/api/remap/1.2/"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *counterparty.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := counterparty.NewService(c)
	return s
}

// dates сравнивает даты строки (пустая строка — null) и возвращает строку с нулевыми датами.
func dates(t *testing.T, row counterparty.Row, first, last, event, updated string) counterparty.Row {
	t.Helper()
	for i, pair := range []struct {
		got  model.DateTime
		want string
	}{{row.FirstDemandDate, first}, {row.LastDemandDate, last}, {row.LastEventDate, event}, {row.Updated, updated}} {
		if pair.want == "" && !pair.got.IsZero() || pair.want != "" && !time.Time(pair.got).Equal(time.Time(documenttest.DateTime(pair.want))) {
			t.Errorf("date %d = %v, want %q", i, time.Time(pair.got), pair.want)
		}
	}
	row.FirstDemandDate, row.LastDemandDate, row.LastEventDate, row.Updated = model.DateTime{}, model.DateTime{}, model.DateTime{}, model.DateTime{}
	return row
}

func meta(path, typ string) model.Meta {
	return model.Meta{Href: host + path, Type: model.MetaType(typ), MediaType: "application/json"}
}

func TestRowJSON(t *testing.T) {
	const id = "df2fdd2d-6934-11e6-8a84-bae500000049"
	party := meta("entity/counterparty/"+id, "counterparty")
	party.MetadataHref = host + "entity/counterparty/metadata"
	want := counterparty.Row{
		Meta: meta("report/counterparty/"+id, "counterparty"),
		Counterparty: counterparty.Counterparty{Meta: party, ID: id, Name: "rtr", ExternalCode: "rRlzrdZmjql9r9dveXPE43",
			CompanyType: model.CompanyTypeLegal},
		DemandsCount: 1, DemandsSum: 60000, AverageReceipt: 60000, Balance: -60000, BonusBalance: 200, Profit: 50000,
		LastEventText: "продажа",
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got counterparty.Row
			if err := codec.Unmarshal(read(t, "upstream/report_counterparty.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got := dates(t, got, "2016-09-01 17:54:00", "2016-09-01 17:54:00", "2016-09-08 13:07:30", "2016-08-23 16:24:08"); got != want {
				t.Fatalf("row = %+v", got)
			}
			var custom counterparty.Row
			if err := codec.Unmarshal(read(t, "custom/reportcounterparty/get.json"), &custom); err != nil {
				t.Fatal(err)
			}
			if got := dates(t, custom, "2016-09-01 17:54:00", "2016-09-01 17:54:00", "2016-09-08 13:07:30", "2016-08-23 16:24:08"); got != want {
				t.Fatalf("custom row = %+v", got)
			}

			var page model.Collection[counterparty.Row]
			if err := codec.Unmarshal(read(t, "custom/reportcounterparty/list.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 5 || page.Meta.Size != 5 || page.Meta.Type != "counterparty" || page.Context == nil || page.Context.Employee == nil {
				t.Fatalf("page = %+v", page)
			}
			// Строка с null-датами и null lastEventText.
			first := dates(t, page.Rows[0], "", "", "", "2016-09-01 18:32:17")
			if first.LastEventText != "" || first.Counterparty.Name != "ООО Радуга" || first.DemandsCount != 0 || first.BonusBalance != 0 {
				t.Fatalf("row 0 = %+v", first)
			}
			if r := page.Rows[2]; r.Profit != -103040600 || r.BonusBalance != 120 || r.DemandsSum != 80000 || r.AverageReceipt != 80000 {
				t.Fatalf("row 2 = %+v", r)
			}
			if r := page.Rows[3]; r.Balance != 1241255000 || r.Counterparty.Name != `ООО "Покупатель"` {
				t.Fatalf("row 3 = %+v", r)
			}
		})
	}
}

// Деньги — копейки как есть: 0, 0.0, дробь и большое значение не теряют точность и знак.
func TestRowMoneyJSON(t *testing.T) {
	data := []byte(`{"meta":{"href":"h"},"counterparty":{"meta":{"href":"c"},"companyType":"futureType"},"averageReceipt":0,"balance":0.0,"demandsSum":150.5,"discountsSum":0.01,` +
		`"profit":-123456789012.25,"returnsSum":999999999999.99,"bonusBalance":12.5,"demandsCount":2147483648,"returnsCount":3,"future":1}`)
	want := counterparty.Row{Meta: model.Meta{Href: "h"}, Counterparty: counterparty.Counterparty{Meta: model.Meta{Href: "c"}, CompanyType: "futureType"},
		DemandsSum: 150.5, DiscountsSum: 0.01, Profit: -123456789012.25, ReturnsSum: 999999999999.99, BonusBalance: 12.5, DemandsCount: 2147483648, ReturnsCount: 3}
	for _, codec := range documenttest.Codecs {
		var got counterparty.Row
		if err := codec.Unmarshal(data, &got); err != nil || got != want {
			t.Fatalf("%s: row = %+v, err = %v", codec.Name, got, err)
		}
	}
}

// Неверные аргументы отклоняются до HTTP; ResponseInfo сброшен.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	check := func(name, operation string, call func(transport.RequestOption) error) {
		t.Helper()
		info := transport.ResponseInfo{StatusCode: 299}
		err := call(transport.WithResponseInfo(&info))
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation) || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	for _, options := range []counterparty.ListOptions{{Limit: -1}, {Limit: 1001}, {Offset: -1}, {Extra: url.Values{"limit": {"5"}}},
		{Extra: url.Values{"offset": {"5"}}}, {Filter: "id=1", Extra: url.Values{"filter": {"id=2"}}}, {Extra: url.Values{"async": {"true"}}}} {
		check("list", "get counterparty report", func(o transport.RequestOption) error {
			_, err := service(c).List(t.Context(), options, o)
			return err
		})
	}
	for name, refs := range map[string][]model.Reference{"nil": nil, "empty": {}, "no href": {{Meta: model.Meta{Href: host + "entity/counterparty/x"}}, {}}} {
		check("selected "+name, "select counterparty report", func(o transport.RequestOption) error {
			_, err := service(c).Selected(t.Context(), refs, o)
			return err
		})
	}
	for _, id := range []string{"", "df2fdd2d-6934-11e6-8a84-bae50000004", "df2fdd2d-6934-11e6-8a84-bae50000004z", "../df2fdd2d-6934-11e6-8a84-bae5000000"} {
		check("get "+id, "get counterparty report", func(o transport.RequestOption) error {
			_, err := service(c).Get(t.Context(), id, o)
			return err
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := counterparty.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
