// Проверка X01: Delete сервисов — один DELETE без тела, опции и ResponseInfo, 200/204, без повторов.
package entity_test

import (
	"context"
	"net/http"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/audit"
	"github.com/eteeet/moysklad-sdk-go/contexts"
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
	"github.com/eteeet/moysklad-sdk-go/entity/customentity"
	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
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
	"github.com/eteeet/moysklad-sdk-go/entity/pricelist"
	"github.com/eteeet/moysklad-sdk-go/entity/pricetype"
	"github.com/eteeet/moysklad-sdk-go/entity/processingplanfolder"
	"github.com/eteeet/moysklad-sdk-go/entity/processingstage"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/entity/productfolder"
	"github.com/eteeet/moysklad-sdk-go/entity/project"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/region"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashin"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashout"
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
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type deleteFunc func(ctx context.Context, c *transport.Client, id string, opts ...transport.RequestOption) error

func deleteCall[S any](newService func(*transport.Client) (*S, error),
	del func(*S, context.Context, string, ...transport.RequestOption) error) deleteFunc {
	return func(ctx context.Context, c *transport.Client, id string, opts ...transport.RequestOption) error {
		s, err := newService(c)
		if err != nil {
			return err
		}
		return del(s, ctx, id, opts...)
	}
}

func deleteCase[S any](name string, newService func(*transport.Client) (*S, error),
	del func(*S, context.Context, string, ...transport.RequestOption) error) optionsCase {
	call := deleteCall(newService, del)
	return optionsCase{name + " delete", http.MethodDelete, "entity/" + name + "/" + testID,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			return call(ctx, c, testID, opts...)
		}}
}

func deleteCases() []optionsCase {
	return []optionsCase{
		deleteCase("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).Delete),
		deleteCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).Delete),
		deleteCase("bundle", bundle.NewService, (*bundle.Service).Delete),
		deleteCase("cashin", cashin.NewService, (*cashin.Service).Delete),
		deleteCase("cashout", cashout.NewService, (*cashout.Service).Delete),
		deleteCase("counterparty", counterparty.NewService, (*counterparty.Service).Delete),
		deleteCase("country", country.NewService, (*country.Service).Delete),
		deleteCase("currency", currency.NewService, (*currency.Service).Delete),
		deleteCase("customerorder", customerorder.NewService, (*customerorder.Service).Delete),
		deleteCase("demand", demand.NewService, (*demand.Service).Delete),
		deleteCase("employee", employee.NewService, (*employee.Service).Delete),
		deleteCase("enter", enter.NewService, (*enter.Service).Delete),
		deleteCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).Delete),
		deleteCase("invoicein", invoicein.NewService, (*invoicein.Service).Delete),
		deleteCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).Delete),
		deleteCase("loss", loss.NewService, (*loss.Service).Delete),
		deleteCase("move", move.NewService, (*move.Service).Delete),
		deleteCase("organization", organization.NewService, (*organization.Service).Delete),
		deleteCase("paymentin", paymentin.NewService, (*paymentin.Service).Delete),
		deleteCase("paymentout", paymentout.NewService, (*paymentout.Service).Delete),
		deleteCase("product", product.NewService, (*product.Service).Delete),
		deleteCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).Delete),
		deleteCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).Delete),
		deleteCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).Delete),
		deleteCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).Delete),
		deleteCase("retailshift", retailshift.NewService, (*retailshift.Service).Delete),
		deleteCase("retailstore", retailstore.NewService, (*retailstore.Service).Delete),
		deleteCase("role", role.NewService, (*role.Service).Delete),
		deleteCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).Delete),
		deleteCase("service", service.NewService, (*service.Client).Delete),
		deleteCase("supply", supply.NewService, (*supply.Service).Delete),
		deleteCase("taxrate", taxrate.NewService, (*taxrate.Service).Delete),
		deleteCase("uom", uom.NewService, (*uom.Service).Delete),
		deleteCase("processingstage", processingstage.NewService, (*processingstage.Service).Delete),
		deleteCase("processingplanfolder", processingplanfolder.NewService, (*processingplanfolder.Service).Delete),
		deleteCase("pricelist", pricelist.NewService, (*pricelist.Service).Delete),
		deleteCase("retaildrawercashout", retaildrawercashout.NewService, (*retaildrawercashout.Service).Delete),
		deleteCase("retaildrawercashin", retaildrawercashin.NewService, (*retaildrawercashin.Service).Delete),
		deleteCase("prepaymentreturn", prepaymentreturn.NewService, (*prepaymentreturn.Service).Delete),
		deleteCase("prepayment", prepayment.NewService, (*prepayment.Service).Delete),
		deleteCase("payroll", payroll.NewService, (*payroll.Service).Delete),
		deleteCase("commissionreportout", commissionreportout.NewService, (*commissionreportout.Service).Delete),
		deleteCase("commissionreportin", commissionreportin.NewService, (*commissionreportin.Service).Delete),
		deleteCase("counterpartyadjustment", counterpartyadjustment.NewService, (*counterpartyadjustment.Service).Delete),
		deleteCase("factureout", factureout.NewService, (*factureout.Service).Delete),
		deleteCase("facturein", facturein.NewService, (*facturein.Service).Delete),
		deleteCase("inventory", inventory.NewService, (*inventory.Service).Delete),
		deleteCase("internalorder", internalorder.NewService, (*internalorder.Service).Delete),
		deleteCase("aggregatepack", aggregatepack.NewService, (*aggregatepack.Service).Delete),
		deleteCase("specialpricediscount", specialpricediscount.NewService, (*specialpricediscount.Service).Delete),
		deleteCase("personaldiscount", personaldiscount.NewService, (*personaldiscount.Service).Delete),
		deleteCase("accumulationdiscount", accumulationdiscount.NewService, (*accumulationdiscount.Service).Delete),
		deleteCase("contentcard", contentcard.NewService, (*contentcard.Service).Delete),
		deleteCase("consignment", consignment.NewService, (*consignment.Service).Delete),
		deleteCase("customentity", customentity.NewService, (*customentity.Service).Delete),
		deleteCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).Delete),
		deleteCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).Delete),
		deleteCase("contract", contract.NewService, (*contract.Service).Delete),
		deleteCase("project", project.NewService, (*project.Service).Delete),
		deleteCase("group", group.NewService, (*group.Service).Delete),
		deleteCase("productfolder", productfolder.NewService, (*productfolder.Service).Delete),
		deleteCase("store", store.NewService, (*store.Service).Delete),
		deleteCase("variant", variant.NewService, (*variant.Service).Delete),
		{"files delete", http.MethodDelete, "entity/product/" + testID + "/files/" + testID, func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			return s.Delete(ctx, "product", testID, testID, opts...)
		}},
		{"state delete", http.MethodDelete, "entity/customerorder/metadata/states/" + testID, func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := metadata.NewService(c)
			return s.DeleteState(ctx, "customerorder", testID, opts...)
		}},
	}
}

// Опции, ResponseInfo и APIError — общей проверкой T02.
func TestDeleteAcceptsRequestOptions(t *testing.T) {
	for _, tc := range deleteCases() {
		t.Run(tc.name, func(t *testing.T) {
			testServiceOptions(t, tc, false)
			testServiceOptions(t, tc, true)
		})
	}
}

// Запрос без тела; 200 и 204 — успех, другой 2xx — ошибка, 403 — APIError с кодом. 503 при
// включённых повторах не повторяется: DELETE — запись.
func TestDeleteStatuses(t *testing.T) {
	for _, tc := range deleteCases() {
		t.Run(tc.name, func(t *testing.T) {
			for status, ok := range map[int]bool{http.StatusOK: true, http.StatusNoContent: true, http.StatusCreated: false, http.StatusForbidden: false, http.StatusServiceUnavailable: false} {
				var calls atomic.Int32
				server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.ContentLength > 0 || r.Header.Get("Content-Type") != "" {
						t.Errorf("body: length %d, type %q", r.ContentLength, r.Header.Get("Content-Type"))
					}
					w.WriteHeader(status)
					if status == http.StatusForbidden {
						w.Write([]byte(`{"errors":[{"error":"forbidden","code":1016}]}`))
					}
				}))
				c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(t.Context(), c)
				server.Close()
				if (err == nil) != ok || calls.Load() != 1 || status == http.StatusForbidden && !sdkerrors.HasCode(err, 1016) {
					t.Fatalf("status %d: err = %v, calls = %d", status, err, calls.Load())
				}
			}
		})
	}
}

// Неверный id каждого Delete — ошибка до запроса; ResponseInfo сброшен.
func TestDeleteRejectsInvalidID(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var invalid []func(context.Context, ...transport.RequestOption) error
	for _, del := range []deleteFunc{
		deleteCall(bonusprogram.NewService, (*bonusprogram.Service).Delete),
		deleteCall(bonustransaction.NewService, (*bonustransaction.Service).Delete),
		deleteCall(bundle.NewService, (*bundle.Service).Delete),
		deleteCall(cashin.NewService, (*cashin.Service).Delete),
		deleteCall(cashout.NewService, (*cashout.Service).Delete),
		deleteCall(counterparty.NewService, (*counterparty.Service).Delete),
		deleteCall(country.NewService, (*country.Service).Delete),
		deleteCall(currency.NewService, (*currency.Service).Delete),
		deleteCall(customerorder.NewService, (*customerorder.Service).Delete),
		deleteCall(demand.NewService, (*demand.Service).Delete),
		deleteCall(employee.NewService, (*employee.Service).Delete),
		deleteCall(enter.NewService, (*enter.Service).Delete),
		deleteCall(expenseitem.NewService, (*expenseitem.Service).Delete),
		deleteCall(invoicein.NewService, (*invoicein.Service).Delete),
		deleteCall(invoiceout.NewService, (*invoiceout.Service).Delete),
		deleteCall(loss.NewService, (*loss.Service).Delete),
		deleteCall(move.NewService, (*move.Service).Delete),
		deleteCall(organization.NewService, (*organization.Service).Delete),
		deleteCall(paymentin.NewService, (*paymentin.Service).Delete),
		deleteCall(paymentout.NewService, (*paymentout.Service).Delete),
		deleteCall(product.NewService, (*product.Service).Delete),
		deleteCall(purchaseorder.NewService, (*purchaseorder.Service).Delete),
		deleteCall(purchasereturn.NewService, (*purchasereturn.Service).Delete),
		deleteCall(retaildemand.NewService, (*retaildemand.Service).Delete),
		deleteCall(retailsalesreturn.NewService, (*retailsalesreturn.Service).Delete),
		deleteCall(retailshift.NewService, (*retailshift.Service).Delete),
		deleteCall(retailstore.NewService, (*retailstore.Service).Delete),
		deleteCall(role.NewService, (*role.Service).Delete),
		deleteCall(salesreturn.NewService, (*salesreturn.Service).Delete),
		deleteCall(service.NewService, (*service.Client).Delete),
		deleteCall(supply.NewService, (*supply.Service).Delete),
		deleteCall(taxrate.NewService, (*taxrate.Service).Delete),
		deleteCall(uom.NewService, (*uom.Service).Delete),
		deleteCall(processingstage.NewService, (*processingstage.Service).Delete),
		deleteCall(processingplanfolder.NewService, (*processingplanfolder.Service).Delete),
		deleteCall(pricelist.NewService, (*pricelist.Service).Delete),
		deleteCall(retaildrawercashout.NewService, (*retaildrawercashout.Service).Delete),
		deleteCall(retaildrawercashin.NewService, (*retaildrawercashin.Service).Delete),
		deleteCall(prepaymentreturn.NewService, (*prepaymentreturn.Service).Delete),
		deleteCall(prepayment.NewService, (*prepayment.Service).Delete),
		deleteCall(payroll.NewService, (*payroll.Service).Delete),
		deleteCall(commissionreportout.NewService, (*commissionreportout.Service).Delete),
		deleteCall(commissionreportin.NewService, (*commissionreportin.Service).Delete),
		deleteCall(counterpartyadjustment.NewService, (*counterpartyadjustment.Service).Delete),
		deleteCall(factureout.NewService, (*factureout.Service).Delete),
		deleteCall(facturein.NewService, (*facturein.Service).Delete),
		deleteCall(inventory.NewService, (*inventory.Service).Delete),
		deleteCall(internalorder.NewService, (*internalorder.Service).Delete),
		deleteCall(aggregatepack.NewService, (*aggregatepack.Service).Delete),
		deleteCall(specialpricediscount.NewService, (*specialpricediscount.Service).Delete),
		deleteCall(personaldiscount.NewService, (*personaldiscount.Service).Delete),
		deleteCall(accumulationdiscount.NewService, (*accumulationdiscount.Service).Delete),
		deleteCall(contentcard.NewService, (*contentcard.Service).Delete),
		deleteCall(consignment.NewService, (*consignment.Service).Delete),
		deleteCall(customentity.NewService, (*customentity.Service).Delete),
		deleteCall(organizationbranch.NewService, (*organizationbranch.Service).Delete),
		deleteCall(saleschannel.NewService, (*saleschannel.Service).Delete),
		deleteCall(contract.NewService, (*contract.Service).Delete),
		deleteCall(project.NewService, (*project.Service).Delete),
		deleteCall(group.NewService, (*group.Service).Delete),
		deleteCall(productfolder.NewService, (*productfolder.Service).Delete),
		deleteCall(store.NewService, (*store.Service).Delete),
		deleteCall(variant.NewService, (*variant.Service).Delete),
	} {
		for _, id := range []string{"bad", testID + "/../x"} {
			invalid = append(invalid, func(ctx context.Context, opts ...transport.RequestOption) error { return del(ctx, c, id, opts...) })
		}
	}
	invalid = append(invalid,
		func(ctx context.Context, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			return s.Delete(ctx, "product", testID, "bad", opts...)
		},
		func(ctx context.Context, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			return s.Delete(ctx, "Product", testID, testID, opts...)
		},
		func(ctx context.Context, opts ...transport.RequestOption) error {
			s, _ := metadata.NewService(c)
			return s.DeleteState(ctx, "customerorder", "bad", opts...)
		},
		func(ctx context.Context, opts ...transport.RequestOption) error {
			s, _ := metadata.NewService(c)
			return s.DeleteState(ctx, "../x", testID, opts...)
		},
	)
	for i, invoke := range invalid {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(t.Context(), transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Fatalf("case %d: err = %v, info = %+v", i, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

// Read-only ресурсы Delete не получают.
func TestDeleteOnlyWhereDocumented(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[*assortment.Service](), reflect.TypeFor[*pricetype.Service](),
		reflect.TypeFor[*thing.Service](),
		reflect.TypeFor[*saleplatform.Service](),
		reflect.TypeFor[*region.Service](),
		reflect.TypeFor[*contexts.Service](), reflect.TypeFor[*audit.Service](), reflect.TypeFor[*turnover.Service]()} {
		if _, ok := typ.MethodByName("Delete"); ok {
			t.Errorf("%s has Delete", typ)
		}
	}
}
