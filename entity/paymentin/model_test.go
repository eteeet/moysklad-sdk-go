package paymentin_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestModelJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name, path string
				want       paymentin.PaymentIn
			}{
				{"upstream", "../../testdata/upstream/payment_in.json", paymentin.PaymentIn{
					ID:           "ad1ffc02-3f8e-11e6-8a84-bae5000000fa",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "00001",
					Description:  "Платеж за красивые глаза",
					Code:         "adsfew4r",
					ExternalCode: "8hBoVBUMiDeumkozzjwX83",
					SyncID:       "ad1ffc02-3f8e-11e6-8a84-bae5000000fa",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2016-07-01 16:21:10"),
					Deleted:      documenttest.DateTime("2016-08-25 19:55:00"),
					Moment:       documenttest.DateTime("2016-07-01 16:19:00"),
					Applicable:   true,
					Shared:       false,
					Printed:      true,
					Published:    true,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/organization/metadata",
							Type:         model.MetaType("organization"),
							MediaType:    "application/json",
						},
					}),
					Owner: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/employee/faba7f37-2e58-11e6-8a84-bae500000028",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata",
							Type:         model.MetaType("employee"),
							MediaType:    "application/json",
						},
					}),
					Group: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/group/f97aa1fb-2e58-11e6-8a84-bae500000002",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/group/metadata",
							Type:         model.MetaType("group"),
							MediaType:    "application/json",
						},
					}),
					Project: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/project/722e39f0-313e-11e6-8a84-bae500000008",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/project/metadata",
							Type:         model.MetaType("project"),
							MediaType:    "application/json",
						},
					}),
					State: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/metadata/states/ab8f5213-3f8e-11e6-8a84-bae5000000f5",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/ad1ffc02-3f8e-11e6-8a84-bae5000000fa",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/metadata",
						Type:         model.MetaType("paymentin"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/metadata/attributes/fd6114e2-3f91-11e6-8a84-bae500000107",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "fd6114e2-3f91-11e6-8a84-bae500000107",
							Name:  "Доля",
							Type:  "long",
							Value: model.Value(jsontext.Value("1")),
						},
					},
					Rate: documenttest.Pointer(model.DocumentRate{
						Currency: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/faf45b9a-2e58-11e6-8a84-bae500000055",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
								Type:         model.MetaType("currency"),
								MediaType:    "application/json",
							},
						}),
						Value: documenttest.Pointer(float64(1)),
					}),
					Sum: float64(124124100),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1489a08d-32ca-11e6-8a84-bae50000000d",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata",
							Type:         model.MetaType("counterparty"),
							MediaType:    "application/json",
						},
					}),
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e/accounts/fae39d66-2e58-11e6-8a84-bae50000004f",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1489a08d-32ca-11e6-8a84-bae50000000d/accounts/1489ad8d-32ca-11e6-8a84-bae50000000e",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/contract/2b34d43f-3f52-11e6-8a84-bae500000066",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/contract/metadata",
							Type:         model.MetaType("contract"),
							MediaType:    "application/json",
						},
					}),
					VATSum:         float64(0),
					PaymentPurpose: "За красивые глаза",
					Operations: []model.FinanceOperation{
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/22591b08-2859-11f1-0a83-1f2500000247",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata",
								Type:         model.MetaType("customerorder"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(200000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata",
								Type:         model.MetaType("purchasereturn"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(45000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/demand/7944ef04-047e-11e9-ac12-000a00000002",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata",
								Type:         model.MetaType("demand"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(150000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/e3f4a5b6-1234-1234-1234-e3f4a5b6c7d8",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/metadata",
								Type:         model.MetaType("invoiceout"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(110000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportin/7944ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportin/metadata",
								Type:         model.MetaType("commissionreportin"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(65000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/retailshift/0c979821-f842-11f0-0a80-03f4000a2c6c",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/retailshift/metadata",
								Type:         model.MetaType("retailshift"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(80000),
						},
					},
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/10b0654b-2859-11f1-0a83-1f2500000244",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/metadata",
							Type:         model.MetaType("saleschannel"),
							MediaType:    "application/json",
						},
					}),
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/factureout/0c979821-f842-11f0-0a80-03f4000a2c6c",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/factureout/metadata",
							Type:         model.MetaType("factureout"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#factureout/edit?id=0c979821-f842-11f0-0a80-03f4000a2c6c",
						},
					}),
					IncomingDate:   documenttest.DateTime("2016-07-22 00:00:00"),
					IncomingNumber: "12412569",
				}},
				{"unexpanded complete", "../../testdata/custom/paymentin/read.json", paymentin.PaymentIn{
					ID:           "read-id",
					AccountID:    "read-accountId",
					Name:         "request-name",
					Description:  "",
					Code:         "",
					ExternalCode: "request-externalCode",
					SyncID:       "request-syncId",
					Created:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Updated:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Deleted:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Moment:       documenttest.DateTime("2026-11-01 12:34:56.000"),
					Applicable:   false,
					Shared:       false,
					Printed:      true,
					Published:    true,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organization",
							Type: model.MetaType("organization"),
						},
					}),
					Owner: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/owner",
							Type: model.MetaType("owner"),
						},
					}),
					Group: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/group",
							Type: model.MetaType("group"),
						},
					}),
					Project: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/project",
							Type: model.MetaType("project"),
						},
					}),
					State: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/state",
							Type: model.MetaType("state"),
						},
					}),
					Meta: model.Meta{
						Href: "https://example.test/entity/paymentin/read-id",
						Type: model.MetaType("paymentin"),
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href: "https://example.test/attributes/long",
							},
							Type:  "long",
							Value: model.Value(jsontext.Value("9007199254740993")),
						},
					},
					Rate: documenttest.Pointer(model.DocumentRate{
						Currency: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/currency",
								Type: model.MetaType("currency"),
							},
						}),
					}),
					Sum: float64(123.25),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/agent",
							Type: model.MetaType("agent"),
						},
					}),
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organizationAccount",
							Type: model.MetaType("organizationAccount"),
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/agentAccount",
							Type: model.MetaType("agentAccount"),
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/contract",
							Type: model.MetaType("contract"),
						},
					}),
					VATSum:         float64(0),
					PaymentPurpose: "request-paymentPurpose",
					Operations: []model.FinanceOperation{
						model.FinanceOperation{
							Meta: model.Meta{
								Href: "https://example.test/operation",
								Type: model.MetaType("futureoperation"),
							},
							LinkedSum: float64(1.25),
						},
					},
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/salesChannel",
							Type: model.MetaType("salesChannel"),
						},
					}),
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureOut",
							Type: model.MetaType("factureOut"),
						},
					}),
					IncomingDate:   documenttest.DateTime("2026-11-01 12:34:56.000"),
					IncomingNumber: "00007",
				}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := os.ReadFile(tc.path)
					if err != nil {
						t.Fatal(err)
					}
					var got paymentin.PaymentIn
					if err := codec.Unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					documenttest.NormalizeAttributeJSON(t, got.Attributes)
					documenttest.NormalizeAttributeJSON(t, tc.want.Attributes)
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("model = %+v, want %+v", got, tc.want)
					}
				})
			}
			for _, data := range []string{`{}`, `{"created":null,"updated":null,"deleted":null,"rate":null,"organization":null,"futureField":1}`} {
				var got paymentin.PaymentIn
				if err := codec.Unmarshal([]byte(data), &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != "" || got.Organization != nil || got.Rate != nil || !got.Created.IsZero() || !got.Updated.IsZero() || !got.Deleted.IsZero() {
					t.Fatalf("invented fields: %+v", got)
				}
			}
		})
	}
}
