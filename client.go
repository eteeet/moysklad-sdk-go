package moysklad

import (
	"github.com/eteeet/moysklad-sdk-go/async"
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
	"github.com/eteeet/moysklad-sdk-go/entity/consignment"
	"github.com/eteeet/moysklad-sdk-go/entity/contentcard"
	"github.com/eteeet/moysklad-sdk-go/entity/contract"
	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/entity/currency"
	"github.com/eteeet/moysklad-sdk-go/entity/customentity"
	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/entity/group"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/entity/organizationbranch"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/entity/personaldiscount"
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
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/images"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/security"
	"github.com/eteeet/moysklad-sdk-go/trackingcode"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Config — настройки общего transport.Client.
type Config = transport.Config

// Client — фасад SDK: группы повторяют путь API (Entity — entity/*, Report — report/*,
// Context — context/*). Все handles работают через один transport.Client: общие HTTPClient, Auth,
// Limiter и Retry. Handle не хранит собственного состояния и создаётся при каждом вызове без HTTP,
// поэтому два вызова дают разные указатели с одинаковым поведением. Client создаётся только через
// NewClient: у нулевого значения нет transport, и получение handle паникует.
type Client struct {
	transport *transport.Client
}

// NewClient не делает HTTP и не получает токен: Config.Auth используется как есть.
func NewClient(cfg Config) (*Client, error) {
	t, err := transport.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{transport: t}, nil
}

// Transport — общий transport.Client для Download, HeaderValues и standalone NewService.
func (c *Client) Transport() *transport.Client { return c.transport }

// Entity — сервисы entity/*.
func (c *Client) Entity() EntityGroup { return EntityGroup{c.transport} }

// Report — отчёты report/*.
func (c *Client) Report() ReportGroup { return ReportGroup{c.transport} }

// Context — context/employee и context/application.
func (c *Client) Context() *contexts.Service { return must(contexts.NewService(c.transport)) }

// Audit — контексты и события аудита.
func (c *Client) Audit() *audit.Service { return must(audit.NewService(c.transport)) }

// Async — асинхронные задачи: статус, список, ожидание, результат и отмена. Создают задачи
// ListAsync сервисов Entity().Counterparty() и Entity().Assortment().
func (c *Client) Async() *async.Service { return must(async.NewService(c.transport)) }

// Security — выпуск токена POST security/token.
func (c *Client) Security() *security.Service { return must(security.NewService(c.transport)) }

// EntityGroup — сервисы entity/*; получается через Client.Entity.
type EntityGroup struct{ transport *transport.Client }

// ReportGroup — отчёты report/*; получается через Client.Report.
type ReportGroup struct{ transport *transport.Client }

// Turnover — отчёт «Обороты».
func (r ReportGroup) Turnover() *turnover.Service { return must(turnover.NewService(r.transport)) }

// Metadata — метаданные сущностей (entity/metadata, entity/<type>/metadata) и статусы.
func (e EntityGroup) Metadata() *metadata.Service { return must(metadata.NewService(e.transport)) }

// Files — файлы операций, документов и других сущностей (entity/<type>/<id>/files).
func (e EntityGroup) Files() *files.Service { return must(files.NewService(e.transport)) }

// Images — изображения товара, комплекта, модификации и карточки контента.
func (e EntityGroup) Images() *images.Service { return must(images.NewService(e.transport)) }

// TrackingCode — коды маркировки позиций документов.
func (e EntityGroup) TrackingCode() *trackingcode.Service {
	return must(trackingcode.NewService(e.transport))
}

// Assortment — Ассортимент: товары, услуги, комплекты, модификации и серии (только чтение).
func (e EntityGroup) Assortment() *assortment.Service {
	return must(assortment.NewService(e.transport))
}

// BonusProgram — Бонусные программы.
func (e EntityGroup) BonusProgram() *bonusprogram.Service {
	return must(bonusprogram.NewService(e.transport))
}

// BonusTransaction — Бонусные операции.
func (e EntityGroup) BonusTransaction() *bonustransaction.Service {
	return must(bonustransaction.NewService(e.transport))
}

// Bundle — Комплекты.
func (e EntityGroup) Bundle() *bundle.Service { return must(bundle.NewService(e.transport)) }

// CashIn — Приходные ордера.
func (e EntityGroup) CashIn() *cashin.Service { return must(cashin.NewService(e.transport)) }

// CashOut — Расходные ордера.
func (e EntityGroup) CashOut() *cashout.Service { return must(cashout.NewService(e.transport)) }

// Counterparty — Контрагенты.
func (e EntityGroup) Counterparty() *counterparty.Service {
	return must(counterparty.NewService(e.transport))
}

// Country — Страны.
func (e EntityGroup) Country() *country.Service { return must(country.NewService(e.transport)) }

// Currency — Валюты.
func (e EntityGroup) Currency() *currency.Service { return must(currency.NewService(e.transport)) }

// CustomerOrder — Заказы покупателей.
func (e EntityGroup) CustomerOrder() *customerorder.Service {
	return must(customerorder.NewService(e.transport))
}

// Demand — Отгрузки.
func (e EntityGroup) Demand() *demand.Service { return must(demand.NewService(e.transport)) }

// Employee — Сотрудники и их права.
func (e EntityGroup) Employee() *employee.Service { return must(employee.NewService(e.transport)) }

// Enter — Оприходования.
func (e EntityGroup) Enter() *enter.Service { return must(enter.NewService(e.transport)) }

// ExpenseItem — Статьи расходов.
func (e EntityGroup) ExpenseItem() *expenseitem.Service {
	return must(expenseitem.NewService(e.transport))
}

// InvoiceIn — Счета поставщиков.
func (e EntityGroup) InvoiceIn() *invoicein.Service { return must(invoicein.NewService(e.transport)) }

// InvoiceOut — Счета покупателям.
func (e EntityGroup) InvoiceOut() *invoiceout.Service {
	return must(invoiceout.NewService(e.transport))
}

// Loss — Списания.
func (e EntityGroup) Loss() *loss.Service { return must(loss.NewService(e.transport)) }

// Move — Перемещения.
func (e EntityGroup) Move() *move.Service { return must(move.NewService(e.transport)) }

// Organization — Юрлица.
func (e EntityGroup) Organization() *organization.Service {
	return must(organization.NewService(e.transport))
}

// PaymentIn — Входящие платежи.
func (e EntityGroup) PaymentIn() *paymentin.Service { return must(paymentin.NewService(e.transport)) }

// PaymentOut — Исходящие платежи.
func (e EntityGroup) PaymentOut() *paymentout.Service {
	return must(paymentout.NewService(e.transport))
}

// PriceType — Типы цен.
func (e EntityGroup) PriceType() *pricetype.Service { return must(pricetype.NewService(e.transport)) }

// Product — Товары.
func (e EntityGroup) Product() *product.Service { return must(product.NewService(e.transport)) }

// PurchaseOrder — Заказы поставщикам.
func (e EntityGroup) PurchaseOrder() *purchaseorder.Service {
	return must(purchaseorder.NewService(e.transport))
}

// PurchaseReturn — Возвраты поставщикам.
func (e EntityGroup) PurchaseReturn() *purchasereturn.Service {
	return must(purchasereturn.NewService(e.transport))
}

// RetailDemand — Розничные продажи.
func (e EntityGroup) RetailDemand() *retaildemand.Service {
	return must(retaildemand.NewService(e.transport))
}

// RetailSalesReturn — Розничные возвраты.
func (e EntityGroup) RetailSalesReturn() *retailsalesreturn.Service {
	return must(retailsalesreturn.NewService(e.transport))
}

// RetailShift — Розничные смены.
func (e EntityGroup) RetailShift() *retailshift.Service {
	return must(retailshift.NewService(e.transport))
}

// RetailStore — Точки продаж.
func (e EntityGroup) RetailStore() *retailstore.Service {
	return must(retailstore.NewService(e.transport))
}

// Role — Роли.
func (e EntityGroup) Role() *role.Service { return must(role.NewService(e.transport)) }

// SalesReturn — Возвраты покупателей.
func (e EntityGroup) SalesReturn() *salesreturn.Service {
	return must(salesreturn.NewService(e.transport))
}

// Service — Услуги.
func (e EntityGroup) Service() *service.Client { return must(service.NewService(e.transport)) }

// Supply — Приёмки.
func (e EntityGroup) Supply() *supply.Service { return must(supply.NewService(e.transport)) }

// TaxRate — Ставки НДС.
func (e EntityGroup) TaxRate() *taxrate.Service { return must(taxrate.NewService(e.transport)) }

// UOM — Единицы измерения.
func (e EntityGroup) UOM() *uom.Service { return must(uom.NewService(e.transport)) }

// AggregatePack — Транспортные упаковки.
func (e EntityGroup) AggregatePack() *aggregatepack.Service {
	return must(aggregatepack.NewService(e.transport))
}

// SpecialPriceDiscount — Специальные цены.
func (e EntityGroup) SpecialPriceDiscount() *specialpricediscount.Service {
	return must(specialpricediscount.NewService(e.transport))
}

// PersonalDiscount — Персональные скидки.
func (e EntityGroup) PersonalDiscount() *personaldiscount.Service {
	return must(personaldiscount.NewService(e.transport))
}

// AccumulationDiscount — Накопительные скидки.
func (e EntityGroup) AccumulationDiscount() *accumulationdiscount.Service {
	return must(accumulationdiscount.NewService(e.transport))
}

// Discount — Скидки: общий список и округление копеек.
func (e EntityGroup) Discount() *discount.Service { return must(discount.NewService(e.transport)) }

// ContentCard — Карточки контента.
func (e EntityGroup) ContentCard() *contentcard.Service {
	return must(contentcard.NewService(e.transport))
}

// Consignment — Партии.
func (e EntityGroup) Consignment() *consignment.Service {
	return must(consignment.NewService(e.transport))
}

// CustomEntity — Пользовательские справочники и их элементы.
func (e EntityGroup) CustomEntity() *customentity.Service {
	return must(customentity.NewService(e.transport))
}

// Thing — Серийные номера (только чтение).
func (e EntityGroup) Thing() *thing.Service { return must(thing.NewService(e.transport)) }

// OrganizationBranch — Подразделения юрлиц.
func (e EntityGroup) OrganizationBranch() *organizationbranch.Service {
	return must(organizationbranch.NewService(e.transport))
}

// SalePlatform — Площадки для продаж (только чтение).
func (e EntityGroup) SalePlatform() *saleplatform.Service {
	return must(saleplatform.NewService(e.transport))
}

// SalesChannel — Каналы продаж.
func (e EntityGroup) SalesChannel() *saleschannel.Service {
	return must(saleschannel.NewService(e.transport))
}

// Region — Регионы (только чтение).
func (e EntityGroup) Region() *region.Service { return must(region.NewService(e.transport)) }

// Contract — Договоры.
func (e EntityGroup) Contract() *contract.Service { return must(contract.NewService(e.transport)) }

// Project — Проекты.
func (e EntityGroup) Project() *project.Service { return must(project.NewService(e.transport)) }

// Group — Отделы.
func (e EntityGroup) Group() *group.Service { return must(group.NewService(e.transport)) }

// ProductFolder — Группы товаров.
func (e EntityGroup) ProductFolder() *productfolder.Service {
	return must(productfolder.NewService(e.transport))
}

// Store — Склады.
func (e EntityGroup) Store() *store.Service { return must(store.NewService(e.transport)) }

// Variant — Модификации и характеристики.
func (e EntityGroup) Variant() *variant.Service { return must(variant.NewService(e.transport)) }

// must: NewService отказывает только при nil client — это нулевое значение Client.
func must[S any](s *S, err error) *S {
	if err != nil {
		panic(err)
	}
	return s
}
