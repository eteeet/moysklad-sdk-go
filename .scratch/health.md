# Здоровье документов

Пересборка 2026-11-11, `dev` @ `3fc370b`; сверка 2026-11-14, `dev` @ `159992b`: изменились `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md` и добавлен тикет `trash-content-type/01`, код не менялся. Только действующая цепочка; история (закрытые спеки и тикеты, research, ревью, handoff, снимки) не сверяется.
Цепочка: 1 md → 2 ADR и glossary → 3 открытые спеки → 4 открытые тикеты и бэклог → 5 docs и README → 6 код.
Строка: `файл:строка — с каким звеном расходится — предложение — кто решает`. «Механика» — однозначная правка вниз; «владелец» — выбор или правка верхнего звена.

## 1. md

Находок нет.

## 2. ADR и glossary

Находок нет.

## 3. Открытые спеки

Находок нет.

## 4. Открытые тикеты и бэклог

Находок нет.

## 5. docs и README

- D3. `docs/coverage.md:9-10,367,401` — реальная запись «только после отдельного согласия», «коммитить после ревью пользователя» — с `~/ETEEET/moysklad/AGENTS.md:16`, `GIT_FLOW.md:12-13` — переписать — механика.
- D4. `docs/coverage.md:404-484` — этапы «Сессия 2 и далее» не отмечены, хотя сделаны (List, заголовки, retry, limiter, metadata, вложенные ресурсы, фасад); `:435` — пути `transport/query.go`, `filter.go`, `order.go` вместо пакета `query/`; `README.md:4` — coverage как «план работ» — привести к коду — механика.
- D5. `docs/status-entities.md:39` — `[ ] gtd` — с `docs/status-backlog.md:12` и кодом `positions.GTD` — отметить — механика.
- D6. `README.md:76-83`, `:107-117`, `:164` — таблицы справочников, документов и Product перечисляют только Get/Create/Update, «счета… не входят» — с кодом (List, Delete, batch, New, Trash, операции строк, accounts, contactpersons, notes, settings, storebalances) — обновить таблицы, сослаться на `docs/dictionaries.md`, `docs/documents-extra.md`, `docs/write-async.md` — механика.
- D7. handle фасада «планируется»/«предлагаемый»: `docs/eventfeed.md:5`, `docs/report-byoperations.md:3`, `docs/report-counterparty.md:3`, `docs/report-dashboard.md:3`, `docs/report-money.md:3`, `docs/report-profit.md:3`, `docs/report-sales.md:3`, `docs/task.md:6` — с кодом (`client.go:180,183,188,191,194,197,213,361`) — «`c.Report().X()`» / «`c.Entity().X()`» — механика.
- D8. «не входит»/TODO про реализованное: `docs/access.md:42` (`metadata/attributes.go:42,56`), `docs/access.md:45` (`metadata/states.go:104`), `docs/catalog-extra.md:80` (`entity/emissionorder`, `entity/retireorder`), `docs/loss.md:50` (`entity/loss/template.go:14`, `wizard/wizard.go:110`), `docs/internalorder.md:89` (`wizard/wizard.go:110`), `docs/files-reports.md:37-38` (пакет `images`), `docs/purchaseorder.md:73` (`positions/model.go:40`), `docs/salesreturn.md:64` (`entity/salesreturn/template.go:14`), `docs/retail.md:21` (`entity/retailstore/cashiers.go:28,33`), `docs/retail.md:53-54` (`New` и операции строк retaildemand/retailsalesreturn, пакеты `retaildrawercashin/out`), `docs/counterparty.md:116` (accounts, contactpersons, notes, settings) — убрать или заменить ссылкой — механика.
- D9. неполные перечни: `docs/documents-extra.md:9-13` (нет retireorder; emissionorder — без удаления строк), `docs/documents-extra.md:54-60` (Trash без emissionorder, retireorder), `docs/documents.md:56` (своя форма строки без EmissionOrder, RetireOrder), `docs/write-async.md:15-16` (без Delete нет discount, `entity/discount/discount.go:108-141`), `docs/transport.md:27` (нет `X-Idempotency-Key`, `transport/options.go:111`) — дописать — механика.
- D10. неверные числа: `docs/transport.md:65` (List «у 6 справочников и 11 документов» — около 80 сервисов), `docs/positions.md:82` (11 → 16), `docs/positions.md:104` (7 → 14), `docs/write-async.md:73` (60 → 61) — исправить — механика.
- D11. неверные утверждения: `docs/transport.md:67` («`fields` в List не принимается» — `query/options.go:50`, `entity/product/service.go:90`, allowlist `internal/list/page.go:44`), `docs/counterparty.md:13` («все методы возвращают `(*Counterparty, error)`» — List, Delete, ListAccounts и др.), `docs/client.md:51-62` (таблица примеров на фасаде без `examples/product_{get,create,update}` на `product.NewService`) — привести к коду — механика.
- D12. glossary в прозе docs (решение владельца по A4: `~/ETEEET/moysklad/meta/GLOSSARY.md` действует на прозу docs SDK; Go-идентификаторы, повторяющие API, — `assortment`, `positions`, `batch` — не меняются): «строка» в значении Position (`GLOSSARY.md:131`) — `README.md:63`, `docs/documents-extra.md`, `docs/positions.md` → «позиция»; «ассортимент» (`GLOSSARY.md:101`) — `README.md:27` и `docs/*` → `assortment` или «SKU» по смыслу; пакет `batch` в прозе (`GLOSSARY.md:141`) → «массовые операции» — механика.

## 6. Код

- K3. `origin/main` = `5f3fb93 Initial commit`, последний выпуск `v0.2.2` — с `GIT_FLOW.md:20` (владелец перематывает `main` на тег) — `git push origin 'v0.2.2^{}:main'` — владелец.

## Открытые тикеты и спеки

- `.scratch/document-name-autonumber/issues/01-create-without-name.md` (ready-for-agent) — оставить: заведён по K1, сверен с `~/ETEEET/moysklad/AGENTS.md:16,18`; включает `docs/documents.md:16-18` (бывший фрагмент D9).
- `.scratch/ci/issues/01-github-actions.md` (ready-for-agent) — оставить: заведён по T3 после «да» владельца на смену способа CI, сверен с `INFRA.md:17-24,31` и `GIT_FLOW.md:13,19`.
- `.scratch/adr-0002-code-refs/issues/01-assortment-comment.md` (ready-for-agent) — оставить: заведён по A5, сверен, согласован.
- `.scratch/trash-content-type/issues/01-trash-content-type.md` (ready-for-agent) — оставить: работа не сделана; выпуск исправления ждёт `ci/01`.
- `.scratch/sdk-full-documents/issues/07-cash-accounts.md` (ready-for-agent) — оставить: работа не сделана.
- `.scratch/sdk-full-documents/spec.md` (needs-triage) — оставить: открыт O07.
- `.scratch/sdk-full/map.md` (needs-triage) — оставить: в карте открыт O07.

## Покрытие

- md: `~/ETEEET/AGENTS.md`, `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md`, `AGENTS.md`, `docs/agents/domain.md`, `docs/agents/triage-labels.md` — согласованы; `docs/agents/issue-tracker.md` — согласован.
- ADR: `docs/adr/0002-baseline.md` — согласован; `0001` заменён им. Glossary: общий — D12; своего `GLOSSARY.md` нет.
- Открытые спеки: `.scratch/sdk-full-documents/spec.md`, `.scratch/sdk-full/map.md` — согласованы.
- Открытые тикеты: `07-cash-accounts`, `trash-content-type/01` и новые тикеты — согласованы. Бэклог: `docs/status-backlog.md` — согласован.
- docs с находками: `README.md`, `docs/coverage.md`, `docs/status-technical.md`, `docs/status-entities.md`, `access`, `catalog-extra`, `client`, `counterparty`, `documents`, `documents-extra`, `eventfeed`, `files-reports`, `internalorder`, `loss`, `notification`, `positions`, `purchaseorder`, `report-byoperations`, `report-counterparty`, `report-dashboard`, `report-money`, `report-profit`, `report-sales`, `retail`, `salesreturn`, `task`, `transport`, `write-async` — D3–D12.
- Остальные `docs/*.md` и `testdata/readme.md` — согласованы.

Не проверено: внешние ссылки на документацию МойСклад, поведение живого API (O07; автонумерация — в `document-name-autonumber/01`), ссылки на API над пакетами в коде (`~/ETEEET/AGENTS.md`, «Комментарии в коде»).
