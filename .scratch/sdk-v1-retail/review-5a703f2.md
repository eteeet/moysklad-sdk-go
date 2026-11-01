# Retail review — 2026-11-02

BASE: `2f21bcf4d7a65e56949e5eed8c1192a2af057f21`.
HEAD: `5a703f2c24d77ee19e58853c1d27d42076aa6ef2`.
Scope: E01–E06, общий test commit и завершающий docs commit; 8 commits, 49 files.
Метод: независимые Standards/Spec reviewers, сверка с pinned первичными источниками.

## Standards

Замечаний нет.

- ADR 0001 соблюдён: RequestOption, сброс ResponseInfo до validation,
  общий Client.Send, сохранение errors.Is/As.
- Readonly отделены от requests; optional поля используют Field/omitzero.
  ChequeRequest сохраняет явные 0 и пустые строки, GiftCardRequest не пишет paymentSum.
- JSON проверяется обоими codecs; HTTP — настоящим transport и httptest.
- Raw environment/state, формы Links и отдельные nested models соответствуют freeze.
- Существенных code smells с учётом repo override для do/validID не обнаружено.

## Spec

Подтверждённых замечаний нет.

- Required Create, собственные Update exclusions, readonly, денежные суммы и баллы
  сверены для шести ресурсов с field matrix и источником.
- RetailDemand требует retailShift; SDK не открывает смену и не получает её скрыто.
- RetailSalesReturn различает based/unbased тела, demand ссылается на розничную
  продажу; явные cashSum/noCashSum=0 сохраняются. Бизнес-ограничения проверяет сервер.
- BonusProgram сохраняет обязательные false/0. BonusTransaction не позволяет
  менять transactionType и писать transactionStatus/categoryType; отказ токена
  возвращается без скрытой смены credentials.
- Отдельных cashier/position mutations, фискализации, расчёта бонусов и скрытых
  запросов не добавлено. Явные поля конфигурации и Cheque остаются input caller.

Неоднозначность источника, не признанная дефектом: RetailDemand.RetailStore не
входит в write matrix E03. Общий текст Update разрешает поля без readonly,
но отдельно изменение retailStore не описано. Поддержка этого поля записи
не обещается текущим freeze; вывод о поведении API не делался.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS.
- gofmt изменённых Go files и `git diff --check 2f21bcf...5a703f2` — PASS.
- Источник: api-remap-1.2-doc revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd,
  _retailstore.md, _retailshift.md, _retaildemand.md, _retail_sales_return.md,
  _bonus_program.md, _bonus_operation.md и общие разделы.
- HEAD оставался 5a703f2; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Результат относится к проверенному профилю Retail E01–E06, а не всей SDK v1.
Standards: 0 замечаний. Spec: 0 замечаний.
