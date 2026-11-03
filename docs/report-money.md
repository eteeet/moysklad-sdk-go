# Деньги (report/money)

Пакет `report/money`; handle фасада планируется как `c.Report().Money()` (пока
`money.NewService(transportClient)`). Источник:
[_report_money.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_money.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод `money.Service` | Endpoint | Ответ |
| --- | --- | --- |
| `PlotSeries(ctx, PlotSeriesOptions, opts...)` | GET report/money/plotseries | `*PlotSeries` |
| `ByAccount(ctx, opts...)` | GET report/money/byaccount | `*model.Collection[AccountRow]` |

Нужны права `viewDashboard` и `viewMoneyDashboard` (проверяет сервер). Каждый метод — один GET;
повтор — только по `Config.Retry` клиента.

## Параметры `PlotSeriesOptions`

| Поле | Query | Ограничения |
| --- | --- | --- |
| `MomentFrom` | momentFrom | обязателен; нулевой — ошибка до запроса |
| `MomentTo` | momentTo | обязателен; нулевой — ошибка до запроса |
| `Interval` | interval | обязателен; `IntervalHour`, `IntervalDay`, `IntervalMonth`; набор открыт — другое непустое значение уходит как есть |
| `Filter` | filter | результат `query.Filter`: organization, project (ссылки, `=`, несколько значений; пустые нельзя — проверяет сервер) |
| `Extra` | любые | без повторов typed-полей, без expand и async |

Моменты уходят в московском времени (`query.ReportOptions`). Заголовок `X-Lognex-Accept-Timezone`
задаёт `transport.WithAcceptTimezone`; ответный `X-Lognex-Content-Timezone` — в
`ResponseInfo.ContentTimezone`. У `ByAccount` параметров источник не описывает.

## Поля (все только для чтения)

`PlotSeries`:

| JSON | Go type | Смысл |
| --- | --- | --- |
| context | `*model.ResponseContext` | сотрудник запроса; нет в таблице, есть в примерах |
| meta | `model.Meta` | нет в таблице, есть в примерах; не коллекция (нет size) |
| credit | `float64` | доход за период, копейки |
| debit | `float64` | расход за период, копейки |
| series | `[]Point` | показатели по интервалам; null — ошибка SDK |

`Point`: `date` (`model.DateTime`), `credit`, `debit`, `balance` (доход − расход) — `float64`, копейки.

`AccountRow` (строка `rows` ответа ByAccount):

| JSON | Go type | Смысл |
| --- | --- | --- |
| account | `*Named` | счёт организации (name — номер счёта); nil у остатка кассы |
| organization | `Named` | организация (name — наименование) |
| balance | `float64` | текущий остаток в валюте счёта, копейки (минимальные единицы); валюту SDK не пересчитывает |

## Противоречия и неизвестное

- Таблица атрибутов plotseries не перечисляет `context` и `meta`; они есть во всех примерах, SDK их
  читает и не требует.
- Ответ ByAccount — `{context, meta, rows}` без size/limit/offset: в `CollectionMeta` они нулевые,
  страниц нет.
- Описание `meta` счёта — «Метаданные организации»; пример содержит meta счёта (`type: account`).
- Пример запроса передаёт даты без времени (`momentFrom=2018-09-01`); SDK всегда передаёт
  `yyyy-MM-dd HH:mm:ss`, как в параметрах источника.
- Async-режим (I23) не входит.

## Ошибки

Один префикс `moysklad:` и операция `get money plotseries` / `get money byaccount`. Ошибки API —
`*errors.APIError`; статус, отличный от 200, ответ PlotSeries без `series` и ByAccount без
`meta.href` или с `rows` null — ошибка.

Фикстуры: `testdata/custom/money` (примеры ответов источника с SHA-256).
