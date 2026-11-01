# Catalog повторное review — 2026-11-02

BASE: `2b91a89bb555634ae7822a4e5ea0a30e27f5a070`.
HEAD: `25a74857b58f06b7045cb5a8f79ff01dc9587499`.
Scope: 8bde787, a9ebfba, b1f3e9c, 25a7485.
Метод: независимые Standards/Spec reviewers и проверка исправлений root.

## Standards

Подтверждённых замечаний нет. ListLegacy проверяет итоговый beta-header с учётом
Config.Header и замены per-call options; конфликт отклоняется до отправки,
ResponseInfo предварительно сбрасывается. HeaderValues использует тот же helper
слияния, что отправка, и копирует значения. Новых обоснованных smell-замечаний нет.

## Spec

Замечаний нет. Исходные P2/P3 закрыты:

- ListLegacy отклоняет assortmentWithoutStock из client defaults и опций вызова.
  Явная замена конфликтующего default другим тегом учитывается. Другие теги не
  запрещены; обычный List по-прежнему выбирает новый режим.
- Публичная документация отделяет неизвестную семантику Service barcodes от
  описанного поведения Bundle/Variant.
- Исключение Quantity записано как решение пользователя в master spec,
  contracts.md, Answer C06 и публичной документации. Коллизия требований разрешена.
- Whitespace-исключение в .gitattributes ограничено verbatim assortment fixtures;
  сами файлы и SHA-256 provenance сохранены.

Конкретных новых регрессий и scope creep не обнаружено.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS; gofmt изменённых Go files — PASS.
- `git diff --check 2b91a89...25a7485` — PASS.
- `git diff --check 8eda31b...25a7485` — PASS для всей захватки.
- SHA-256 list.json и consignment.json совпадают с README; testdata diff пуст.
- HEAD в начале и конце проверки — 25a7485; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Все замечания предыдущего ревью Catalog закрыты в проверенном срезе.
Это результат проверки Catalog, не оставшихся захваток SDK v1.

Standards: 0 замечаний. Spec: 0 замечаний.
