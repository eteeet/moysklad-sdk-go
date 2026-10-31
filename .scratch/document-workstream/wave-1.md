Status: resolved
Stage: COMPLETED

# Первая волна: C01/C02/F01/P01/P02

Пользователь 2026-11-01: «внеси правки и закоммить все изменения перед следующим этапом.
потом запускай следующий этап». Команда следует согласованному отзыву: первая волна
из пяти tickets, три потока, независимые review, локальные коммиты. Entity и P03 не запускать.

Этап завершён: planning BASE a280428, реализация до 1b28603, пять tickets resolved.
[Итоги и отдельные отчёты Standards/Spec](review-wave-1.md). Ниже сохранён исторический handoff,
повторно запускать его не требуется. Следующая группа отдельно согласуется пользователем.

## Перед стартом

Сначала сохранить все текущие изменения spec/tickets/docs/temp_plan отдельным planning
коммитом в dev. Его HEAD — фиксированный BASE реализации и итоговых двух review.
Только после чистого status создавать рабочие ветки/worktree и менять Go-код.
Не читать .env, не обращаться к живому API, не выполнять push/PR/VPS.

## Ownership

| Поток | Tickets | Разрешённые новые файлы |
| --- | --- | --- |
| model JSON | C01 → C02 → F01; отдельные коммиты | model/document_rate.go, model/document_rate_test.go; model/document_overhead.go, model/document_overhead_test.go; model/finance_operation.go, model/finance_operation_test.go; соответствующие три issue-файла |
| positions JSON | P01 | positions/doc.go, positions/model.go, positions/request.go, positions/model_test.go, positions/request_test.go; positions/testdata/json-*.json при необходимости; issue P01 |
| positions reader | P02 | positions/service.go, positions/service_test.go, positions/pagination.go; positions/testdata/page-*.json при необходимости; issue P02 |

P01/P02 создают один новый package, но не общие файлы или test helpers. P02 тестирует
пользовательский T, поэтому не зависит от моделей P01. Имена test helpers префиксовать
по своему срезу, не вводить общий helper-файл параллельно. Если нужен новый путь — согласовать
с интегратором до записи. Existing model helpers/transport/entity не рефакторить.

Все потоки не меняют README, docs/coverage.md, docs/temp_plan.md, общие spec и карты.
Согласованная документация публичного interface — GoDoc в принадлежащих им файлах;
интегратор после слияния сделает общий docs commit. Отдельные handoff-заметки в issue
допустимы, общие инструкции не дублировать.

## Выполнение и проверка

1. Прочитать AGENTS.md, CLAUDE.md, docs/agents/, общий glossary и свои spec/ticket.
   Прочитать tdd SKILL.md, tests.md и mocking.md полностью. Согласованные seams:
   public JSON обоими codecs; для P02 public GetPage через httptest.Server + transport.Client.
2. Установить свой Status: claimed. Делать red→green вертикально, фиксировать результат
   первого падающего теста до добавления реализации; не писать все тесты горизонтальным слоем.
3. Проверить все acceptance criteria, gofmt, targeted tests, go test -race -count=1 ./...,
   go vet ./... ./testdata и git diff --check. Никаких live-запросов, retry или expand.
4. Локально коммитить только принадлежащие потоку код/тесты и claim своего issue.
   В Answer отразить фактическую реализацию/проверки, но Status оставлять claimed до review.
   Передать интегратору SHA, red/green evidence, проверки и замеченные риски.
5. Интегратор объединяет ветки без перетирания файлов, запускает полный race/vet,
   выполняет параллельные независимые Standards и Spec review от planning BASE.
   Reviewer код не меняет; исправления — основной implementer с последующей проверкой.
6. После чистых review интегратор ставит пять tickets resolved, фиксирует результаты
   и обновляет coverage/README/map/temp_plan одним отдельным docs commit. Затем остановка.

## Контроль спорных мест

- Деньги float64 в копейках независимо от Int/Float таблицы, включая будущие cost/overhead;
  тесты 0/0.0/дробь. Не требовать точности денежных integer >2^53.
- DocumentRate.Value *float64 с omitzero: nil в присутствующем rate означает курс 1.
  Комментарий: «nil — курс 1: сервер не передаёт value, равный 1». Не подставлять 1
  в структуру/JSON; отсутствие всего rate или unset request не означает тот же default.
- Общий reader сохраняет открытый documentType и T any; в будущих entity packages
  будет типизированный GetPositions. Эти entity tickets зависят от P02, но сюда не входят.
- Counterparty/Organization SyncID Update — TODO в coverage, не реализация этой волны.
- P03 Movement остаётся needs-triage; числовые исправления в его DRAFT не запускают его.
