package profit_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/profit"
)

func TestTotalsJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			employee := decode[model.Collection[profit.EmployeeRow]](t, codec, read(t, "upstream/report_profit_by_employee_list.json"))
			want := profit.Totals{SalesCount: 3, SalesAvgCheck: 104498, SellSum: 313495, SellCostSum: 180000, ReturnCount: 1,
				ReturnAvgCheck: 15000, ReturnSum: 15000, ReturnCostSum: 9000, Profit: 124495, Margin: 0.6916388888888889, SalesMargin: 0.3971224421442144}
			if len(employee.Rows) != 1 || employee.Meta.Type != "salesbyemployee" || employee.Rows[0].Employee.Name != "Администратор" ||
				employee.Rows[0].Employee.Meta.Type != "employee" || employee.Rows[0].Totals != want {
				t.Fatalf("employee = %+v", employee)
			}

			counterparty := decode[model.Collection[profit.CounterpartyRow]](t, codec, read(t, "upstream/report_profit_by_counterparty_list.json"))
			if len(counterparty.Rows) != 1 || counterparty.Rows[0].Counterparty.Name != `ООО "Покупатель"` ||
				counterparty.Rows[0].SalesCount != 1 || counterparty.Rows[0].Profit != 110390 || counterparty.Rows[0].Margin != 0.55195 {
				t.Fatalf("counterparty = %+v", counterparty)
			}

			channel := decode[model.Collection[profit.SalesChannelRow]](t, codec, read(t, "upstream/report_profit_by_sales_channel_list.json"))
			if len(channel.Rows) != 1 || channel.Rows[0].SalesChannel.Name != "ijk-shop.ru" || channel.Rows[0].SalesChannel.Type != "ECOMMERCE" ||
				channel.Rows[0].SellCostSum != 150000 || channel.Rows[0].SalesMargin != 0.40594059405940597 {
				t.Fatalf("sales channel = %+v", channel)
			}

			// Счётчики приходят и целыми, и 1.0; деньги — 0, 0.0, дробь, большое значение;
			// неизвестный тип канала не отвергается.
			sparse := decode[profit.SalesChannelRow](t, codec, []byte(`{"salesChannel":{"meta":{"href":"h"},"name":"n","type":"FUTURE_CHANNEL"},`+
				`"salesCount":2,"returnCount":1.0,"sellSum":0,"returnSum":0.0,"salesAvgCheck":1234.56,"profit":310390615323104640.0}`))
			if sparse.SalesChannel.Type != "FUTURE_CHANNEL" || sparse.SalesCount != 2 || sparse.ReturnCount != 1 || sparse.SellSum != 0 ||
				sparse.ReturnSum != 0 || sparse.SalesAvgCheck != 1234.56 || sparse.Profit != 310390615323104640.0 {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

func TestByEmployeeHTTP(t *testing.T) {
	testHTTP(t, "byemployee", "byemployee.json", (*profit.Service).ByEmployee, func(t *testing.T, got *model.Collection[profit.EmployeeRow]) {
		// meta.href примера — byEmployee, путь запроса — byemployee.
		if len(got.Rows) != 1 || got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/report/profit/byEmployee" ||
			got.Rows[0].SalesCount != 3 || got.Rows[0].SellSum != 313495109156366700.0 || got.Rows[0].Employee.Name != "Администратор" {
			t.Fatalf("page = %+v", got)
		}
	})
}

func TestByCounterpartyHTTP(t *testing.T) {
	testHTTP(t, "bycounterparty", "bycounterparty.json", (*profit.Service).ByCounterparty, func(t *testing.T, got *model.Collection[profit.CounterpartyRow]) {
		if len(got.Rows) != 1 || got.Meta.Type != "salesbyCounterparty" || got.Meta.Size != 3 ||
			got.Rows[0].Counterparty.Meta.Type != "counterparty" || got.Rows[0].SalesAvgCheck != 310390615323104640.0 {
			t.Fatalf("page = %+v", got)
		}
	})
}

func TestBySalesChannelHTTP(t *testing.T) {
	testHTTP(t, "bysaleschannel", "bysaleschannel.json", (*profit.Service).BySalesChannel, func(t *testing.T, got *model.Collection[profit.SalesChannelRow]) {
		if len(got.Rows) != 1 || got.Rows[0].SalesChannel.Type != "ECOMMERCE" || got.Rows[0].SalesChannel.Meta.UUIDHref == "" ||
			got.Rows[0].SellSum != 252500 || got.Rows[0].SalesMargin != 1 {
			t.Fatalf("page = %+v", got)
		}
	})
}

func TestNewService(t *testing.T) {
	if _, err := profit.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
