# Discount, Cost planning и дублирующий тест

Дата: 2026-11-01. BASE eba6f75154929790c5b9d2cb1454dd80783b14db.
Commit 59b16e73a7395ffdc40c2fadb8caeee381de9a35.
Diff: git diff eba6f75...59b16e7 плюс согласованные Cost-правки рабочего planning.

## Scope и результат

- Coverage/guide: Discount — процент скидки/наценки; Price/Overhead в копейках.
- S02: общий Cost float64 и Cost Field[float64] omitzero, только для услуг;
  применимость проверяет сервер без Get. Production Cost этой задачей не добавлялся.
  Связанные specs согласованы; будущие JSON/HTTP cases и primary sources закреплены.
- Удалён один absent fields из TestPositionOptionalFieldsJSON; в TestPositionJSON
  этот case и остальные уникальные проверки сохранены.

## Проверки и параллельная работа

Во время исполнения внешний процесс добавил 11 document packages, internal и поля
Cost/Shipped/InTransit в positions/model.go/request.go. Эти изменения не входят
в scope и не включались в commit 59b16e7. Общий race-прогон завершился FAIL в positions:
новые нулевые JSON-поля расходятся с прежними golden. Следующий go vet в цепочке
не запускался. Чужие поля и ожидания не правились.

Для независимой проверки создан snapshot eba6f75 без рабочего diff через git archive
в /private/tmp/moysklad-doc-fixes.YNlCMm; применено только удаление одного case.
Полный go test -race -count=1 ./... и go vet ./... ./testdata там PASS.
Gofmt изменённого test-файла и git diff --check PASS.
Точечный go test -race -count=1 ./positions -run '^TestPositionOptionalFieldsJSON$'
в общем рабочем каталоге также PASS.
Это подтверждение own delta, а не зелёного состояния общего рабочего каталога.

## Standards

Независимый read-only reviewer unify_review_standards: 0 замечаний.
Discount отделён от денег; удалён только дублирующий case; Cost в общем типе и Field
соответствует контракту, серверная применимость и будущие tests/sources заданы.
Нарушений правил и существенных smells нет. Параллельный Go-код не проверялся.

## Spec

Независимый read-only reviewer unify_review_spec: 0 замечаний.
Все три согласованные темы выполнены, scope creep нет. PASS относится к изолированному
snapshot; текущий main зелёным не объявляется.

Итого: Standards 0, Spec 0. Planning/tracker не коммитились; worktree не затронуты.
