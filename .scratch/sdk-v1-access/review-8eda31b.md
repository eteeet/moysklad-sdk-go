# Access повторное review — 2026-11-01

BASE: `33065182cb688fb0ef45bc3ceef7e9d414207c89`.
HEAD: `8eda31b6346b2b9f2e675e4806995ac4dc3d72ff`.
Scope: fix commits 251f638, d446807 и docs 8eda31b.
Метод: независимые Standards/Spec reviewers, прежние публичные воспроизведения.

## Standards

Подтверждённых замечаний нет. Оба исходных P2 закрыты: Permission сохраняет
объект с нестроковыми actions в Raw, AttributeList использует корректную
array/collection JSON-форму. Новых конкретных нарушений стандартов и
обоснованных smell-замечаний нет.

## Spec

Замечаний нет. Actions типизируется только для объектов из JSON-строк;
null, числа и вложенные объекты сохраняются без потерь. AttributeList.MarshalJSON
совместим с UnmarshalJSON; строки сохраняются после round-trip.
Регрессии охватывают обе формы списка, пустой массив, metadata-only collection,
Entity и вложенные CompanySettings, а также non-string/пустые/строковые права.
Конкретных новых регрессий или scope creep не найдено.

## Verification

- `/private/tmp/access_review_permissions.go` — null сохранён обоими codecs.
- `/private/tmp/access_spec_attributes.go` — before=1, after=1 обоими codecs.
- `go vet ./... ./testdata` — PASS.
- `go test -race -count=1 -p 1 ./...` — PASS.
- gofmt изменённых Go files и `git diff --check 3306518...8eda31b` — PASS.
- HEAD в начале и конце независимых проверок — 8eda31b; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Standards: 0 замечаний. Spec: 0 замечаний.
