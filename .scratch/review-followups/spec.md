Status: needs-info
Stage: APPROVED

# Правки после review P03 и первой волны

Пользователь 2026-11-01 согласовал реализацию замечаний после проверки кода.
BASE review: f26ec329f4375e64782dfec2b24bc98b1ad240d5 (текущий HEAD до правок).

## Согласованный scope

- SourceSlot/TargetSlot модели positions.Movement используют omitzero: nil не
  сериализуется, непустая ссылка сохраняется. Request не меняется.
- Публичные JSON-тесты обоими codecs проверяют обе nil-ссылки, каждую отдельно
  и обе заполненные, с независимыми ожидаемыми JSON.
- readerValidID переименован в validID без изменения поведения/публичного API;
  имя существующих тестов и внутреннее устройство helper не закрепляются.
- Counterparty.UpdateRequest и Organization.UpdateRequest получают SyncID Field[string]
  с omitzero. Unset не отправляется; значение, пустая строка и null кодируются явно.
  Это разрешение попытки записи: сервер проверяет write-once после заполнения;
  SDK не обещает очистку/перезапись и не выполняет скрытый GET или retry.
  Это заменяет прежнее ограничение Create-only из spec core-dictionaries.
- JSON-тесты SyncID обоими codecs и локальные HTTP-тесты через Service.Update:
  SyncID-only request, точный PUT/path/body, один запрос, ответ сервера и APIError.
  Существующие нулевые UpdateRequest по-прежнему отклоняются до HTTP.
- Удалить восемь вспомогательных worktree только после проверки чистоты и переноса
  изменений; не применять --force, не удалять branches. Prune только устаревших
  записей после dry-run. Основной worktree и его локальные изменения сохраняются.
- Код и публичная документация сохраняются локальными коммитами; tracker/spec/map/
  temp_plan остаются незакоммиченными до отдельной команды согласно действующему правилу.

## Проверки и границы

Использовать существующие публичные JSON/HTTP seams и структуру тестов пакетов.
gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, git diff --check;
независимые Standards/Spec review от фиксированного BASE с учётом рабочего diff.
Не добавлять entity-документы, expand, новые transport abstractions, live API,
.env, push/PR/VPS или новое поведение записи помимо SyncID.
