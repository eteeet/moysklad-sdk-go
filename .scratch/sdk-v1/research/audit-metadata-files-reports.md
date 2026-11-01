# SDK v1: Audit, контексты, metadata, файлы и Turnover

Primary source: official api-remap-1.2-doc, revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`. Таблицы JSON markers в
[contracts.md](../contracts.md); HTTP fixtures SDK сохраняются неизменными.

## Audit

Подтверждённые маршруты:

- GET `audit`: контексты, limit default25, 1…100, offset default0, filter.
- GET `audit/{contextId}/events`: события контекста.
- GET `entity/{entityType}/{id}/audit`: события сущности.
- GET `audit/metadata/filters`: каталог фильтров.

Отдельный GET `audit/{id}` не документирован и в v1 не придумывается.
Контекст и событие — разные модели: у event не требуется искусственный top-level ID.
Events имеют meta/context/rows в примерах, но их собственная таблица limit/order/filter
не подтверждена; default25 в example не доказательство полной pagination API.
[Audit source](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md).

Фильтры source/entityType/eventType и UID следуют source; общий раздел также
документирует filter=application=<installation href>. Не путать App UID и installation ID.
Права/доступность событий проверяет сервер. Источник не даёт durable cursor,
exactly-once, snapshot order или гарантии неограниченной retention: SDK этого не обещает.
[App/Audit filters](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L3618).

Diff зависит от event type: значения могут быть вложенными, unknown, null, содержать
большие integer. Рекомендация — typed event header + lossless raw diff и осторожные
accessors, без привязки каждой старой записи к текущей read-модели entity. Агрегатор
реплики, durable checkpoint и автоматическое последующее чтение events вне SDK.
[Diff format](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L260).

## Контексты

GET `context/employee` не принимает caller ID, возвращает данные Employee и permissions.
Это не `model.ResponseContext` из collection wrapper. GET `context/application`
возвращает App installation meta/id/appUid по токену решения; отдельно существует
GET `entity/application/{installationId}`. Не смешивать их и полный Vendor API.
[Employee context](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L2898),
[Application context](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L3618).

## Metadata и States

GET `entity/metadata` получает общий каталог, filter=type=product;type=service
ограничивает сущности. GET `entity/{type}/metadata` — metadata данного типа;
attribute definitions и states имеют свои shapes и права. «Все metadata» не означает
автоматический CRUD всех определений дополнительных полей или шаблонов.
[General metadata](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L97).

States читаются в metadata, создаются POST `entity/{type}/metadata/states` и меняются
PUT `.../states/{id}`. Required Create: name/color/stateType; color — ARGB integer,
stateType Regular/Successful/Unsuccessful (unknown read values не отвергать).
Запись требует admin по источнику. Не выдумывать GET list states вместо metadata;
отдельный GET state добавлять только при подтверждённом endpoint.
[States source](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_states.md).

## Files и download

List/Upload — owner-scoped `entity/{type}/{id}/files`. List даёт metadata файлов;
POST принимает массив {filename,content Base64}, максимум10 за вызов. Это специальный
upload protocol, не разрешение generic batch всех сущностей. Поля metadata:
created/createdBy/filename/meta/miniature/size/tiny/title. Не выдумывать entity/file
CRUD или GetFile-by-ID, не подтверждённый source.
[Files source](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md).

Ограничения 100 files/object и 10 upload/request не требуют скрытого GET ради
подсчёта. Проверка encoded JSON body включает Base64 overhead; decoded file limit
и общий HTTP body limit не взаимозаменяемы. Signed download URL и expiration не
делают доступ к нему безопасным для логирования токенов/Authorization.
[Limits и размер](transport.md#размеры-files-и-streams),
[временные ссылки](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-ssylki-na-fajly).

Текущий transport.Do закрывает body и предназначен для JSON. Рекомендация — отдельный
stream или writer-based download с явным ownership, close/cancellation и shared
limiter lifetime. Доверенный API download может перенаправить на другое хранилище:
Authorization должен удаляться при смене origin; произвольный external URL не
принимается как trusted API endpoint. После partial output нельзя молча retry.
Images остаются отдельным набором методов/лимитов, а не alias Files.
[Workbook files](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/workbook/_workbook_files.md),
[redirect policy исходного SDK](../../../transport/client.go).

## Turnover: три разных результата

- GET `report/turnover/all`: rows assortment/onPeriodStart/onPeriodEnd/income/outcome;
  metrics sum в копейках и quantity; groupBy default product, variant поддержан.
  limit/offset документированы у All.
- GET `report/turnover/bystore`: filter product либо variant обязателен; rows имеют
  stockByStore. momentFrom/To optional; limit/offset/groupBy не наследуются от All.
- GET `report/turnover/byoperations`: filter product либо variant обязателен, по одному
  непустому значению каждого фильтра. Rows assortment/store/operation/quantity/cost/sum.
  Таблица называет momentFrom/To required, но соседние пояснения и пример допускают
  пропуск. Это открытый source conflict, не доказательство любой из двух трактовок.

[Turnover source](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md).

Каждый endpoint имеет собственный row/options type и success criteria; top-level
entity ID не нужен. Даты периода и report timezone header не создают глобальное
timezone состояние в JSON decoder. По source header — дата клиента с numeric zone
(пример Wed, 16 Aug 2017 23:07:01 +0700); обозначение «RFC 3522» выглядит опечаткой.
До уточнения ориентироваться на явно приведённый формат и report-specific semantics.
Рекомендация examples — передавать explicit range, source uncertainty не скрывать.
[Headers ByOperations](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L587).
