# Здоровье документов

Скан 2026-11-11, `dev` @ `c7f33d0`. Только регистрация, ничего не исправлено.
Источник правды: `~/ETEEET/AGENTS.md`, `~/ETEEET/GIT_FLOW.md`, `~/ETEEET/INFRA.md`, `~/ETEEET/moysklad/AGENTS.md`, `AGENTS.md` и `docs/agents/` → ADR → спеки и тикеты → docs.
Формат находки: `файл:строка — что с чем — предложение`.

## Находки

### A. Противоречия источнику правды (AGENTS, GIT_FLOW, INFRA)

- A1. репозиторий: нет `.github/workflows/`; решение «CI — только `make check`» в `.scratch/sdk-v1-release/issues/03-ci-release.md:43,86`, `.scratch/sdk-v1/map.md:157`, `docs/status-technical.md:104`, `docs/status-backlog.md:162`, `docs/coverage.md:484` — с `INFRA.md:15-17` (CI на GitHub Actions у каждого репозитория с кодом) — тикет на workflow по INFRA после «да» владельца (`INFRA.md:29`) либо ADR-исключение.
- A2. теги `v0.1.0`…`v0.2.2` — с `GIT_FLOW.md:19` (тег — на коммит с зелёным CI) — следствие A1, закрывается вместе с ним.
- A3. `origin/main` = `5f3fb93 Initial commit`, последний выпуск `v0.2.2` — с `GIT_FLOW.md:20` (владелец перематывает `main` на тег) — действие владельца: `git push origin 'v0.2.2^{}:main'`.
- A4. «Planning/tracker не коммитить без отдельной команды»: `.scratch/sdk-full-dictionaries/spec.md:347`, `sdk-full-catalog/spec.md:193`, `sdk-full-documents/spec.md:270`, `sdk-full-production/spec.md:138`, `sdk-full-reports/spec.md:303`, шаблон в 143 тикетах `sdk-full-*/issues/*` (строка 39, в т.ч. открытый `sdk-full-documents/issues/07-cash-accounts.md:39`), `.scratch/review-followups/spec.md:28-29`, `.scratch/sdk-v1/map.md:11`, `docs/temp_plan.md:61,136-137` — с `GIT_FLOW.md:12` (документы коммитятся прямо в `dev`) — убрать из открытых; закрытые не трогать (история).
- A5. «Live API — только с согласия / не входит»: `docs/status-backlog.md:160-161`, `docs/status-technical.md:100`, `docs/coverage.md:9-10,367,401`, `docs/temp_plan.md:135`, `.scratch/sdk-full-documents/issues/07-cash-accounts.md:18`, `.scratch/sdk-full-documents/spec.md:138`, `.scratch/sdk-full-{dictionaries:348,catalog:194,documents:271,production:139,reports:304}/spec.md` — с `~/ETEEET/moysklad/AGENTS.md:16` (тестовый аккаунт `admin@documan` — работать без запроса) — переписать: live-проверки на тестовом аккаунте разрешены.
- A6. Create требует `Name != ""`: `entity/invoiceout/service.go:43`, `entity/retailsalesreturn/service.go:45`, `entity/processingorder/service.go:42`; описано в `docs/documents.md:16-18`, `.scratch/sdk-v1-retail/spec.md:88`, `.scratch/invoice-documents/spec.md:68-70,88` — с `~/ETEEET/moysklad/AGENTS.md:18` (документам не задавать `name`) — тикет: проверить автонумерацию на `admin@documan` и снять проверку, как сделано для invoicein в v0.2.2.
- A7. `.scratch/core-dictionaries/handoff.md:21,23,34,47,61,75` — «локальный commit» кода в основном каталоге на `dev` — с `GIT_FLOW.md:11` (код — через ветку от `dev`) — пометить handoff историческим или удалить.
- A8. `docs/temp_plan.md:117-137,254-270` — правила: worktree на ticket, коммиты только локально, push и коммит планирования только по команде; `docs/temp_plan.md:198-216` — свой pipeline (grill-with-docs → to-spec → to-tickets → implement) — с `GIT_FLOW.md:11-14` и `~/ETEEET/AGENTS.md` «Скилы Pocock» (выбор через `ask-matt`) — удалить `temp_plan.md` (см. D1).

### B. Противоречия и устаревание ADR

- B1. `docs/adr/0001-request-options.md:24-25` — список заголовков transport без `Idempotency-Key`/`X-Idempotency-Key`; код их блокирует (`transport/options.go:110-111`); `docs/transport.md:27` тоже без `X-Idempotency-Key` — дописать в ADR и docs.
- B2. `docs/adr/0001-request-options.md:41-42` — «T09 добавит число попыток» — `ResponseInfo.Attempts` уже есть (`transport/options.go:35`) — настоящее время.
- B3. `docs/adr/0001-request-options.md:46-47` — «компилируются сейчас только методы UOM, остальных ещё нет» — все прототипы реализованы (`entity/customerorder/service.go:76,88`, `report/turnover/service.go:100`) — убрать оговорку.
- B4. `docs/notification.md:16` и `notification/settings.go:40` — `Settings(ctx)` возвращает `(notification.Settings, error)` без указателя — с ADR 0001:12-13 (`(*T, error)`) — записать исключение для map-типа в ADR или вернуть указатель.
- B5. `.scratch/document-positions/spec.md:74,78-79` — `GetPage(ctx, documentID, page)` без `opts`, query только limit/offset — с ADR 0001:12 и кодом `positions/service.go:50`, `positions/pagination.go:7` — пометить таблицу историей.
- B6. `.scratch/sdk-v1-foundation/spec.md:81` — «альтернативу `Result[T]` сравнить» — ADR 0001:69-73 её отверг — пометить историей.
- B7. решения ADR-уровня без ADR (кандидаты; оформить через `domain-modeling`): фасад группами `c.Entity()/Report()/…` и `Transport()` как мост (`.scratch/sdk-v1/map.md:156`); одна `Position/PositionRequest` для товарных документов (`.scratch/universal-positions/spec.md:11-13`); деньги — `float64` в копейках (`.scratch/sdk-v1/spec.md:92`); без скрытых запросов, автоповтор только GET/HEAD (`.scratch/sdk-v1/spec.md:96-97,122-126`); limiter в пределах процесса (`.scratch/sdk-v1/spec.md:113-118`); batch без autosplit (`.scratch/sdk-v1-optional/spec.md:95,97`); syncId — server write-once, не Idempotency-Key (`.scratch/sdk-v1/spec.md:94-95`); vendorapi без хранения состояния (`.scratch/vendor-api/spec.md:30-31`); expand — последним (`.scratch/sdk-v1/spec.md:83,156`).

### C. Расхождения docs и спек с кодом

- C1. `README.md:76-83` — таблица справочников «Get, Create, Update» и «счета… не входят» — у counterparty/organization есть List, Delete, batch, accounts, contactpersons, notes, settings (`entity/counterparty/accounts.go`, `entity/organization/accounts.go`) — обновить, сослаться на `docs/dictionaries.md`.
- C2. `README.md:107-117` — документы «Get/Create/Update» — у всех есть List, New, Trash, Delete, batch и операции строк (`entity/paymentin`, `entity/supply` и др.) — обновить колонку.
- C3. `README.md:164` — Product «Get, Create, Update» — есть List, Delete, batch, storebalances (`entity/product`) — обновить.
- C4. `docs/status-backlog.md:36` — `[x]` «fsrarId, utmUrl» — не добавлены: `.scratch/sdk-full-dictionaries/issues/16-organization-accounts.md:55`, `docs/organization.md:106` — снять отметку.
- C5. `docs/status-entities.md:39` — `[ ] gtd` — покрыт `positions.GTD` (`docs/status-backlog.md:12`) — отметить.
- C6. `docs/coverage.md:435` — `transport/query.go`, `filter.go`, `order.go` — это пакет `query/` — исправить.
- C7. `docs/access.md:42` — «создание и изменение описаний доп. полей не входит» — `metadata/attributes.go:42,56` — сослаться на `dictionaries.md`.
- C8. `docs/access.md:45` — «массовое создание статусов не входит» — `metadata/states.go:104` — сослаться на `dictionaries.md`.
- C9. `docs/catalog-extra.md:80` — G12/G13 не выполнены — `entity/emissionorder`, `entity/retireorder` есть — убрать из «Не входит».
- C10. `docs/eventfeed.md:5` — handle «предлагаемый» — `client.go:213` — «`c.Entity().EventFeed()`».
- C11. `docs/loss.md:50` — «`PUT /new` и wizard не входят» — `entity/loss/template.go:14`, `wizard/wizard.go:110`, `docs/documents-extra.md:33,49` — удалить строку.
- C12. `docs/internalorder.md:89` — «wizard не входит» — `docs/documents-extra.md:49`, `wizard/wizard.go:110` — убрать wizard.
- C13. `docs/documents-extra.md:54-60` — Trash без emissionorder/retireorder — у них `Trash` есть — дописать.
- C14. `docs/documents-extra.md:9-13` — операции строк без retireorder (полный набор) и emissionorder (без удаления) — `entity/retireorder/positions.go`, `entity/emissionorder/positions.go` — дописать с оговоркой.
- C15. `docs/documents.md:16-18` — список Create с проверкой Name без ProcessingOrder — `entity/processingorder/service.go:42` — дописать (или снять проверку по A6).
- C16. `docs/documents.md:56` — своя форма строки без EmissionOrder/RetireOrder — `entity/emissionorder/model.go:40`, `entity/retireorder/model.go:55` — дописать.
- C17. `docs/files-reports.md:37-38` — «изображения товаров не входят» — пакет `images` (`images/images.go:84-174`) — сослаться на `catalog-extra.md`.
- C18. `docs/client.md:51-62` — таблица примеров на фасаде без `examples/product_{get,create,update}` (они на `product.NewService`, у product_get нет `main_test.go`) — перевести на фасад и внести в таблицу или удалить.
- C19. `docs/counterparty.md:13` — «все методы возвращают `(*Counterparty, error)`» — List, Delete, ListAccounts и др. (`entity/counterparty/service.go:88-145`) — «Get/Create/Update возвращают…».
- C20. `docs/counterparty.md:116` — раздел «TODO» перечисляет реализованное (accounts, contactpersons, notes, settings) — переименовать в «См. также».
- C21. `docs/transport.md:65` — «List у 6 справочников и 11 документов» — List есть примерно у 80 сервисов — «у всех списочных сервисов».
- C22. `docs/transport.md:67` — «`fields` в List не принимается, minimumStock отложен» — `query/options.go:50`, `entity/product/service.go:90`, `internal/list/page.go:44` — переписать: allowlist по endpoint.
- C23. `docs/report-byoperations.md:3`, `docs/report-counterparty.md:3`, `docs/report-dashboard.md:3`, `docs/report-money.md:3` — «handle фасада планируется» — `client.go:197,183,194,180` — «`c.Report().X()`».
- C24. `docs/report-profit.md:3`, `docs/report-sales.md:3` — «после подключения в `client.go`» — `client.go:191,188` — убрать оговорку.
- C25. `docs/task.md:6` — «предлагаемый handle» — `client.go:361` — убрать «предлагаемый».
- C26. `docs/positions.md:82` — GetPositions «для одиннадцати документов» — их 16 (`type Position = positions.Position` в `entity/*/model.go`) — исправить.
- C27. `docs/positions.md:104` — «Create/Update всех семи товарных документов» — `PositionRequest` у 14 — исправить.
- C28. `docs/purchaseorder.md:73` — «Pack остаётся TODO» — `positions/model.go:40`, `positions/request.go:37` — удалить.
- C29. `docs/salesreturn.md:64` — «шаблоны /new не входят» — `entity/salesreturn/template.go:14` — удалить.
- C30. `docs/retail.md:53-54` — «не входят» /new, операции строк, внесения/выплаты — есть `New` и операции строк у retaildemand/retailsalesreturn, пакеты `retaildrawercashin/out` — переписать.
- C31. `docs/retail.md:21` — «ресурс кассиров не входит» — `entity/retailstore/cashiers.go:28,33` — удалить.
- C32. `docs/write-async.md:15-16` — список без Delete не содержит discount (`entity/discount/discount.go:108-141`); `entity/delete_test.go:366` не проверяет discount и application — дописать в doc и тест.
- C33. `docs/write-async.md:73` — «остальные 60 сервисов» — полный набор у 61 — исправить.
- C34. `.scratch/invoice-documents/spec.md:68-70,88`, `issues/01-invoicein.md:38`, `.scratch/document-workstream/map.md:93-94`, `approval.md:60` — Name обязателен для invoicein — `entity/invoicein/service.go:43` (v0.2.2) — пометить отклонение.
- C35. `.scratch/invoice-documents/spec.md:185`, `.scratch/stock-documents/spec.md:61`, `.scratch/document-workstream/map.md:99` — files в Create вне scope — `entity/invoicein/request.go:39`, `entity/supply/request.go:43` (v0.2.1) — дописать.
- C36. `.scratch/document-positions/spec.md:129-131` — Cost/Shipped/InTransit, Country, GTD, Pack и др. «нет в коде» — `positions/model.go:29-47` — пометить историей.
- C37. `.scratch/sdk-v1-release/spec.md:40-43` — `client.Product`, `client.States`, `client.Reports.Turnover` — реализовано `c.Entity().Product()`, `c.Report().Turnover()`, без States (`client.go:135-138`) — обновить.
- C38. `.scratch/sdk-v1-foundation/spec.md:121` — `CreateToken(ctx, options)` — `security/service.go:38` `(ctx, username, password, opts...)`; `:238` SyncID в Product Update — нет (`entity/product/request.go:63`) — дописать «отклонение, см. Answer T04/T13».
- C39. `.scratch/core-dictionaries/spec.md:152,361,383` — syncId «только Create» — Update Counterparty/Organization получили SyncID (`entity/counterparty/request.go:86`, `entity/organization/request.go:89`, решение `.scratch/review-followups/spec.md:17-21`) — пометить заменённым.
- C40. `.scratch/review-followups/worktree-audit.md:9,19-20,46`, `issues/01-followups.md:38` — сохранены ветки `feat/core-*`, `feat/document-model-foundation-20261101` и `stash@{0..4}` — фактически теги `archive/feat/core-*`, `archive/stash-0..4`; stash пуст, `feat/document-model-foundation-*` нет — обновить имена.
- C41. `.scratch/sdk-v1/coverage.md:67,70,76` — report/stock «заглушка», изображений нет, Vendor API вне — `report/stock/service.go`, пакеты `images`, `vendorapi` — пометить снимком 2026-11-02.

### D. Устаревшее

- D1. `docs/temp_plan.md` целиком — снимок 2026-11-01 (core-dictionaries, первая волна); `:7` «подключения из AGENTS.md нет», `:194` taxrate «выяснить» (реализован); `:276` сам требует удалить файл после работ — удалить, живые правила перенести (см. A8).
- D2. `docs/coverage.md:404-484` — «Сессия 2 и далее»: неотмеченные этапы выполнены (List, заголовки, retry, limiter, metadata, вложенные, фасад); `README.md:4` называет coverage «план работ» — отметить или свернуть.
- D3. `docs/status-backlog.md:10` — «Spec и tickets (DRAFT, needs-triage)» — все захватки реализованы (`.scratch/sdk-full/map.md:39-43`) — исправить; заголовки захваток 2–5 без «сделано».
- D4. `docs/status-backlog.md:163-164`, `docs/status-technical.md:106` — «тег и релиз — отдельной командой», «Vendor API — отдельный контур» — есть теги `v0.1.0`…`v0.2.2` и пакет `vendorapi` — обновить.
- D5. `docs/status-technical.md:4`, `docs/status-entities.md:3` — revision `9a811cef` — `docs/status-backlog.md:3` `0e1b3d61`, `:165` сам требует обновить — обновить.
- D6. `CHANGELOG.md:88-132` — «v1.0.0 — предложение», «не входит в v1: корзина, операции позиций, async-отчёты, Vendor API» — всё реализовано и выпущено в 0.x; `.scratch/sdk-v1-release/spec.md:112`, `issues/03-ci-release.md:44` «до v1.0 tag» — снять или переписать.
- D7. `.scratch/sdk-v1/spec.md:19,21,47,70-73,162-169` — «DRAFT/needs-triage», «предположение», Open Decisions — решены (ADR 0001, Answer тикетов) — свернуть в «Принятые решения».
- D8. `.scratch/sdk-v1/map.md:140` — «X01–X03 не входят в v1» — с `:164` (включены 2026-11-02); `:183-187` «Not yet specified» — пометить отменённым.
- D9. `sdk-v1-*/spec.md:16` «текущая работа — планирование», «Общий контракт реализации» и пустые Open Decisions (foundation 249-266, access 151-168, catalog 152-169, documents 105-122, files-reports 108-125, optional 73-90, retail 125-142, release 91-108) — пометить историей.
- D10. `.scratch/document-positions/spec.md:10`, `.scratch/document-workstream/approval.md:96` — entity tickets needs-triage/не запускались — все resolved (`.scratch/document-workstream/map.md:6-7`) — исправить.
- D11. `.scratch/document-workstream/map.md:187-188` — cleanup worktree не выполнен — выполнен (`.scratch/review-followups/issues/01-followups.md`) — исправить.
- D12. `.scratch/review-followups/issues/01-followups.md:27,31-34` — «needs-info, resolved не установлен», «5 dirty ожидают» — при `Status: resolved` и выполненном stash — убрать абзац.
- D13. `.scratch/core-dictionaries/handoff.md:11` — «пока локальные, без commit» — все 5 тикетов resolved — пометить историческим.
- D14. `docs/cashin.md:49`, `cashout.md:52`, `demand.md:78`, `invoicein.md:61`, `invoiceout.md:59`, `move.md:54`, `paymentin.md:53`, `paymentout.md:54`, `purchasereturn.md:58`, `supply.md:64` — заголовок «Ограничения и TODO» без TODO — «Ограничения».
- D15. `docs/transport.md` — вес запроса Employee «сейчас 3» — устареет 2026-12-01 — указать дату или убрать «сейчас».

### E. Формат трекера (`docs/agents/issue-tracker.md`)

- E1. ~200 файлов `.scratch/` — строка `Stage:` (DRAFT/APPROVED/SUPERSEDED) — в `docs/agents/issue-tracker.md:11-12` только `Status:` — описать `Stage:` в issue-tracker.md или перестать писать.
- E2. 52 тикета `sdk-v1-*/issues/*.md:25-40` — 7 неотмеченных `- [ ]` при `Status: resolved` и `## Answer` — отметить или пометить шаблоном.
- E3. 142 resolved-тикета (напр. `.scratch/sdk-v1-release/issues/03-ci-release.md:82`) — «Этот DRAFT не является командой реализации» — удалить.
- E4. `.scratch/review-followups/issues/01-followups.md:20`, `.scratch/document-positions/issues/01-json-contracts.md:83`, `02-get-page.md:82` — `## Comments` не в конце — с `docs/agents/issue-tracker.md:13` — перенести.
- E5. открытые спеки и research при закрытых тикетах — см. «Открытые тикеты».

## Открытые тикеты

- `.scratch/core-dictionaries/spec.md` (ready-for-agent) — закрыть: 5/5 тикетов resolved.
- `.scratch/uom/spec.md` (ready-for-agent) — закрыть: 3/3 тикета resolved.
- `.scratch/review-followups/spec.md` (needs-info) — закрыть: тикет 01 resolved, cleanup выполнен.
- `.scratch/sdk-full-dictionaries/spec.md` (needs-triage) — закрыть: 27/27 resolved.
- `.scratch/sdk-full-catalog/spec.md` (needs-triage) — закрыть: 13/13 resolved.
- `.scratch/sdk-full-production/spec.md` (needs-triage) — закрыть: 8/8 resolved.
- `.scratch/sdk-full-reports/spec.md` (needs-triage) — закрыть: 23/23 resolved.
- `.scratch/sdk-full-documents/spec.md` (needs-triage) — оставить до O07 (19/20 resolved); статус needs-triage неверен.
- `.scratch/sdk-full/map.md` (needs-triage) — оставить до O07; все захватки реализованы, статус needs-triage неверен.
- `.scratch/sdk-full-documents/issues/07-cash-accounts.md` (needs-info) — оставить: блокер «live-проверка с согласия» снят `~/ETEEET/moysklad/AGENTS.md:16`, перевести в ready-for-agent и проверить счета cashin/cashout на `admin@documan`.
- `.scratch/document-positions/research.md` (needs-triage) — закрыть: спека resolved.
- `.scratch/document-workstream/triage-research.md` (needs-triage) — закрыть: карта resolved.
- `.scratch/financial-documents/research.md` (needs-triage) — закрыть: спека resolved.
- `.scratch/stock-documents/research.md` (needs-triage) — закрыть: спека resolved.
- `.scratch/sdk-v1/contracts.md` (needs-triage) — закрыть как справочный каталог (resolved): sdk-v1 resolved, на него ссылается `.scratch/sdk-v1/spec.md:89`.

## Черновик baseline

Действующие решения ADR. ADR в `docs/adr/` один; общих ADR в `~/ETEEET/moysklad/meta/docs/adr/` нет.

- ADR 0001: последний параметр каждого публичного метода — `opts ...transport.RequestOption`; результаты существующих методов `(*T, error)` не меняются.
- ADR 0001: `transport.WithHeader` добавляет значение; повтор ключа даёт несколько значений, не CSV.
- ADR 0001: `transport.WithResponseInfo` заполняет StatusCode и копию заголовков последнего ответа, в т.ч. при `APIError` и ошибке decode; без ответа — обнуляется.
- ADR 0001: публичная операция обнуляет ResponseInfo на входе (`ResetResponseInfo`) до локальной валидации; общего `LastResponse` нет.
- ADR 0001: приоритет заголовков — transport < `Config.Header` < опции вызова; ключ из опций заменяет ключ клиента целиком.
- ADR 0001: служебные заголовки (`Authorization`, `Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`, `Host`) задаёт только transport; попытка задать — ошибка до отправки.
- ADR 0001: `transport.WithAuth` заменяет `Config.Auth` на один вызов; `Authorization` по-прежнему ставит transport.
- ADR 0001: `transport.Client.Send` (и `SendNoRedirect`) — единственный seam JSON-вызова; limiter и retry живут в transport, не в services.
- ADR 0001: `Download` пишет тело в `io.Writer`, закрывает его сам; слоты limiter держатся до конца копирования; retry — только пока в writer ничего не записано.
- ADR 0001: policy (настройки, limiter, retry) принадлежит `*transport.Client`; services и фасад на одном Client делят её.
- ADR 0001: `NewRequest`/`Do` — низкоуровневый API без опций и без ResponseInfo.
- ADR 0001: ResponseInfo содержит число попыток; StatusCode и Header — последней попытки.
- ADR 0001: query-параметры — отдельный аргумент конкретного метода, не RequestOption.
- ADR 0001: `Result[T]{Value, Info}` отвергнут.

## Покрытие

- ADR: `docs/adr/0001-request-options.md` — B1, B2, B3.
- Спеки:
  - `core-dictionaries/spec.md` — C39, открытые; `handoff.md` — A7, D13.
  - `document-common/spec.md` — согласован.
  - `document-positions/spec.md` — B5, C36, D10, E4.
  - `document-workstream/map.md` — C34, C35, D11; `approval.md` — C34, D10; `wave-1.md` — согласован.
  - `financial-documents/spec.md`, `contract.md` — согласованы.
  - `invoice-documents/spec.md` — A6, C34, C35.
  - `purchase-return/spec.md` — согласован.
  - `review-followups/spec.md` — A4, C40, D12, E4, открытые.
  - `sdk-full/map.md` — открытые.
  - `sdk-full-dictionaries`, `sdk-full-catalog`, `sdk-full-production`, `sdk-full-reports` `/spec.md` — A4, A5, открытые.
  - `sdk-full-documents/spec.md` — A4, A5, открытые.
  - `sdk-v1/spec.md` — D7, B7; `map.md` — A1, A4, D8; `coverage.md` — C41.
  - `sdk-v1-foundation/spec.md` — B6, C38, D9, E2.
  - `sdk-v1-access`, `sdk-v1-catalog`, `sdk-v1-documents`, `sdk-v1-files-reports`, `sdk-v1-optional` `/spec.md` — D9, E2, E3.
  - `sdk-v1-retail/spec.md` — A6, D9, E2.
  - `sdk-v1-release/spec.md` — A1, C37, D6, D9, E2, E3.
  - `stock-documents/spec.md` — C35.
  - `stock-move/spec.md`, `universal-positions/spec.md`, `vendor-api/spec.md` — согласованы.
  - `uom/spec.md` — согласован, открытые.
- Открытый тикет: `sdk-full-documents/issues/07-cash-accounts.md` — A4, A5, открытые.
- Docs с находками: `README.md` (C1–C3), `CHANGELOG.md` (D6), `docs/coverage.md` (A1, A5, C6, D2), `docs/status-backlog.md` (A1, A5, C4, D3–D5), `docs/status-technical.md` (A1, A5, D4, D5), `docs/status-entities.md` (C5, D5), `docs/temp_plan.md` (A4, A5, A8, D1), `docs/agents/issue-tracker.md` (E1), `access`, `catalog-extra`, `client`, `counterparty`, `documents`, `documents-extra`, `eventfeed`, `files-reports`, `internalorder`, `loss`, `notification`, `positions`, `purchaseorder`, `report-byoperations`, `report-counterparty`, `report-dashboard`, `report-money`, `report-profit`, `report-sales`, `retail`, `salesreturn`, `task`, `transport`, `write-async` (C7–C33, B1, B4, D15), а также заголовки D14.
- Docs согласованы: `docs/agents/domain.md`, `docs/agents/triage-labels.md`, `docs/upstream.md` (проверена шапка снимка), `application`, `catalog`, `commissionreportin`, `commissionreportout`, `counterpartyadjustment`, `country`, `customerorder`, `dictionaries`, `document-common`, `emissionorder`, `enter`, `expenseitem`, `facturein`, `factureout`, `inventory`, `namedfilter`, `organization`, `payroll`, `prepayment`, `prepaymentreturn`, `pricelist`, `print`, `processing`, `processingorder`, `processingplan`, `processingplanfolder`, `processingprocess`, `processingstage`, `productionstagecompletion`, `productiontask`, `publication`, `report-stock`, `requisites`, `retaildrawercashin`, `retaildrawercashout`, `retireorder`, `vendorapi`, `webhook`, `webhookstock`, `testdata/readme.md`.

Не проверено: `.scratch/sdk-v1/contracts.md` построчно (2151 строка), тела review-файлов, поведение живого API (A6, O07), внешние ссылки на документацию МойСклад.
