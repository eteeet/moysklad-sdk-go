# Review захватки 4 «Производство» — 2026-11-03

BASE: `c3bace0`. HEAD review: `2dbdbf3`. Исправления: `3778088`.
Scope: M01–M08, 10 коммитов, 116 файлов. Метод: независимые Standards и Spec reviewers (subagents;
Spec разделён на справочники M01–M04 и документы M05–M08), проверка root по источнику,
make check.

## Standards

Жёстких нарушений нет: ResetResponseInfo, один префикс `moysklad:`, UUID до HTTP, запись без
повторов, Field+omitzero, memnet, отсутствие `shortcut:`, success-проверки по значениям.

На суждение:
1. processingprocess/positions.go: комментарий «лимит 100 проверяет сервер» противоречил коду
   (Create/Update проверяют 1–100 до запроса). — Исправлено: комментарий описывает, что Create/Update
   проверяют лимит, CreatePositions — сервер. Лимит 100 оставлен локальным: это лимит размера тела
   из текста источника, как 1000 у документов.
2. Duplicated Code: локальные копии nested.List (processing, processingprocess, productiontask) и
   pageCases в тестах (processingplan, productionstagecompletion). — Не исправлено: общий контракт
   волны запрещает новые общие абстракции; копии повторяют принятую практику (validID/do).
3. processingplan GetProduct терял причину ошибки разбора. — Исправлено.
4. Имена: AddMaterial (M08) против CreateMaterial (M05, M07). — Исправлено: M08 CreateMaterial,
   MaterialCreateRequest. `positions.PageOptions` в processing — псевдоним `query.PageOptions`,
   тот же тип, как у документов прошлых волн; не менялось.
5. Разнобой текста операций в ошибках. — Не исправлено (каждый текст содержит тип и id; на поведение
   не влияет).
6. Дублирование проверки href этапа в processingprocess. — Не исправлено (две строки).
7. processingorder: ошибка positions.NewService проглатывается. — Не исправлено: тот же приём и
   комментарий у internalorder (ошибка возможна только для nil client и неверного кода).
8. ResetResponseInfo в ветке ошибки перед nested.*. — Не исправлено: nested сбрасывает сам,
   прецеденты bundle.AddComponent, counterparty.CreateAccount.

## Spec

M01–M04:
- (b) GetParametricMaterial: отдельного раздела нет, есть curl-пример в разделе списка
  параметрических материалов. Оставлен, записан в godoc и docs; ticket называет вложенные методы.
- (c) ID строк техкарты помечен readonly. — Не ошибка: правило проекта допускает id/meta существующего
  объекта для адресации (payroll, commissionreportin, positions.CoreRequest), общее правило позиций
  `_general.md` сопоставляет строки по id. Docs processingplan дополнены этим основанием.
- Остальное сверено с таблицами: поля моделей, readonly, деньги, пути, meta.type.

M05–M08:
- (a) Порядок коммитов M05–M08 раньше M04 и M08 раньше M07: зависимости между пакетами только через
  JSON-ссылки (model.Reference), кодовой зависимости нет; каждый коммит проходит тесты пакета.
  Историю не переписывали.
- (a) M07 — два коммита: 486d0d6 и 9fb5e31 (пустая строка в конце фикстуры из источника ломала
  diff-check). Не переписывалось (не соседние коммиты).
- (a) docs/status-backlog.md не отслеживается git по решению пользователя; пункты отмечены локально.
- (c1) M07 CreateMaterial/CreateProducts не проверяли assortment и planQuantity. — Исправлено,
  тесты отказа до HTTP добавлены (мутация роняет оба случая).
- (c2) processing: условная обязательность (техкарта или products+processingSum) — сервер. Маркера
  «Необходимо при создании» у этих полей нет, правило условное; записано в godoc и docs.
- (c3) productiontask: `updated` позиции без маркера readonly не пишется — серверная метка времени;
  записано в docs.

## Решения пользователя (не реализовано, ждёт)

- M01: поле родительской группы техкарт (в таблице нет; текст «Изменить» упоминает processingplanfolder).
- M04: `characteristicValue` параметрических строк — readonly, но примеры Create передают; без него
  параметрические материалы не сопоставить значениям характеристики.
- M06: `attributes` (и `shared`) заказа на производство — readonly, примеры Create/Update передают.
- M07: `productionRow` продукта — readonly, но раздел «Создать продукт» требует его; путь создания
  продукта `…/product` (пример) против `…/products` (SDK).
- M08: `materials`/`products` в теле Create/Update — не readonly, но формы записи и примеров нет.

## Verification

- make check (CHECK_P=2) на `2dbdbf3` — PASS; на `3778088` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., diff-check HEAD и рабочего дерева.
- Live API не вызывался.
