package retailstore

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы непустые name, organization, store и priceType. Без active сервер
// включает точку, если позволяет тариф. allowCreateProducts и controlShippingStock не могут быть
// оба true, у основного счёта организации должна быть валюта учёта — проверяет сервер.
// externalCode, environment, state, allowCustomPrice, authTokenAttached, ofdEnabled таблица
// помечает только для чтения — их здесь нет, хотя пример Create передаёт allowCustomPrice и
// authTokenAttached. Cashiers (тоже в примере) не пишется: администрирование кассиров в v1 не входит. createAgentsTags таблица называет
// Array(Object), описание и примеры — массивом строк; SDK пишет строки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-sozdat-tochku-prodazh
type CreateRequest struct {
	Name                                string                            `json:"name"`
	Organization                        model.Reference                   `json:"organization"`
	Store                               model.Reference                   `json:"store"`
	PriceType                           model.PriceType                   `json:"priceType"`
	Address                             model.Field[string]               `json:"address,omitzero"`
	DemandPrefix                        model.Field[string]               `json:"demandPrefix,omitzero"`
	Description                         model.Field[string]               `json:"description,omitzero"`
	IDQR                                model.Field[string]               `json:"idQR,omitzero"`
	QRTerminalID                        model.Field[string]               `json:"qrTerminalId,omitzero"`
	Active                              model.Field[bool]                 `json:"active,omitzero"`
	AllowCreateProducts                 model.Field[bool]                 `json:"allowCreateProducts,omitzero"`
	AllowDeleteReceiptPositions         model.Field[bool]                 `json:"allowDeleteReceiptPositions,omitzero"`
	AllowSellTobaccoWithoutMRC          model.Field[bool]                 `json:"allowSellTobaccoWithoutMRC,omitzero"`
	Archived                            model.Field[bool]                 `json:"archived,omitzero"`
	ControlCashierChoice                model.Field[bool]                 `json:"controlCashierChoice,omitzero"`
	ControlShippingStock                model.Field[bool]                 `json:"controlShippingStock,omitzero"`
	CreateCashInOnRetailShiftClosing    model.Field[bool]                 `json:"createCashInOnRetailShiftClosing,omitzero"`
	CreatePaymentInOnRetailShiftClosing model.Field[bool]                 `json:"createPaymentInOnRetailShiftClosing,omitzero"`
	DiscountEnable                      model.Field[bool]                 `json:"discountEnable,omitzero"`
	EnableReturnsWithNoReason           model.Field[bool]                 `json:"enableReturnsWithNoReason,omitzero"`
	IssueOrders                         model.Field[bool]                 `json:"issueOrders,omitzero"`
	OnlyInStock                         model.Field[bool]                 `json:"onlyInStock,omitzero"`
	PrintAlways                         model.Field[bool]                 `json:"printAlways,omitzero"`
	QRPayEnabled                        model.Field[bool]                 `json:"qrPayEnabled,omitzero"`
	RequiredFio                         model.Field[bool]                 `json:"requiredFio,omitzero"`
	RequiredPhone                       model.Field[bool]                 `json:"requiredPhone,omitzero"`
	RequiredEmail                       model.Field[bool]                 `json:"requiredEmail,omitzero"`
	RequiredBirthdate                   model.Field[bool]                 `json:"requiredBirthdate,omitzero"`
	RequiredSex                         model.Field[bool]                 `json:"requiredSex,omitzero"`
	RequiredDiscountCardNumber          model.Field[bool]                 `json:"requiredDiscountCardNumber,omitzero"`
	ReservePrepaidGoods                 model.Field[bool]                 `json:"reservePrepaidGoods,omitzero"`
	ReturnFromClosedShiftEnabled        model.Field[bool]                 `json:"returnFromClosedShiftEnabled,omitzero"`
	SellReserves                        model.Field[bool]                 `json:"sellReserves,omitzero"`
	SendMarksForCheck                   model.Field[bool]                 `json:"sendMarksForCheck,omitzero"`
	SendMarksToChestnyZnakOnCloud       model.Field[bool]                 `json:"sendMarksToChestnyZnakOnCloud,omitzero"`
	SyncAgents                          model.Field[bool]                 `json:"syncAgents,omitzero"`
	Shared                              model.Field[bool]                 `json:"shared,omitzero"`
	ShowBeerOnTap                       model.Field[bool]                 `json:"showBeerOnTap,omitzero"`
	DefaultTaxSystem                    model.Field[string]               `json:"defaultTaxSystem,omitzero"`
	FiscalType                          model.Field[string]               `json:"fiscalType,omitzero"`
	MarkingSellingMode                  model.Field[string]               `json:"markingSellingMode,omitzero"`
	MarksCheckMode                      model.Field[string]               `json:"marksCheckMode,omitzero"`
	MinionToMasterType                  model.Field[string]               `json:"minionToMasterType,omitzero"`
	OrderTaxSystem                      model.Field[string]               `json:"orderTaxSystem,omitzero"`
	PriorityOFDSend                     model.Field[string]               `json:"priorityOfdSend,omitzero"`
	TobaccoMRCControlType               model.Field[string]               `json:"tobaccoMrcControlType,omitzero"`
	BankPercent                         model.Field[float64]              `json:"bankPercent,omitzero"`
	QRBankPercent                       model.Field[float64]              `json:"qrBankPercent,omitzero"`
	DiscountMaxPercent                  model.Field[float64]              `json:"discountMaxPercent,omitzero"`
	Acquire                             model.Field[model.Reference]      `json:"acquire,omitzero"`
	CreateOrderWithState                model.Field[model.Reference]      `json:"createOrderWithState,omitzero"`
	OrderToState                        model.Field[model.Reference]      `json:"orderToState,omitzero"`
	OrganizationBranch                  model.Field[model.Reference]      `json:"organizationBranch,omitzero"`
	Owner                               model.Field[model.Reference]      `json:"owner,omitzero"`
	Group                               model.Field[model.Reference]      `json:"group,omitzero"`
	QRAcquire                           model.Field[model.Reference]      `json:"qrAcquire,omitzero"`
	ReceiptTemplate                     model.Field[model.Reference]      `json:"receiptTemplate,omitzero"`
	AddressFull                         model.Field[model.AddressRequest] `json:"addressFull,omitzero"`
	CreateAgentsTags                    model.Field[[]string]             `json:"createAgentsTags,omitzero"`
	FilterAgentsTags                    model.Field[[]string]             `json:"filterAgentsTags,omitzero"`
	CustomerOrderStates                 model.Field[[]model.Reference]    `json:"customerOrderStates,omitzero"`
	MasterRetailStores                  model.Field[[]model.Reference]    `json:"masterRetailStores,omitzero"`
	ProductFolders                      model.Field[[]model.Reference]    `json:"productFolders,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-izmenit-tochku-prodazh
type UpdateRequest struct {
	Name                                model.Field[string]               `json:"name,omitzero"`
	Organization                        model.Field[model.Reference]      `json:"organization,omitzero"`
	Store                               model.Field[model.Reference]      `json:"store,omitzero"`
	PriceType                           model.Field[model.PriceType]      `json:"priceType,omitzero"`
	Address                             model.Field[string]               `json:"address,omitzero"`
	DemandPrefix                        model.Field[string]               `json:"demandPrefix,omitzero"`
	Description                         model.Field[string]               `json:"description,omitzero"`
	IDQR                                model.Field[string]               `json:"idQR,omitzero"`
	QRTerminalID                        model.Field[string]               `json:"qrTerminalId,omitzero"`
	Active                              model.Field[bool]                 `json:"active,omitzero"`
	AllowCreateProducts                 model.Field[bool]                 `json:"allowCreateProducts,omitzero"`
	AllowDeleteReceiptPositions         model.Field[bool]                 `json:"allowDeleteReceiptPositions,omitzero"`
	AllowSellTobaccoWithoutMRC          model.Field[bool]                 `json:"allowSellTobaccoWithoutMRC,omitzero"`
	Archived                            model.Field[bool]                 `json:"archived,omitzero"`
	ControlCashierChoice                model.Field[bool]                 `json:"controlCashierChoice,omitzero"`
	ControlShippingStock                model.Field[bool]                 `json:"controlShippingStock,omitzero"`
	CreateCashInOnRetailShiftClosing    model.Field[bool]                 `json:"createCashInOnRetailShiftClosing,omitzero"`
	CreatePaymentInOnRetailShiftClosing model.Field[bool]                 `json:"createPaymentInOnRetailShiftClosing,omitzero"`
	DiscountEnable                      model.Field[bool]                 `json:"discountEnable,omitzero"`
	EnableReturnsWithNoReason           model.Field[bool]                 `json:"enableReturnsWithNoReason,omitzero"`
	IssueOrders                         model.Field[bool]                 `json:"issueOrders,omitzero"`
	OnlyInStock                         model.Field[bool]                 `json:"onlyInStock,omitzero"`
	PrintAlways                         model.Field[bool]                 `json:"printAlways,omitzero"`
	QRPayEnabled                        model.Field[bool]                 `json:"qrPayEnabled,omitzero"`
	RequiredFio                         model.Field[bool]                 `json:"requiredFio,omitzero"`
	RequiredPhone                       model.Field[bool]                 `json:"requiredPhone,omitzero"`
	RequiredEmail                       model.Field[bool]                 `json:"requiredEmail,omitzero"`
	RequiredBirthdate                   model.Field[bool]                 `json:"requiredBirthdate,omitzero"`
	RequiredSex                         model.Field[bool]                 `json:"requiredSex,omitzero"`
	RequiredDiscountCardNumber          model.Field[bool]                 `json:"requiredDiscountCardNumber,omitzero"`
	ReservePrepaidGoods                 model.Field[bool]                 `json:"reservePrepaidGoods,omitzero"`
	ReturnFromClosedShiftEnabled        model.Field[bool]                 `json:"returnFromClosedShiftEnabled,omitzero"`
	SellReserves                        model.Field[bool]                 `json:"sellReserves,omitzero"`
	SendMarksForCheck                   model.Field[bool]                 `json:"sendMarksForCheck,omitzero"`
	SendMarksToChestnyZnakOnCloud       model.Field[bool]                 `json:"sendMarksToChestnyZnakOnCloud,omitzero"`
	SyncAgents                          model.Field[bool]                 `json:"syncAgents,omitzero"`
	Shared                              model.Field[bool]                 `json:"shared,omitzero"`
	ShowBeerOnTap                       model.Field[bool]                 `json:"showBeerOnTap,omitzero"`
	DefaultTaxSystem                    model.Field[string]               `json:"defaultTaxSystem,omitzero"`
	FiscalType                          model.Field[string]               `json:"fiscalType,omitzero"`
	MarkingSellingMode                  model.Field[string]               `json:"markingSellingMode,omitzero"`
	MarksCheckMode                      model.Field[string]               `json:"marksCheckMode,omitzero"`
	MinionToMasterType                  model.Field[string]               `json:"minionToMasterType,omitzero"`
	OrderTaxSystem                      model.Field[string]               `json:"orderTaxSystem,omitzero"`
	PriorityOFDSend                     model.Field[string]               `json:"priorityOfdSend,omitzero"`
	TobaccoMRCControlType               model.Field[string]               `json:"tobaccoMrcControlType,omitzero"`
	BankPercent                         model.Field[float64]              `json:"bankPercent,omitzero"`
	QRBankPercent                       model.Field[float64]              `json:"qrBankPercent,omitzero"`
	DiscountMaxPercent                  model.Field[float64]              `json:"discountMaxPercent,omitzero"`
	Acquire                             model.Field[model.Reference]      `json:"acquire,omitzero"`
	CreateOrderWithState                model.Field[model.Reference]      `json:"createOrderWithState,omitzero"`
	OrderToState                        model.Field[model.Reference]      `json:"orderToState,omitzero"`
	OrganizationBranch                  model.Field[model.Reference]      `json:"organizationBranch,omitzero"`
	Owner                               model.Field[model.Reference]      `json:"owner,omitzero"`
	Group                               model.Field[model.Reference]      `json:"group,omitzero"`
	QRAcquire                           model.Field[model.Reference]      `json:"qrAcquire,omitzero"`
	ReceiptTemplate                     model.Field[model.Reference]      `json:"receiptTemplate,omitzero"`
	AddressFull                         model.Field[model.AddressRequest] `json:"addressFull,omitzero"`
	CreateAgentsTags                    model.Field[[]string]             `json:"createAgentsTags,omitzero"`
	FilterAgentsTags                    model.Field[[]string]             `json:"filterAgentsTags,omitzero"`
	CustomerOrderStates                 model.Field[[]model.Reference]    `json:"customerOrderStates,omitzero"`
	MasterRetailStores                  model.Field[[]model.Reference]    `json:"masterRetailStores,omitzero"`
	ProductFolders                      model.Field[[]model.Reference]    `json:"productFolders,omitzero"`
}
