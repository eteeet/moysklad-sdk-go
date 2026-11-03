package money_test

import (
	"context"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/money"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const organization = "https://api.moysklad.ru/api/remap/1.2/entity/organization/850c8195-f504-11e5-8a84-bae50000015e"

func TestByAccountJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, path := range []string{"upstream/report_money_by_account_list.json", "custom/money/byaccount.json"} {
				var got model.Collection[money.AccountRow]
				if err := codec.Unmarshal(read(t, path), &got); err != nil {
					t.Fatal(err)
				}
				if len(got.Rows) != 3 || got.Meta.Type != "moneyreport" || got.Context == nil || got.Meta.Size != 0 {
					t.Fatalf("%s: page = %+v", path, got)
				}
				cash := got.Rows[0]
				if cash.Account != nil || cash.Balance != 100 || cash.Organization.Name != "OOO Ромашка" || cash.Organization.Meta.Href != organization {
					t.Fatalf("%s: cash = %+v", path, cash)
				}
				account := got.Rows[2]
				if account.Account == nil || account.Account.Name != "00000" || account.Account.Meta.Type != "account" || account.Balance != 200 ||
					account.Organization.Name != "OOO Серьезное Юридическое Лицо" {
					t.Fatalf("%s: account = %+v", path, account)
				}
			}
		})
	}
}

// Деньги — копейки как есть: 0, 0.0, дробь и большое значение.
func TestByAccountMoneyJSON(t *testing.T) {
	for raw, want := range map[string]float64{"0": 0, "0.0": 0, "150.5": 150.5, "-0.01": -0.01, "123456789012345.67": 123456789012345.67} {
		for _, codec := range documenttest.Codecs {
			var got money.AccountRow
			if err := codec.Unmarshal([]byte(`{"organization":{"meta":{"href":"o"},"name":"n"},"balance":`+raw+`}`), &got); err != nil || got.Balance != want || got.Account != nil {
				t.Fatalf("%s %s: row = %+v, err = %v", codec.Name, raw, got, err)
			}
		}
	}
}

func TestByAccountHTTP(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[money.AccountRow], error) {
		return service(c).ByAccount(ctx)
	}
	for _, rc := range documenttest.ResponseCases(string(read(t, "custom/money/byaccount.json"))) {
		t.Run(rc.Name, func(t *testing.T) {
			documenttest.HTTP(t, rc, "GET", "report/money/byaccount", nil, invoke, func(got *model.Collection[money.AccountRow]) {
				if len(got.Rows) != 3 || got.Rows[1].Organization.Name != "OOO Серьезное Юридическое Лицо" || got.Rows[1].Balance != 100 {
					t.Fatalf("page = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
