package retailstore

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// RetailStore — точка продаж. Модель чтения без expand; проценты (bankPercent, qrBankPercent,
// discountMaxPercent) — не деньги. discountMaxPercent таблица называет Int, но ответы содержат
// 100.0 — SDK читает и пишет float64.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore
type RetailStore struct {
	Meta      model.Meta     `json:"meta"`
	ID        string         `json:"id"`
	AccountID string         `json:"accountId"`
	Updated   model.DateTime `json:"updated,omitzero"`
	Name      string         `json:"name"`
	// ExternalCode у точки продаж только для чтения.
	ExternalCode                        string `json:"externalCode"`
	Address                             string `json:"address"`
	DemandPrefix                        string `json:"demandPrefix"`
	Description                         string `json:"description"`
	IDQR                                string `json:"idQR"`
	QRTerminalID                        string `json:"qrTerminalId"`
	Active                              bool   `json:"active"`
	AllowCreateProducts                 bool   `json:"allowCreateProducts"`
	AllowDeleteReceiptPositions         bool   `json:"allowDeleteReceiptPositions"`
	AllowSellTobaccoWithoutMRC          bool   `json:"allowSellTobaccoWithoutMRC"`
	Archived                            bool   `json:"archived"`
	ControlCashierChoice                bool   `json:"controlCashierChoice"`
	ControlShippingStock                bool   `json:"controlShippingStock"`
	CreateCashInOnRetailShiftClosing    bool   `json:"createCashInOnRetailShiftClosing"`
	CreatePaymentInOnRetailShiftClosing bool   `json:"createPaymentInOnRetailShiftClosing"`
	DiscountEnable                      bool   `json:"discountEnable"`
	EnableReturnsWithNoReason           bool   `json:"enableReturnsWithNoReason"`
	IssueOrders                         bool   `json:"issueOrders"`
	OnlyInStock                         bool   `json:"onlyInStock"`
	PrintAlways                         bool   `json:"printAlways"`
	QRPayEnabled                        bool   `json:"qrPayEnabled"`
	RequiredFio                         bool   `json:"requiredFio"`
	RequiredPhone                       bool   `json:"requiredPhone"`
	RequiredEmail                       bool   `json:"requiredEmail"`
	RequiredBirthdate                   bool   `json:"requiredBirthdate"`
	RequiredSex                         bool   `json:"requiredSex"`
	RequiredDiscountCardNumber          bool   `json:"requiredDiscountCardNumber"`
	ReservePrepaidGoods                 bool   `json:"reservePrepaidGoods"`
	ReturnFromClosedShiftEnabled        bool   `json:"returnFromClosedShiftEnabled"`
	SellReserves                        bool   `json:"sellReserves"`
	SendMarksForCheck                   bool   `json:"sendMarksForCheck"`
	SendMarksToChestnyZnakOnCloud       bool   `json:"sendMarksToChestnyZnakOnCloud"`
	SyncAgents                          bool   `json:"syncAgents"`
	Shared                              bool   `json:"shared"`
	ShowBeerOnTap                       bool   `json:"showBeerOnTap"`
	// Только для чтения.
	AllowCustomPrice  bool `json:"allowCustomPrice"`
	AuthTokenAttached bool `json:"authTokenAttached"`
	OFDEnabled        bool `json:"ofdEnabled"`
	// Перечисления открыты.
	DefaultTaxSystem      string            `json:"defaultTaxSystem"`
	FiscalType            string            `json:"fiscalType"`
	MarkingSellingMode    string            `json:"markingSellingMode"`
	MarksCheckMode        string            `json:"marksCheckMode"`
	MinionToMasterType    string            `json:"minionToMasterType"`
	OrderTaxSystem        string            `json:"orderTaxSystem"`
	PriorityOFDSend       string            `json:"priorityOfdSend"`
	TobaccoMRCControlType string            `json:"tobaccoMrcControlType"`
	BankPercent           float64           `json:"bankPercent"`
	QRBankPercent         float64           `json:"qrBankPercent"`
	DiscountMaxPercent    float64           `json:"discountMaxPercent"`
	Organization          *model.Reference  `json:"organization,omitzero"`
	Store                 *model.Reference  `json:"store,omitzero"`
	PriceType             *model.PriceType  `json:"priceType,omitzero"`
	Acquire               *model.Reference  `json:"acquire,omitzero"`
	CreateOrderWithState  *model.Reference  `json:"createOrderWithState,omitzero"`
	OrderToState          *model.Reference  `json:"orderToState,omitzero"`
	OrganizationBranch    *model.Reference  `json:"organizationBranch,omitzero"`
	Owner                 *model.Reference  `json:"owner,omitzero"`
	Group                 *model.Reference  `json:"group,omitzero"`
	QRAcquire             *model.Reference  `json:"qrAcquire,omitzero"`
	ReceiptTemplate       *model.Reference  `json:"receiptTemplate,omitzero"`
	AddressFull           *model.Address    `json:"addressFull,omitzero"`
	CreateAgentsTags      []string          `json:"createAgentsTags"`
	FilterAgentsTags      []string          `json:"filterAgentsTags"`
	CustomerOrderStates   []model.Reference `json:"customerOrderStates"`
	// MasterRetailStores и ProductFolders пишутся массивом ссылок.
	MasterRetailStores *Links `json:"masterRetailStores,omitzero"`
	ProductFolders     *Links `json:"productFolders,omitzero"`
	// Cashiers — метаданные коллекции кассиров; ресурс кассиров в SDK v1 не входит.
	Cashiers *model.Collection[model.Reference] `json:"cashiers,omitzero"`
	// Environment и State — окружение и статус кассы, только для чтения; JSON как пришёл
	// (вложенные таблицы источника не типизируются: SDK их не пишет и не интерпретирует).
	Environment jsontext.Value `json:"environment,omitzero"`
	State       jsontext.Value `json:"state,omitzero"`
	// LastOperationNames устарело в источнике.
	LastOperationNames []LastOperation `json:"lastOperationNames"`
}

// Links — ссылки, которые фикстура отдаёт коллекцией {meta[, rows]}, а пример ответа Create —
// массивом. Meta заполнен для формы коллекции; запись возвращает ту же форму.
type Links struct {
	Meta *model.CollectionMeta
	Rows []model.Reference
}

func (l *Links) UnmarshalJSON(data []byte) error {
	if v := jsontext.Value(data); v.Kind() == '[' {
		*l = Links{}
		return json.Unmarshal(data, &l.Rows)
	}
	var collection model.Collection[model.Reference]
	if err := json.Unmarshal(data, &collection); err != nil {
		return err
	}
	*l = Links{Meta: &collection.Meta, Rows: collection.Rows}
	return nil
}

func (l Links) MarshalJSON() ([]byte, error) {
	if l.Meta == nil {
		return json.Marshal(l.Rows)
	}
	return json.Marshal(struct {
		Meta *model.CollectionMeta `json:"meta"`
		Rows []model.Reference     `json:"rows,omitempty"`
	}{l.Meta, l.Rows})
}

// LastOperation — последняя операция на точке: код сущности и номер.
type LastOperation struct {
	Entity string `json:"entity"`
	Name   string `json:"name"`
}
