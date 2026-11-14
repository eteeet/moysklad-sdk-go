# Показатели (report/dashboard)

Пакет `report/dashboard`, handle фасада — `c.Report().Dashboard()`. Источник:
[_dashboard.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_dashboard.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод `dashboard.Service` | Endpoint | Ответ |
| --- | --- | --- |
| `Day(ctx, opts...)` | GET report/dashboard/day | `*Dashboard` |
| `Week(ctx, opts...)` | GET report/dashboard/week | `*Dashboard` |
| `Month(ctx, opts...)` | GET report/dashboard/month | `*Dashboard` |

Параметров запроса источник не описывает. Ответ — один объект без `meta`; каждый метод — один
GET. Нужно право `viewDashboard` (проверяет сервер).

## Поля (все только для чтения, все обязательны в ответе)

| JSON | Go type | Тип в источнике | Смысл |
| --- | --- | --- | --- |
| sales / orders | `Sales` / `Orders` | Object | раздел обязателен; null или отсутствие — ошибка SDK |
| …count | `int` | Int | количество |
| …amount | `float64` | Int | сумма, копейки (источник: «Прибыль») |
| …movementAmount | `float64` | Int | изменение к прошлому такому же периоду, копейки |
| money | `Money` | Object | раздел обязателен |
| money.income | `float64` | Int | доходы, копейки |
| money.outcome | `float64` | Float | расходы, копейки |
| money.balance | `float64` | Float | текущий баланс, копейки |
| money.todayMovement | `float64` | Float | дельта за сегодня; не объявлена деньгами |
| money.movement | `float64` | Float | дельта за период; не объявлена деньгами |

## Асинхронно

«Отчет Показатели» входит в список async источника целиком. [_async.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md) (revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод | Endpoint |
| --- | --- |
| `DayAsync(ctx)` | GET report/dashboard/day?async=true |
| `WeekAsync(ctx)` | GET report/dashboard/week?async=true |
| `MonthAsync(ctx)` | GET report/dashboard/month?async=true |

Каждый метод — один GET с теми же параметрами, что у синхронного, плюс `async=true`; ответ 202 с
Content-Location (…/async/{id}) и Location (…/async/{id}/result) — `*async.Submitted` (ID,
StatusHref, ResultHref), другой статус или 202 без этих адресов — ошибка. Создание задачи —
побочный эффект: transport не повторяет `async=true` даже с `Config.Retry`; после сбоя сети задача
могла создаться — проверьте `c.Async().List`. Результат: `c.Async().Wait`, затем
`c.Async().Result(ctx, task.ID, w)` — JSON того же вида, что ответ синхронного метода (SDK его не
разбирает). `limit`/`offset` источник у async запрещает: параметров у методов нет, в запросе только `async=true`.

## Противоречия и неизвестное

- Ticket I06: «sales/orders/money — копейки, movement — проценты». Источник единиц не называет.
  Для `movementAmount` пример противоречит процентам: `amount` 0 и `movementAmount` -80000 —
  это абсолютное изменение суммы (было 80000 копеек), поэтому SDK документирует его в копейках.
  `money.movement` и `money.todayMovement` в примере равны 0; по ticket SDK не объявляет их
  деньгами и отдаёт число как пришло. Тип `float64` одинаков при любой трактовке.
- Таблица «Заказы за период» повторяет описания продаж («Количество продаж», «Прибыль»);
  SDK считает это опечаткой источника.
- Тип Int у сумм противоречит правилу «деньги — float64»; SDK читает `float64`, целые значения
  читаются без потерь.

## Ошибки

Ошибки содержат один префикс `moysklad:` и операцию `get dashboard day|week|month`; ошибки
API — `*errors.APIError`; статус, отличный от 200, — ошибка.
