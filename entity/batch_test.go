// Проверка X02: CreateMany, UpdateMany и DeleteMany — один POST с массивом, ответ по элементам.
package entity_test

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/accumulationdiscount"
	"github.com/eteeet/moysklad-sdk-go/entity/aggregatepack"
	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	"github.com/eteeet/moysklad-sdk-go/entity/bonusprogram"
	"github.com/eteeet/moysklad-sdk-go/entity/bonustransaction"
	"github.com/eteeet/moysklad-sdk-go/entity/bundle"
	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/entity/commissionreportin"
	"github.com/eteeet/moysklad-sdk-go/entity/commissionreportout"
	"github.com/eteeet/moysklad-sdk-go/entity/consignment"
	"github.com/eteeet/moysklad-sdk-go/entity/contentcard"
	"github.com/eteeet/moysklad-sdk-go/entity/contract"
	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/entity/counterpartyadjustment"
	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/entity/currency"
	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/entity/facturein"
	"github.com/eteeet/moysklad-sdk-go/entity/factureout"
	"github.com/eteeet/moysklad-sdk-go/entity/group"
	"github.com/eteeet/moysklad-sdk-go/entity/internalorder"
	"github.com/eteeet/moysklad-sdk-go/entity/inventory"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/entity/organizationbranch"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/entity/payroll"
	"github.com/eteeet/moysklad-sdk-go/entity/personaldiscount"
	"github.com/eteeet/moysklad-sdk-go/entity/prepayment"
	"github.com/eteeet/moysklad-sdk-go/entity/prepaymentreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/pricetype"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/entity/productfolder"
	"github.com/eteeet/moysklad-sdk-go/entity/project"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/region"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retailshift"
	"github.com/eteeet/moysklad-sdk-go/entity/retailstore"
	"github.com/eteeet/moysklad-sdk-go/entity/role"
	"github.com/eteeet/moysklad-sdk-go/entity/saleplatform"
	"github.com/eteeet/moysklad-sdk-go/entity/saleschannel"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/service"
	"github.com/eteeet/moysklad-sdk-go/entity/specialpricediscount"
	"github.com/eteeet/moysklad-sdk-go/entity/store"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/entity/taxrate"
	"github.com/eteeet/moysklad-sdk-go/entity/thing"
	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const otherID = "88888888-8888-8888-8888-888888888888"

// batchCase отправляет два элемента (testID, otherID) и возвращает по элементу: id записанного
// объекта или код ошибки элемента.
type batchCase struct {
	name, path string
	update     bool // элементы — meta по id с полями
	invoke     func(context.Context, *transport.Client, ...transport.RequestOption) ([]string, error)
}

func createManyCase[S, C, T any](name string, newService func(*transport.Client) (*S, error),
	m func(*S, context.Context, []C, ...transport.RequestOption) ([]batch.Result[T], error)) batchCase {
	return batchCase{name + " create many", "entity/" + name, false,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) ([]string, error) {
			s, err := newService(c)
			if err != nil {
				return nil, err
			}
			res, err := m(s, ctx, []C{fillRequest[C](), fillRequest[C]()}, opts...)
			return results(res), err
		}}
}

func updateManyCase[S, U, T any](name string, newService func(*transport.Client) (*S, error),
	m func(*S, context.Context, []batch.Update[U], ...transport.RequestOption) ([]batch.Result[T], error)) batchCase {
	return batchCase{name + " update many", "entity/" + name, true,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) ([]string, error) {
			s, err := newService(c)
			if err != nil {
				return nil, err
			}
			res, err := m(s, ctx, []batch.Update[U]{{ID: testID, Request: fillRequest[U]()}, {ID: otherID, Request: fillRequest[U]()}}, opts...)
			return results(res), err
		}}
}

func deleteManyCase[S any](name string, newService func(*transport.Client) (*S, error),
	m func(*S, context.Context, []string, ...transport.RequestOption) ([]batch.Deleted, error)) batchCase {
	return batchCase{name + " delete many", "entity/" + name + "/delete", true,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) ([]string, error) {
			s, err := newService(c)
			if err != nil {
				return nil, err
			}
			res, err := m(s, ctx, []string{testID, otherID}, opts...)
			var out []string
			for _, r := range res {
				out = append(out, r.Info+codes(r.Errors))
			}
			return out, err
		}}
}

func results[T any](res []batch.Result[T]) []string {
	var out []string
	for _, r := range res {
		id := ""
		if r.Value != nil {
			id = reflect.ValueOf(r.Value).Elem().FieldByName("ID").String()
		}
		out = append(out, id+codes(r.Errors))
	}
	return out
}

func codes(details []sdkerrors.Detail) string {
	var out []string
	for _, d := range details {
		out = append(out, fmt.Sprint(d.Code))
	}
	return strings.Join(out, ",")
}

func batchCases() []batchCase {
	return []batchCase{
		deleteManyCase("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).DeleteMany),
		createManyCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).CreateMany),
		updateManyCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).UpdateMany),
		deleteManyCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).DeleteMany),
		createManyCase("bundle", bundle.NewService, (*bundle.Service).CreateMany),
		updateManyCase("bundle", bundle.NewService, (*bundle.Service).UpdateMany),
		deleteManyCase("bundle", bundle.NewService, (*bundle.Service).DeleteMany),
		createManyCase("cashin", cashin.NewService, (*cashin.Service).CreateMany),
		updateManyCase("cashin", cashin.NewService, (*cashin.Service).UpdateMany),
		deleteManyCase("cashin", cashin.NewService, (*cashin.Service).DeleteMany),
		createManyCase("cashout", cashout.NewService, (*cashout.Service).CreateMany),
		updateManyCase("cashout", cashout.NewService, (*cashout.Service).UpdateMany),
		deleteManyCase("cashout", cashout.NewService, (*cashout.Service).DeleteMany),
		createManyCase("counterparty", counterparty.NewService, (*counterparty.Service).CreateMany),
		updateManyCase("counterparty", counterparty.NewService, (*counterparty.Service).UpdateMany),
		deleteManyCase("counterparty", counterparty.NewService, (*counterparty.Service).DeleteMany),
		createManyCase("country", country.NewService, (*country.Service).CreateMany),
		updateManyCase("country", country.NewService, (*country.Service).UpdateMany),
		deleteManyCase("country", country.NewService, (*country.Service).DeleteMany),
		createManyCase("currency", currency.NewService, (*currency.Service).CreateMany),
		updateManyCase("currency", currency.NewService, (*currency.Service).UpdateMany),
		deleteManyCase("currency", currency.NewService, (*currency.Service).DeleteMany),
		createManyCase("customerorder", customerorder.NewService, (*customerorder.Service).CreateMany),
		updateManyCase("customerorder", customerorder.NewService, (*customerorder.Service).UpdateMany),
		deleteManyCase("customerorder", customerorder.NewService, (*customerorder.Service).DeleteMany),
		createManyCase("demand", demand.NewService, (*demand.Service).CreateMany),
		updateManyCase("demand", demand.NewService, (*demand.Service).UpdateMany),
		deleteManyCase("demand", demand.NewService, (*demand.Service).DeleteMany),
		updateManyCase("employee", employee.NewService, (*employee.Service).UpdateMany),
		deleteManyCase("employee", employee.NewService, (*employee.Service).DeleteMany),
		createManyCase("enter", enter.NewService, (*enter.Service).CreateMany),
		updateManyCase("enter", enter.NewService, (*enter.Service).UpdateMany),
		deleteManyCase("enter", enter.NewService, (*enter.Service).DeleteMany),
		createManyCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).CreateMany),
		updateManyCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).UpdateMany),
		deleteManyCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).DeleteMany),
		createManyCase("invoicein", invoicein.NewService, (*invoicein.Service).CreateMany),
		updateManyCase("invoicein", invoicein.NewService, (*invoicein.Service).UpdateMany),
		deleteManyCase("invoicein", invoicein.NewService, (*invoicein.Service).DeleteMany),
		createManyCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).CreateMany),
		updateManyCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).UpdateMany),
		deleteManyCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).DeleteMany),
		createManyCase("loss", loss.NewService, (*loss.Service).CreateMany),
		updateManyCase("loss", loss.NewService, (*loss.Service).UpdateMany),
		deleteManyCase("loss", loss.NewService, (*loss.Service).DeleteMany),
		createManyCase("move", move.NewService, (*move.Service).CreateMany),
		updateManyCase("move", move.NewService, (*move.Service).UpdateMany),
		deleteManyCase("move", move.NewService, (*move.Service).DeleteMany),
		createManyCase("organization", organization.NewService, (*organization.Service).CreateMany),
		updateManyCase("organization", organization.NewService, (*organization.Service).UpdateMany),
		deleteManyCase("organization", organization.NewService, (*organization.Service).DeleteMany),
		createManyCase("paymentin", paymentin.NewService, (*paymentin.Service).CreateMany),
		updateManyCase("paymentin", paymentin.NewService, (*paymentin.Service).UpdateMany),
		deleteManyCase("paymentin", paymentin.NewService, (*paymentin.Service).DeleteMany),
		createManyCase("paymentout", paymentout.NewService, (*paymentout.Service).CreateMany),
		updateManyCase("paymentout", paymentout.NewService, (*paymentout.Service).UpdateMany),
		deleteManyCase("paymentout", paymentout.NewService, (*paymentout.Service).DeleteMany),
		createManyCase("product", product.NewService, (*product.Service).CreateMany),
		updateManyCase("product", product.NewService, (*product.Service).UpdateMany),
		deleteManyCase("product", product.NewService, (*product.Service).DeleteMany),
		createManyCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).CreateMany),
		updateManyCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).UpdateMany),
		deleteManyCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).DeleteMany),
		createManyCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).CreateMany),
		updateManyCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).UpdateMany),
		deleteManyCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).DeleteMany),
		createManyCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).CreateMany),
		updateManyCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).UpdateMany),
		deleteManyCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).DeleteMany),
		createManyCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).CreateMany),
		updateManyCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).UpdateMany),
		createManyCase("retailshift", retailshift.NewService, (*retailshift.Service).CreateMany),
		updateManyCase("retailshift", retailshift.NewService, (*retailshift.Service).UpdateMany),
		deleteManyCase("retailshift", retailshift.NewService, (*retailshift.Service).DeleteMany),
		createManyCase("retailstore", retailstore.NewService, (*retailstore.Service).CreateMany),
		updateManyCase("retailstore", retailstore.NewService, (*retailstore.Service).UpdateMany),
		deleteManyCase("retailstore", retailstore.NewService, (*retailstore.Service).DeleteMany),
		createManyCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).CreateMany),
		updateManyCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).UpdateMany),
		deleteManyCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).DeleteMany),
		createManyCase("service", service.NewService, (*service.Client).CreateMany),
		updateManyCase("service", service.NewService, (*service.Client).UpdateMany),
		deleteManyCase("service", service.NewService, (*service.Client).DeleteMany),
		createManyCase("supply", supply.NewService, (*supply.Service).CreateMany),
		updateManyCase("supply", supply.NewService, (*supply.Service).UpdateMany),
		deleteManyCase("supply", supply.NewService, (*supply.Service).DeleteMany),
		createManyCase("taxrate", taxrate.NewService, (*taxrate.Service).CreateMany),
		updateManyCase("taxrate", taxrate.NewService, (*taxrate.Service).UpdateMany),
		deleteManyCase("taxrate", taxrate.NewService, (*taxrate.Service).DeleteMany),
		createManyCase("uom", uom.NewService, (*uom.Service).CreateMany),
		updateManyCase("uom", uom.NewService, (*uom.Service).UpdateMany),
		deleteManyCase("uom", uom.NewService, (*uom.Service).DeleteMany),
		createManyCase("payroll", payroll.NewService, (*payroll.Service).CreateMany),
		updateManyCase("payroll", payroll.NewService, (*payroll.Service).UpdateMany),
		deleteManyCase("payroll", payroll.NewService, (*payroll.Service).DeleteMany),
		createManyCase("commissionreportout", commissionreportout.NewService, (*commissionreportout.Service).CreateMany),
		updateManyCase("commissionreportout", commissionreportout.NewService, (*commissionreportout.Service).UpdateMany),
		deleteManyCase("commissionreportout", commissionreportout.NewService, (*commissionreportout.Service).DeleteMany),
		createManyCase("commissionreportin", commissionreportin.NewService, (*commissionreportin.Service).CreateMany),
		updateManyCase("commissionreportin", commissionreportin.NewService, (*commissionreportin.Service).UpdateMany),
		deleteManyCase("commissionreportin", commissionreportin.NewService, (*commissionreportin.Service).DeleteMany),
		createManyCase("counterpartyadjustment", counterpartyadjustment.NewService, (*counterpartyadjustment.Service).CreateMany),
		updateManyCase("counterpartyadjustment", counterpartyadjustment.NewService, (*counterpartyadjustment.Service).UpdateMany),
		deleteManyCase("counterpartyadjustment", counterpartyadjustment.NewService, (*counterpartyadjustment.Service).DeleteMany),
		createManyCase("factureout", factureout.NewService, (*factureout.Service).CreateMany),
		updateManyCase("factureout", factureout.NewService, (*factureout.Service).UpdateMany),
		deleteManyCase("factureout", factureout.NewService, (*factureout.Service).DeleteMany),
		createManyCase("facturein", facturein.NewService, (*facturein.Service).CreateMany),
		updateManyCase("facturein", facturein.NewService, (*facturein.Service).UpdateMany),
		deleteManyCase("facturein", facturein.NewService, (*facturein.Service).DeleteMany),
		createManyCase("inventory", inventory.NewService, (*inventory.Service).CreateMany),
		updateManyCase("inventory", inventory.NewService, (*inventory.Service).UpdateMany),
		deleteManyCase("inventory", inventory.NewService, (*inventory.Service).DeleteMany),
		createManyCase("internalorder", internalorder.NewService, (*internalorder.Service).CreateMany),
		updateManyCase("internalorder", internalorder.NewService, (*internalorder.Service).UpdateMany),
		deleteManyCase("internalorder", internalorder.NewService, (*internalorder.Service).DeleteMany),
		createManyCase("aggregatepack", aggregatepack.NewService, (*aggregatepack.Service).CreateMany),
		updateManyCase("aggregatepack", aggregatepack.NewService, (*aggregatepack.Service).UpdateMany),
		deleteManyCase("aggregatepack", aggregatepack.NewService, (*aggregatepack.Service).DeleteMany),
		updateManyCase("discount", discount.NewService, (*discount.Service).UpdateMany),
		createManyCase("contentcard", contentcard.NewService, (*contentcard.Service).CreateMany),
		updateManyCase("contentcard", contentcard.NewService, (*contentcard.Service).UpdateMany),
		deleteManyCase("contentcard", contentcard.NewService, (*contentcard.Service).DeleteMany),
		createManyCase("consignment", consignment.NewService, (*consignment.Service).CreateMany),
		updateManyCase("consignment", consignment.NewService, (*consignment.Service).UpdateMany),
		deleteManyCase("consignment", consignment.NewService, (*consignment.Service).DeleteMany),
		createManyCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).CreateMany),
		updateManyCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).UpdateMany),
		deleteManyCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).DeleteMany),
		createManyCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).CreateMany),
		updateManyCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).UpdateMany),
		deleteManyCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).DeleteMany),
		createManyCase("contract", contract.NewService, (*contract.Service).CreateMany),
		updateManyCase("contract", contract.NewService, (*contract.Service).UpdateMany),
		deleteManyCase("contract", contract.NewService, (*contract.Service).DeleteMany),
		createManyCase("project", project.NewService, (*project.Service).CreateMany),
		updateManyCase("project", project.NewService, (*project.Service).UpdateMany),
		deleteManyCase("project", project.NewService, (*project.Service).DeleteMany),
		createManyCase("productfolder", productfolder.NewService, (*productfolder.Service).CreateMany),
		updateManyCase("productfolder", productfolder.NewService, (*productfolder.Service).UpdateMany),
		deleteManyCase("productfolder", productfolder.NewService, (*productfolder.Service).DeleteMany),
		createManyCase("store", store.NewService, (*store.Service).CreateMany),
		updateManyCase("store", store.NewService, (*store.Service).UpdateMany),
		deleteManyCase("store", store.NewService, (*store.Service).DeleteMany),
		createManyCase("variant", variant.NewService, (*variant.Service).CreateMany),
		updateManyCase("variant", variant.NewService, (*variant.Service).UpdateMany),
		deleteManyCase("variant", variant.NewService, (*variant.Service).DeleteMany),
	}
}

// Тело — массив из двух элементов; у update/delete — meta объекта с абсолютным href от BaseURL.
// Ответ: первый элемент записан, второй — с ошибкой; опции вызова и ResponseInfo работают.
func TestBatchRequests(t *testing.T) {
	for _, tc := range batchCases() {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				var body []map[string]jsontext.Value
				if err := json.UnmarshalRead(r.Body, &body); err != nil || r.Method != http.MethodPost || r.URL.RequestURI() != "/"+tc.path ||
					len(body) != 2 || r.Header.Get("Accept-Language") != "en" {
					t.Errorf("request = %s %s %v, body = %v, err = %v", r.Method, r.URL, r.Header, body, err)
				}
				for i, id := range []string{testID, otherID} {
					if !tc.update {
						break
					}
					var meta model.Meta
					entity := strings.TrimSuffix(strings.TrimPrefix(tc.path, "entity/"), "/delete")
					want, metadata := "http://"+r.Host+"/entity/"+entity, entity
					if entity == "bundle" || entity == "service" {
						metadata = "product" // метаданные товаров, комплектов и услуг общие
					}
					if err := json.Unmarshal(body[i]["meta"], &meta); err != nil || meta.Href != want+"/"+id ||
						meta.MetadataHref != "http://"+r.Host+"/entity/"+metadata+"/metadata" || string(meta.Type) != entity || meta.MediaType != "application/json" {
						t.Errorf("items[%d].meta = %+v, err = %v", i, meta, err)
					}
					if strings.HasSuffix(tc.path, "/delete") != (len(body[i]) == 1) {
						t.Errorf("items[%d] fields = %d", i, len(body[i]))
					}
				}
				w.Header().Set("X-Lognex-Content-Region", "kz")
				if strings.HasSuffix(tc.path, "/delete") {
					w.Write([]byte(`[{"info":"ok"},{"errors":[{"error":"x","code":1021}]}]`))
					return
				}
				w.Write([]byte(`[{"id":"` + testID + `"},{"errors":[{"error":"x","code":1000}]}]`))
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			var info transport.ResponseInfo
			got, err := tc.invoke(t.Context(), c, transport.WithAcceptLanguage("en"), transport.WithResponseInfo(&info))
			want := "[" + testID + " 1000]"
			if strings.HasSuffix(tc.path, "/delete") {
				want = "[ok 1021]"
			}
			if err != nil || fmt.Sprint(got) != want || calls.Load() != 1 || info.StatusCode != http.StatusOK ||
				info.Header.Get("X-Lognex-Content-Region") != "kz" {
				t.Fatalf("got = %v, err = %v, calls = %d, info = %+v", got, err, calls.Load(), info)
			}
		})
	}
}

// Ошибка всего запроса: APIError (фикстура batch_global_error), другой статус, ответ не той
// длины, null, null-элементы. 503 при включённых повторах не повторяется: POST — запись.
func TestBatchRequestErrors(t *testing.T) {
	global, err := os.ReadFile("../testdata/upstream/batch/batch_global_error.json")
	if err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct {
		status int
		body   string
		check  func(error) bool
	}{
		"api error":  {http.StatusUnauthorized, string(global), func(err error) bool { return sdkerrors.HasCode(err, 1056) }},
		"503":        {http.StatusServiceUnavailable, "", func(err error) bool { return sdkerrors.StatusCode(err) == http.StatusServiceUnavailable }},
		"201":        {http.StatusCreated, "[{},{}]", func(err error) bool { return err != nil }},
		"short":      {http.StatusOK, "[{}]", func(err error) bool { return err != nil }},
		"null":       {http.StatusOK, "null", func(err error) bool { return err != nil }},
		"not array":  {http.StatusOK, "{}", func(err error) bool { return err != nil }},
		"null items": {http.StatusOK, "[null,null]", func(err error) bool { return err != nil }},
	} {
		t.Run(name, func(t *testing.T) {
			for _, bc := range batchCases() {
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					w.WriteHeader(tc.status)
					w.Write([]byte(tc.body))
				}))
				c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				got, err := bc.invoke(t.Context(), c)
				server.Close()
				if got != nil || !tc.check(err) || calls.Load() != 1 {
					t.Fatalf("%s: got = %v, err = %v, calls = %d", bc.name, got, err, calls.Load())
				}
			}
		})
	}
}

// Неизвестный исход: отмена ctx во время запроса и ошибка 1088 (таймаут обработки массива) —
// ошибка вызова без результатов и без повтора, даже при включённых повторах и статусе 500.
func TestBatchAmbiguousOutcome(t *testing.T) {
	for _, bc := range batchCases() {
		var calls atomic.Int32
		received, done := make(chan struct{}, 1), make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			if r.Header.Get("X-Case") == "cancel" {
				// Тело не дочитано: сервер не заметит обрыв сам, ждём конца вызова.
				received <- struct{}{}
				<-done
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"errors":[{"error":"timeout","code":1088}]}`))
		}))
		c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(t.Context())
		go func() { <-received; cancel() }()
		got, err := bc.invoke(ctx, c, transport.WithHeader("X-Case", "cancel"))
		close(done)
		if got != nil || !errors.Is(err, context.Canceled) || calls.Load() != 1 {
			t.Fatalf("%s cancel: got = %v, err = %v, calls = %d", bc.name, got, err, calls.Load())
		}
		calls.Store(0)
		got, err = bc.invoke(t.Context(), c)
		if got != nil || !sdkerrors.HasCode(err, 1088) || calls.Load() != 1 {
			t.Fatalf("%s 1088: got = %v, err = %v, calls = %d", bc.name, got, err, calls.Load())
		}
		cancel()
		server.Close()
	}
}

// Отмена ctx до отправки: ошибка context.Canceled, запроса нет.
func TestBatchCanceled(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, bc := range batchCases() {
		if _, err := bc.invoke(ctx, c); !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: err = %v", bc.name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

// Лимит 1..1000 элементов, UUID и непустой Request проверяются до запроса; ResponseInfo сброшен.
// 1000 элементов уходят одним запросом.
func TestBatchLocalChecks(t *testing.T) {
	var calls atomic.Int32
	var items atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var body []jsontext.Value
		json.UnmarshalRead(r.Body, &body)
		items.Store(int32(len(body)))
		w.Write([]byte("[" + strings.Repeat(`{"info":"ok"},`, len(body)-1) + `{"info":"ok"}]`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := product.NewService(c)
	ids := func(n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = testID
		}
		return out
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"create none": func(opts ...transport.RequestOption) error {
			_, err := s.CreateMany(t.Context(), nil, opts...)
			return err
		},
		"create 1001": func(opts ...transport.RequestOption) error {
			_, err := s.CreateMany(t.Context(), make([]product.CreateRequest, batch.MaxItems+1), opts...)
			return err
		},
		"update bad id": func(opts ...transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[product.UpdateRequest]{{ID: "bad", Request: fillRequest[product.UpdateRequest]()}}, opts...)
			return err
		},
		"update no fields": func(opts ...transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[product.UpdateRequest]{{ID: testID}}, opts...)
			return err
		},
		"delete none": func(opts ...transport.RequestOption) error {
			_, err := s.DeleteMany(t.Context(), []string{}, opts...)
			return err
		},
		"delete 1001": func(opts ...transport.RequestOption) error {
			_, err := s.DeleteMany(t.Context(), ids(batch.MaxItems+1), opts...)
			return err
		},
		"delete bad id": func(opts ...transport.RequestOption) error {
			_, err := s.DeleteMany(t.Context(), []string{testID, testID + "/../x"}, opts...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Fatalf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if res, err := s.DeleteMany(t.Context(), ids(batch.MaxItems)); err != nil || len(res) != batch.MaxItems || calls.Load() != 1 || items.Load() != batch.MaxItems {
		t.Fatalf("1000: len = %d, err = %v, calls = %d, items = %d", len(res), err, calls.Load(), items.Load())
	}
}

// Массовые методы есть ровно там, где их описывает источник.
func TestBatchOnlyWhereDocumented(t *testing.T) {
	for typ, want := range map[reflect.Type]string{
		reflect.TypeFor[*bonusprogram.Service]():           "DeleteMany",
		reflect.TypeFor[*bonustransaction.Service]():       "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*bundle.Service]():                 "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*cashin.Service]():                 "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*cashout.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*counterparty.Service]():           "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*country.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*currency.Service]():               "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*customerorder.Service]():          "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*demand.Service]():                 "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*employee.Service]():               "UpdateMany DeleteMany",
		reflect.TypeFor[*enter.Service]():                  "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*expenseitem.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*invoicein.Service]():              "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*invoiceout.Service]():             "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*loss.Service]():                   "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*move.Service]():                   "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*organization.Service]():           "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*paymentin.Service]():              "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*paymentout.Service]():             "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*product.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*purchaseorder.Service]():          "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*purchasereturn.Service]():         "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*retaildemand.Service]():           "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*retailsalesreturn.Service]():      "CreateMany UpdateMany",
		reflect.TypeFor[*retailshift.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*retailstore.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*role.Service]():                   "",
		reflect.TypeFor[*salesreturn.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*service.Client]():                 "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*supply.Service]():                 "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*taxrate.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*uom.Service]():                    "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*prepaymentreturn.Service]():       "",
		reflect.TypeFor[*prepayment.Service]():             "",
		reflect.TypeFor[*payroll.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*commissionreportout.Service]():    "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*commissionreportin.Service]():     "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*counterpartyadjustment.Service](): "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*factureout.Service]():             "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*facturein.Service]():              "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*inventory.Service]():              "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*internalorder.Service]():          "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*aggregatepack.Service]():          "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*specialpricediscount.Service]():   "",
		reflect.TypeFor[*personaldiscount.Service]():       "",
		reflect.TypeFor[*accumulationdiscount.Service]():   "",
		reflect.TypeFor[*discount.Service]():               "UpdateMany",
		reflect.TypeFor[*contentcard.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*consignment.Service]():            "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*thing.Service]():                  "",
		reflect.TypeFor[*organizationbranch.Service]():     "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*saleplatform.Service]():           "",
		reflect.TypeFor[*saleschannel.Service]():           "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*region.Service]():                 "",
		reflect.TypeFor[*contract.Service]():               "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*project.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*group.Service]():                  "",
		reflect.TypeFor[*productfolder.Service]():          "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*store.Service]():                  "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*variant.Service]():                "CreateMany UpdateMany DeleteMany",
		reflect.TypeFor[*assortment.Service]():             "DeleteMany",
		reflect.TypeFor[*pricetype.Service]():              "",
	} {
		var got []string
		for _, name := range []string{"CreateMany", "UpdateMany", "DeleteMany"} {
			if _, ok := typ.MethodByName(name); ok {
				got = append(got, name)
			}
		}
		if strings.Join(got, " ") != want {
			t.Errorf("%s: %v, want %q", typ, got, want)
		}
	}
}

// Ассортимент удаляет объекты разных типов по готовым meta; пустые href или type отклоняются
// до запроса.
func TestAssortmentDeleteMany(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var body []struct {
			Meta model.Meta `json:"meta"`
		}
		if err := json.UnmarshalRead(r.Body, &body); err != nil || r.Method != http.MethodPost || r.URL.Path != "/entity/assortment/delete" ||
			len(body) != 2 || body[0].Meta.Type != "product" || body[1].Meta.Type != "service" || body[1].Meta.Href != "https://api.moysklad.ru/x/s1" {
			t.Errorf("request = %s %s, body = %+v, err = %v", r.Method, r.URL, body, err)
		}
		w.Write([]byte(`[{"info":"service"},{"info":"product"}]`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := assortment.NewService(c)
	metas := []model.Meta{{Href: "https://api.moysklad.ru/x/p1", Type: "product"}, {Href: "https://api.moysklad.ru/x/s1", Type: "service"}}
	if res, err := s.DeleteMany(t.Context(), metas); err != nil || len(res) != 2 || res[0].Info != "service" {
		t.Fatalf("res = %+v, err = %v", res, err)
	}
	for _, bad := range [][]model.Meta{nil, {{Href: "h"}}, {{Type: "product"}}, make([]model.Meta, batch.MaxItems+1)} {
		if _, err := s.DeleteMany(t.Context(), bad); err == nil {
			t.Fatalf("metas %d: no error", len(bad))
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
