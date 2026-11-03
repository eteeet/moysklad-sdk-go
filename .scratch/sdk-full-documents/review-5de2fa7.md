# Review реализации документов — 2026-11-03

BASE: `d21d84d44c10078ff6e0e28a966d82d065cefc88`.
HEAD: `5de2fa76148380386a7dec4cce9f81c2ff3b92c0`.
Scope: O01–O06, O08–O20, 250 files; O07 намеренно остаётся needs-info.
Метод: независимые Standards/Spec reviewers, дополнительная проверка root,
полный make check и отдельная перепроверка упавших пакетов.
Предыдущие review-wave.md и исправления 2e9e03d/3b62b18/341a4fd учтены.

## Standards

### S1. P2 — неполная HTTP-проверка wizard

`wizard/wizard_test.go:92–94`: TestFill проверяет лишь non-nil/error/число запросов.
Возврат пустого Result вместо разобранного HTTP-ответа не сломает успешный test.
Нет gzip-сценария и проверки headers/ResponseInfo успешного вызова. TestJSON
проверяет модели отдельно и не подтверждает результат публичного Fill.

Правило: sdk-full-documents/spec.md:261–264 — exact headers, plain/gzip и
фактический server response. Добавить assertions рассчитанных значений и Raw
из результата HTTP-вызова и применимые transport cases. Это test coverage gap,
а не подтверждённая ошибка production-кода.

### S2. P3 — docs обещают отсутствующие способы использования

- `docs/documents-extra.md:31–32` предлагает доработать результат New и передать
  его в Create. New возвращает read model, Create принимает отдельный CreateRequest;
  прямого преобразования нет. Описать перенос выбранных полей в CreateRequest.
- `docs/documents.md:39–41` обещает alias PositionRequest у Prepayment/PrepaymentReturn,
  хотя эти пакеты экспортируют только Position и поддерживают чтение строк.

Правило: публичная документация должна отражать фактический interface.

Других подтверждённых Standards нарушений и существенных baseline smells нет.
Приватные do/validID соответствуют явному решению coverage.md.

## Spec

Замечаний нет в проверенном scope. Основные таблицы 13 новых read-моделей,
required Create fields, собственные формы строк inventory/pricelist/payroll/
commissionreportin/out, ID существующих строк, оба набора позиций commissionreportin
и meta.type pricelistrow сверены с источником.

Четыре группы readonly-исключений явно разрешены пользователем; вопрос не
переоткрывается. Реализация соответствует зафиксированной Create/Update matrix:
shared commissionreportout и priceType pricelist остаются Create-only.

Источник payroll задаёт employee/baseSalary/pieceworkSalary, а старый текст ticket
говорит employee/sum/type: реализация следует источнику, расхождение записано в Answer.
Inventory/pricelist пишут {"rows":[…]} по примерам, решение отражено в docs/godoc.

Root дополнительно проверил positions.Service mutations, nested.DeleteRows,
write.Template/Trash, wizard и общие связи. Подтверждённых ошибок реализации
в этих механизмах не найдено.

## Verification

- `make check` на HEAD с CHECK_P=2 — **FAIL** на двух пакетах:
  - FactureOut: TestFactureOutHTTP/create/no_content, calls=0 вместо 1;
  - Group: TestServiceHTTP/get/html_error, dial 127.0.0.1:56903 — connect timeout.
- Остальные пакеты полного race-прогона прошли, включая wizard и positions.
  gofmt и go vet ./... ./testdata в начале make check — PASS.
- `go test -race -count=1 -p 1 ./entity/factureout ./entity/group` — PASS при
  отдельной перепроверке. Это не превращает первый полный прогон в PASS и не
  устраняет причину нестабильности. Для FactureOut непосредственно зафиксирован
  нулевой счётчик, а не текст TCP-ошибки; общую причину двух сбоев review не доказывает.
- Поскольку make check остановился до diff-check, root отдельно выполнил
  diff-check полного tracked HEAD, working tree и диапазона BASE...HEAD — PASS.
- Первый неуспешный прогон сохранён как результат проверки. Диагностика TCP
  нестабильности остаётся отдельной задачей, исправления в review не вносились.
- Источник — локальный архив 9a811cef, контракты совпадают с 0e1b3d61
  (ранее проверенный официальный compare меняет лишь fences labels/print).
- HEAD остался 5de2fa7. Code, planning, tracker statuses и commits reviewer не менял.
  Создан только этот отчёт; .DS_Store не читались/не изменялись; live API не вызывался.

## Следующая захватка

Начинать производство можно: O01/O02/O04 реализованы и resolved, новых документных
blockers для M05–M08 не найдено. Внутренний порядок производства сохраняется:
M01/M02 → M03 → M04 → M05–M07 → M08. Field/method freeze соответствующих tickets
выполняется обычным порядком до кода; review не меняет их стадии автоматически.

G12/G13 также разблокированы O01/O04. Отчёты/печать/интеграции можно вести
независимо с одним владельцем общих transport/positions/facade файлов.
O07 остаётся needs-info и эти работы не блокирует.

Standards: 2 findings, максимум P2. Spec: 0 findings.
Следующая независимая работа допустима; текущей волне ещё нужны тестовые/docs
правки и устранение нестабильности полного проверочного прогона.

## Answer

Исправлено в коммите после 5de2fa7 (`fix: замечания review 5de2fa7`):
- S1: TestFill сверяет поля Result и Raw с ответом сервера, есть случай gzip, проверяются
  заголовки запроса, кодирование `action` (`%2C`) и ResponseInfo. Мутация «Fill возвращает
  пустой Result» ловится в случаях ok и gzip.
- S2: «доработать и передать в Create» убрано из docs/documents-extra.md и godoc всех 22 New;
  docs/documents.md больше не обещает PositionRequest у Prepayment/PrepaymentReturn.
- Verification: calls=0 в FactureOut create/no_content — сбой соединения. Ошибка подключения
  удовлетворяет ожиданию ошибки в случае no content, падает только счётчик вызовов. Теперь
  documenttest.HTTP печатает ошибку вместе с calls. Причина нестабильности (исчерпание
  локальных портов) — отдельная задача.
- make check CHECK_P=1 на новом HEAD — PASS.
- Нестабильность прогона устранена в c939971: тестовые HTTP-серверы — internal/memnet (net.Pipe),
  без loopback TCP. Причина — поток новых loopback-соединений на macOS (~1 сбой на 3000),
  воспроизведена стендом на stdlib; исчерпание портов, очередь listen и сторона TIME_WAIT
  опровергнуты. go test -p 4: было 9–12 таймаутов за прогон, стало 0 (3 прогона);
  make check с P=2 и P=4 — PASS.
