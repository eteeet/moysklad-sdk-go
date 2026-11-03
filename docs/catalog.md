# Товары, услуги, комплекты, модификации и цены

Источник контрактов — `api-remap-1.2-doc` revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0001](adr/0001-request-options.md)).
Деньги в ценах — копейки (float64), ставки и проценты — нет. Ссылки — `model.Reference` без expand.

## Справочники

| Пакет / метод | Маршрут |
| --- | --- |
| `entity/currency`: `Get`, `List`, `Create`, `Update` | `entity/currency[/{id}]` |
| `entity/pricetype`: `List`, `Get`, `GetDefault`, `ReplaceAll` | `context/companysettings/pricetype[/{id},/default]` |
| `entity/taxrate`: `Get`, `List`, `Create`, `Update` | `entity/taxrate[/{id}]` |
| `entity/service`: `Get`, `List`, `Create`, `Update` (клиент `service.Client`) | `entity/service[/{id}]` |
| `entity/bundle`: `Get`, `List`, `Create`, `Update`, `GetComponents`, `GetComponent` | `entity/bundle[/{id}[/components[/{id}]]]` |
| `entity/variant`: `Get`, `List`, `Create`, `Update` | `entity/variant[/{id}]` |
| `variant.Service`: `ListCharacteristics`, `GetCharacteristic`, `CreateCharacteristic` | `entity/variant/metadata[/characteristics[/{id}]]` |
| `entity/assortment`: `List`, `ListLegacy` (Deprecated) | `entity/assortment` |

- Currency: системную валюту создают `System: true` с code или isoCode; `rateUpdateType` пишется только
  при Create (исключение из readonly-маркера источника). rate 0 и пустые name/code/isoCode отклоняются.
- PriceType: `ReplaceAll` заменяет весь список (1–100 типов); существующие передаются через
  `pricetype.Existing(p)` с meta, иначе сервер создаст новый тип.
- Service, Bundle, Variant: штрихкод — `model.Barcode` (объект с одним форматом); доп. поля и цены
  продажи сервер сливает по метаданным и типу цены. `mod__tasnif__uz` и `mod__marking__uz` — только
  с `transport.WithAcceptRegion(transport.RegionUZ)`.
- Штрихкоды при Update: у Bundle и Variant список заменяется целиком, пустой удаляет все (правила
  создания у Variant скопированы источником из раздела комплекта и не подтверждены). У Service
  источник замену не описывает: что делает пустой или неполный список, неизвестно.
- Bundle: Create требует name и непустые components; `ComponentRequest{ID, Assortment, Quantity}` —
  с ID изменяет существующий компонент. Quantity пишется, хотя таблица источника помечает его
  readonly: исключение из правила «readonly сильнее примеров» утверждено пользователем — без
  количества компонент не имеет смысла, его передают примеры Create и добавления компонента. Компоненты — свой тип, не Position;
  мутации компонентов — [dictionaries.md](dictionaries.md). `GetComponents` читает одну страницу (`query.PageOptions`).
- Variant: Create требует product и непустые characteristics. Значение характеристики
  (`CharacteristicValue`) сервер связывает с описанием по id, иначе по name; Update заменяет набор
  характеристик целиком. Описания характеристик — отдельные методы, скрыто не создаются.
- Null у Field отправляется как есть: семантику null источник для этих сущностей не описывает.

## Ассортимент

- `List` — новый endpoint, всегда с `X-Lognex-Remap-Beta-Feature: assortmentWithoutStock` (поверх
  опций вызова). Остатков в строках нет, SDK их не догружает: используйте отчёт об остатках по id.
  Все страницы — `All(ctx, options, opts...)`: List и обход по nextHref с тем же beta-тегом на каждом
  запросе. Прямой `pages.All` по странице List без `transport.WithBetaFeature(assortment.BetaFeature)`
  уводит GET по nextHref в legacy endpoint.
- `ListLegacy` — старый endpoint с stock, reserve, inTransit, quantity (`Row.Balances`); источник
  выводит его из эксплуатации в 2027 году.
- `assortment.Options{ListOptions, GroupBy}`; фильтры, умолчания и недоступные параметры каждого
  режима перечислены в godoc. fields не принимается.
- `Row` — union по meta.type: Product, Service, Bundle, Variant, Consignment; неизвестный тип — только
  `Type` и `Raw`. Запись Row возвращает `Raw` байт в байт.

## Не входит

Expand, изображения и файлы (только метаданные коллекций), `minimumStock` у Product и
Variant, склады остатков (storebalances). Мутации компонентов комплекта и массовое создание
характеристик — [dictionaries.md](dictionaries.md). Delete, массовые операции и `ListAsync` ассортимента —
[write-async.md](write-async.md). Противоречия источника записаны в godoc и Answer tickets C01–C08.
