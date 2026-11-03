# Review реализации справочников — 2026-11-03

BASE: `30e5d769fb27ec7a19ffe47394ffd796120aff0b`.
HEAD: `dd628cad53627b4390aad3cf680b234a188e1458`.
Scope: `git diff BASE...HEAD`, 115 файлов, S01–S27 реализованного среза.
S10 Consignment намеренно ждёт G01; это не пропуск реализованной волны.
Метод: независимые Standards/Spec reviewers, дополнительная проверка root,
полный make check. Предыдущие внутренние reviews 6b3e272/246383e прочитаны;
исправленные замечания не переоткрывались.

## Standards

### S1. P3 — устаревший godoc массового создания характеристик

`entity/variant/characteristic.go:85–86` утверждает, что массовое создание
характеристик не входит в SDK v1. Эта же волна добавляет публичный
CreateManyCharacteristics на строке 136.

Убрать устаревшее ограничение или сослаться на новый метод. Основание:
Definition of Done общего spec — публичная документация отражает реализованный scope.

P1/P2 по Standards не обнаружены. Проверенные публичные методы используют transport,
сохраняют причины ошибок и сбрасывают ResponseInfo до валидации. Дубли do/validID
допустимы по coverage.md. Существенных новых baseline smells не найдено.

## Spec

### P1. P2 — Slot.Zone пишет read-only поле без утверждённого исключения

`entity/store/request.go:63,71`: Zone есть в SlotCreateRequest и SlotUpdateRequest.
Первичный `_store.md:1814` явно помечает zone `+Только для чтения`.
`docs/documents.md:67–70` задаёт правило: явная пометка в таблице имеет приоритет
над противоречащими примерами запросов. Новый spec также требует разделить readonly
и исключения до реализации.

S01 Answer:55, godoc request и docs/dictionaries.md:21–22 признают расхождение,
но объясняют запись только примерами API. Отдельного утверждения пользователем
исключения Slot.Zone в доступных решениях нет. Согласованное ранее исключение
Bundle.Component.Quantity на другие поля не распространяется.

Исправление: убрать Zone из write requests, сохранив чтение, либо явно согласовать
исключение и зафиксировать его в нормативном контракте с соответствующими тестами.
Это замечание к соблюдению выбранной policy; оно не доказывает, что живой API
отклоняет Zone. Live API в review не вызывался.

Других подтверждённых Spec findings нет. Просмотренные модели и capability matrices
соответствуют таблицам и согласованным решениям. Исправления Contract name/rewardPercent,
региональных реквизитов Organization и форм UserSettings Default* присутствуют.
Root дополнительно проверил internal/nested, internal/single, расширение internal/write,
settings, attributes и customentity; новых подтверждённых дефектов не найдено.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., diff-check полного tracked HEAD и working tree.
- Источники — локальный архив api-remap-1.2-doc 9a811cef. Ранее официальным GitHub
  compare подтверждено, что переход на 0e1b3d61 меняет только code-fence разметку
  labels/print; проверенные контракты совпадают.
- HEAD оставался dd628ca. Tracked code, planning и tracker statuses reviewer не менял.
- Создан только этот отчёт. Существующие untracked planning и .DS_Store
  не изменялись; .DS_Store не читались. Live API, push и release не выполнялись.

Standards: 1 finding P3. Spec: 1 finding P2.

## Answer

Оба замечания приняты, исправлено в `c81f2ff`: zone убран из запросов ячейки (чтение сохранено),
godoc CreateCharacteristic ссылается на CreateManyCharacteristics. make check — PASS.

Позже пользователь утвердил исключение для slot.zone: запись возвращена с фиксацией противоречия и решения.
