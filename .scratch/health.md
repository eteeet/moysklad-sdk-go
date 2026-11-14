# Здоровье документов

Пересборка 2026-11-11, `dev` @ `3fc370b`; сверка 2026-11-14, `dev` @ `159992b`: изменились `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md` и добавлен тикет `trash-content-type/01`, код не менялся. Только действующая цепочка; история (закрытые спеки и тикеты, research, ревью, handoff, снимки) не сверяется.
Цепочка: 1 md → 2 ADR и glossary → 3 открытые спеки → 4 открытые тикеты и бэклог → 5 docs и README → 6 код.
Строка: `файл:строка — с каким звеном расходится — предложение — кто решает`. «Механика» — однозначная правка вниз; «владелец» — выбор или правка верхнего звена.

## 1. md

- M2. `.scratch/sdk-full-documents/spec.md:2`, `.scratch/sdk-full-documents/issues/07-cash-accounts.md:2`, `.scratch/sdk-full/map.md:2` — строка `Stage:` не описана в `docs/agents/issue-tracker.md:11-12` (только `Status:`) — убрать `Stage:` из открытых файлов — механика.
- M3. `.scratch/trash-content-type/issues/01-trash-content-type.md:2` — `Type: bug` нет в `docs/agents/issue-tracker.md:29` (`research`, `prototype`, `grilling`, `task`) — `Type: task` — механика.

## 2. ADR и glossary

Все находки — ADR 0001 устарел относительно кода, который его исполняет; правка верхнего звена.

- A1. `docs/adr/0001-request-options.md:24-25` — список заголовков transport без `Idempotency-Key`, `X-Idempotency-Key`; код их блокирует (`transport/options.go:110-111`) — дописать в ADR — владелец.
- A2. `docs/adr/0001-request-options.md:41-42` — «T09 добавит число попыток»; `ResponseInfo.Attempts` есть (`transport/options.go:35`) — настоящее время — владелец.
- A3. `docs/adr/0001-request-options.md:46-47` — «компилируются только методы UOM, остальных ещё нет»; все прототипы реализованы (`entity/customerorder/service.go:76,88`, `report/turnover/service.go:100`) — убрать оговорку — владелец.
- A4. `~/ETEEET/moysklad/meta/GLOSSARY.md:131,101,141` (Position: _Avoid_ «строка»; SKU: _Avoid_ «ассортимент»; Batch — файлы одного Ingest) — docs SDK называют позиции «строками» (`README.md:63`, `docs/documents-extra.md`, `docs/positions.md`), пишут «ассортимент» (`README.md:27`), пакет массовых операций назван `batch` — решить: glossary действует на SDK, называющий ресурсы API, или SDK — исключение; при «действует» — правка docs (механика) и вопрос о пакете `batch` (ломающее переименование) — владелец.

## 3. Открытые спеки

- S1. `.scratch/sdk-full-documents/spec.md:138` — «needs-info: ответ МойСклад или live-проверка с согласия» — с `~/ETEEET/moysklad/AGENTS.md:16` (тестовый аккаунт — без запроса) — «live-проверка на `admin@documan`» — механика.
- S2. `.scratch/sdk-full-documents/spec.md:270-271` — «Planning не коммитить без отдельной команды», «Live API… не входят» — с `GIT_FLOW.md:12` и `~/ETEEET/moysklad/AGENTS.md:16` — убрать обе оговорки — механика.

## 4. Открытые тикеты и бэклог

- T1. `.scratch/sdk-full-documents/issues/07-cash-accounts.md:18,39` — «live-проверка с согласия», «Planning без команды не коммитится» — с `~/ETEEET/moysklad/AGENTS.md:16`, `GIT_FLOW.md:12` — переписать, Status → `ready-for-agent` — механика.
- T2. `docs/status-backlog.md:160-161` — live-проверки «только с согласия» — с `~/ETEEET/moysklad/AGENTS.md:16` — «на `admin@documan` без запроса» — механика.
- T3. `docs/status-backlog.md:162`, `docs/status-technical.md:104`, `docs/coverage.md:484` и сам репозиторий: нет `.github/workflows/`, записано «CI для v1 не делается, только `make check`» — с `INFRA.md:17-24` (CI на GitHub Actions у каждого репозитория с кодом: проверки на push ветки и PR, на push в `dev`; тег у библиотеки ничего не собирает), `GIT_FLOW.md:13` (ветка сливается в `dev` при зелёном CI ветки) и `GIT_FLOW.md:19` (тег — на коммит с зелёным CI; `v0.1.0`…`v0.2.2` без CI; выпуск исправления `trash-content-type/01` упирается в это же) — workflow по INFRA отдельным тикетом либо исключение в ADR; способ CI меняется только после «да» владельца (`INFRA.md:31`) — владелец.
- T4. `docs/status-backlog.md:163`, `docs/status-technical.md:106` — «тег и релиз — отдельной командой» — с `GIT_FLOW.md:16-24` (процедура выпуска; теги уже есть) — ссылка на GIT_FLOW — механика.
- T5. `docs/status-backlog.md:164` — «Vendor API — отдельный контур» — с кодом (`vendorapi/`) — отметить сделанным — механика.
- T6. `docs/status-backlog.md:36` — `[x]` «fsrarId, utmUrl» — с кодом (`entity/organization` их не содержит) и `docs/organization.md:106` — отметку оставить счетам, fsrarId/utmUrl — «не добавлены: нет в таблице атрибутов» — механика.
- T7. `docs/status-backlog.md:10` и заголовки захваток 2–5 — «Spec и tickets (DRAFT, needs-triage)», нет «сделано» — с `.scratch/sdk-full/map.md:39-43` (реализованы) — исправить — механика.
- T8. `docs/status-backlog.md:165` — пункт «обновить revision в status-*» не выполнен: `docs/status-technical.md:4`, `docs/status-entities.md:3` — `9a811cef` против `docs/status-backlog.md:3` `0e1b3d61` — обновить — механика.
- T9. `.scratch/trash-content-type/issues/01-trash-content-type.md:17` — в критерии «ms_documan убирает обход в `.scratch/replica/livecheck/main.go`, когда поднимет версию» — работа другого репозитория в тикете SDK, с `~/ETEEET/moysklad/AGENTS.md:26`; код продукта ms_documan `Trash` не вызывает (нет вызовов в `internal/`, `cmd/`), обход — только в одноразовом скрипте на SDK v0.1.1, тикет в ms_documan не нужен — убрать эту часть, «Выпуск SDK с исправлением» оставить — механика.

## 5. docs и README

- D1. `docs/temp_plan.md` — правила (`:117-137`, `:254-270`: worktree на ticket, коммиты только локально, push и коммит planning по команде; `:198-216`: свой pipeline grill-with-docs → to-spec → to-tickets → implement) — с `GIT_FLOW.md:11-14` и `~/ETEEET/AGENTS.md` «Скилы Pocock»; `:276` сам требует удалить файл — решить, переносить ли правила тестов `:237-252` (upstream-фикстуры без изменений, независимые golden, expand последним) в `AGENTS.md` или ADR, затем удалить файл — владелец.
- D2. `CHANGELOG.md:88-132` — «v1.0.0 — предложение»: «не входит в v1: корзина, операции позиций, async-отчёты, Vendor API» — с кодом (всё есть, выпущено в 0.x) и `GIT_FLOW.md:18` (версию предлагает агент при выпуске) — снять раздел или переписать предложение — владелец.
- D3. `docs/coverage.md:9-10,367,401` — реальная запись «только после отдельного согласия», «коммитить после ревью пользователя» — с `~/ETEEET/moysklad/AGENTS.md:16`, `GIT_FLOW.md:12-13` — переписать — механика.
- D4. `docs/coverage.md:404-484` — этапы «Сессия 2 и далее» не отмечены, хотя сделаны (List, заголовки, retry, limiter, metadata, вложенные ресурсы, фасад); `:435` — пути `transport/query.go`, `filter.go`, `order.go` вместо пакета `query/`; `README.md:4` — coverage как «план работ» — привести к коду — механика.
- D5. `docs/status-entities.md:39` — `[ ] gtd` — с `docs/status-backlog.md:12` и кодом `positions.GTD` — отметить — механика.
- D6. `README.md:76-83`, `:107-117`, `:164` — таблицы справочников, документов и Product перечисляют только Get/Create/Update, «счета… не входят» — с кодом (List, Delete, batch, New, Trash, операции строк, accounts, contactpersons, notes, settings, storebalances) — обновить таблицы, сослаться на `docs/dictionaries.md`, `docs/documents-extra.md`, `docs/write-async.md` — механика.
- D7. handle фасада «планируется»/«предлагаемый»: `docs/eventfeed.md:5`, `docs/report-byoperations.md:3`, `docs/report-counterparty.md:3`, `docs/report-dashboard.md:3`, `docs/report-money.md:3`, `docs/report-profit.md:3`, `docs/report-sales.md:3`, `docs/task.md:6` — с кодом (`client.go:180,183,188,191,194,197,213,361`) — «`c.Report().X()`» / «`c.Entity().X()`» — механика.
- D8. «не входит»/TODO про реализованное: `docs/access.md:42` (`metadata/attributes.go:42,56`), `docs/access.md:45` (`metadata/states.go:104`), `docs/catalog-extra.md:80` (`entity/emissionorder`, `entity/retireorder`), `docs/loss.md:50` (`entity/loss/template.go:14`, `wizard/wizard.go:110`), `docs/internalorder.md:89` (`wizard/wizard.go:110`), `docs/files-reports.md:37-38` (пакет `images`), `docs/purchaseorder.md:73` (`positions/model.go:40`), `docs/salesreturn.md:64` (`entity/salesreturn/template.go:14`), `docs/retail.md:21` (`entity/retailstore/cashiers.go:28,33`), `docs/retail.md:53-54` (`New` и операции строк retaildemand/retailsalesreturn, пакеты `retaildrawercashin/out`), `docs/counterparty.md:116` (accounts, contactpersons, notes, settings) — убрать или заменить ссылкой — механика.
- D9. неполные перечни: `docs/documents-extra.md:9-13` (нет retireorder; emissionorder — без удаления строк), `docs/documents-extra.md:54-60` (Trash без emissionorder, retireorder), `docs/documents.md:16-18` (проверка Name без ProcessingOrder, `entity/processingorder/service.go:42`), `docs/documents.md:56` (своя форма строки без EmissionOrder, RetireOrder), `docs/write-async.md:15-16` (без Delete нет discount, `entity/discount/discount.go:108-141`), `docs/transport.md:27` (нет `X-Idempotency-Key`, `transport/options.go:111`) — дописать — механика.
- D10. неверные числа: `docs/transport.md:65` (List «у 6 справочников и 11 документов» — около 80 сервисов), `docs/positions.md:82` (11 → 16), `docs/positions.md:104` (7 → 14), `docs/write-async.md:73` (60 → 61) — исправить — механика.
- D11. неверные утверждения: `docs/transport.md:67` («`fields` в List не принимается» — `query/options.go:50`, `entity/product/service.go:90`, allowlist `internal/list/page.go:44`), `docs/counterparty.md:13` («все методы возвращают `(*Counterparty, error)`» — List, Delete, ListAccounts и др.), `docs/client.md:51-62` (таблица примеров на фасаде без `examples/product_{get,create,update}` на `product.NewService`) — привести к коду — механика.

## 6. Код

- K1. Create отказывает при пустом `Name`: `entity/invoiceout/service.go:43`, `entity/retailsalesreturn/service.go:45`, `entity/processingorder/service.go:42` — с `~/ETEEET/moysklad/AGENTS.md:18` (документам не задавать `name`, МойСклад нумерует сам) — подтвердить, что правило распространяется на валидацию SDK; тогда тикет: проверить автонумерацию на `admin@documan` и снять проверку, как для invoicein в v0.2.2 — владелец.
- K2. `notification/settings.go:40` (описан в `docs/notification.md:16`) — `Settings(ctx)` возвращает `(notification.Settings, error)` без указателя — с ADR 0001:12-13 (`(*T, error)`) — исключение для map-типа в ADR или тикет на указатель (ломающее изменение) — владелец.
- K3. `origin/main` = `5f3fb93 Initial commit`, последний выпуск `v0.2.2` — с `GIT_FLOW.md:20` (владелец перематывает `main` на тег) — `git push origin 'v0.2.2^{}:main'` — владелец.

## Открытые тикеты и спеки

- `.scratch/trash-content-type/issues/01-trash-content-type.md` (ready-for-agent) — оставить: работа не сделана; M3, T9; выпуск исправления ждёт T3.
- `.scratch/sdk-full-documents/issues/07-cash-accounts.md` (needs-info) — оставить: работа не сделана; блокер «live-проверка с согласия» снят `~/ETEEET/moysklad/AGENTS.md:16`, после T1 — `ready-for-agent`.
- `.scratch/sdk-full-documents/spec.md` (needs-triage) — оставить: открыт O07.
- `.scratch/sdk-full/map.md` (needs-triage) — оставить: в карте открыт O07.
- `Status:` открыт, но все тикеты resolved — спеки уже закрыты, исправить строку на `resolved` (механика): `.scratch/core-dictionaries/spec.md` (5/5), `.scratch/uom/spec.md` (3/3), `.scratch/review-followups/spec.md` (1/1), `.scratch/sdk-full-dictionaries/spec.md` (27/27), `.scratch/sdk-full-catalog/spec.md` (13/13), `.scratch/sdk-full-production/spec.md` (8/8), `.scratch/sdk-full-reports/spec.md` (23/23).

## Черновик baseline

Действующие решения ADR; других ADR нет ни в `docs/adr/`, ни в `~/ETEEET/moysklad/meta/docs/adr/`. Сверены с md: противоречий нет. Устаревшие формулировки самого ADR — A1–A3.

- ADR 0001: последний параметр каждого публичного метода — `opts ...transport.RequestOption`; результаты существующих методов `(*T, error)` не меняются.
- ADR 0001: `transport.WithHeader` добавляет значение; повтор ключа даёт несколько значений, не CSV.
- ADR 0001: `transport.WithResponseInfo` заполняет StatusCode и копию заголовков последнего ответа, в т.ч. при `APIError` и ошибке decode; без ответа — обнуляется.
- ADR 0001: публичная операция обнуляет ResponseInfo на входе (`ResetResponseInfo`) до локальной валидации; общего `LastResponse` нет.
- ADR 0001: приоритет заголовков — transport < `Config.Header` < опции вызова; ключ из опций заменяет ключ клиента целиком.
- ADR 0001: служебные заголовки (`Authorization`, `Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`, `Host`; в коде ещё `Idempotency-Key`, `X-Idempotency-Key` — A1) задаёт только transport; попытка задать — ошибка до отправки.
- ADR 0001: `transport.WithAuth` заменяет `Config.Auth` на один вызов; `Authorization` по-прежнему ставит transport.
- ADR 0001: `transport.Client.Send` (и `SendNoRedirect`) — единственный seam JSON-вызова; limiter и retry живут в transport, не в services.
- ADR 0001: `Download` пишет тело в `io.Writer` и закрывает его сам; слоты limiter держатся до конца копирования; retry — только пока в writer ничего не записано.
- ADR 0001: policy (настройки, limiter, retry) принадлежит `*transport.Client`; services и фасад на одном Client делят её.
- ADR 0001: `NewRequest`/`Do` — низкоуровневый API без опций и без ResponseInfo.
- ADR 0001: ResponseInfo содержит число попыток; StatusCode и Header — последней попытки.
- ADR 0001: query-параметры — отдельный аргумент конкретного метода, не RequestOption.
- ADR 0001: `Result[T]{Value, Info}` отвергнут.

## Покрытие

- md: `~/ETEEET/AGENTS.md`, `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md`, `AGENTS.md`, `docs/agents/domain.md`, `docs/agents/triage-labels.md` — согласованы; `docs/agents/issue-tracker.md` — M2, M3.
- ADR: `docs/adr/0001-request-options.md` — A1–A3, K2. Glossary: общий — A4; своего `GLOSSARY.md` нет.
- Открытые спеки: `.scratch/sdk-full-documents/spec.md` — M2, S1, S2; `.scratch/sdk-full/map.md` — M2.
- Открытые тикеты: `.scratch/sdk-full-documents/issues/07-cash-accounts.md` — M2, T1; `.scratch/trash-content-type/issues/01-trash-content-type.md` — M3, T9. Бэклог: `docs/status-backlog.md` — T2–T8.
- docs с находками: `README.md`, `CHANGELOG.md`, `docs/temp_plan.md`, `docs/coverage.md`, `docs/status-technical.md`, `docs/status-entities.md`, `access`, `catalog-extra`, `client`, `counterparty`, `documents`, `documents-extra`, `eventfeed`, `files-reports`, `internalorder`, `loss`, `notification`, `positions`, `purchaseorder`, `report-byoperations`, `report-counterparty`, `report-dashboard`, `report-money`, `report-profit`, `report-sales`, `retail`, `salesreturn`, `task`, `transport`, `write-async` — D1–D11, K2.
- Остальные `docs/*.md` и `testdata/readme.md` — согласованы.

Не проверено: внешние ссылки на документацию МойСклад, поведение живого API (K1, O07), ссылки на API над пакетами в коде (`~/ETEEET/AGENTS.md`, «Комментарии в коде»).
