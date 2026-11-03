# Review реализации catalog и решений по документации — 2026-11-03

BASE: `f2f5f098e0440ea118ca04fd34eb27984e264e21`.
HEAD: `2a76d132d0f742b6ae764c9532e32dfe6d7e167d`.
Scope: 80 файлов, G01–G11 и S10 Consignment; G12/G13 намеренно ждут O01/O04.
Метод: независимые Standards/Spec reviewers, проверка root и make check.
Предыдущий внутренний review-catalog-wave.md и исправления 2a76d13 учтены.

## Standards

### S1. P3 — неполная матрица денежных JSON-контрактов

`entity/discount/discount_test.go:49`,
`entity/accumulationdiscount/accumulationdiscount_test.go:48`,
`entity/specialpricediscount/specialpricediscount_test.go:48`.

Новые денежные Level.Amount и SpecialPrice.Value читаются из integer fixtures;
нет независимых wire cases 0.0 и дроби обоими codecs. Marshal model.Value(0.0)
пишет 0 и не проверяет приём JSON 0.0. docs/coverage.md:22 прямо требует для денег
0, 0.0 и дробь. Добавить адресные read/write cases этих денежных форм.

Других подтверждённых Standards findings текущего HEAD нет. do/validID намеренно
приватные согласно coverage.md; дополнительных обоснованных baseline smells нет.

Процесс, уже исправленный в HEAD: commit 419e52a прямо фиксирует, что snapshot
9f94905 был закоммичен с падающими общими batch tests. Это нарушило правило проверки
каждого commit в .scratch/sdk-v1/spec.md:137–139, но текущий HEAD исправлен.
Отдельной runtime ошибкой HEAD этот исторический факт не считается.

## Spec

### P1. P2 — concrete-модели скидок теряют productfolders при JSON v2

`entity/accumulationdiscount/model.go:20`,
`entity/personaldiscount/model.go:17`,
`entity/specialpricediscount/model.go:22` принимают только productFolders.
Таблица источника использует productfolders. Общий discount.Discount поддерживает
оба ключа; конкретные типы — нет. JSON v1 сопоставляет регистр нечувствительно,
JSON v2, используемый transport, пропускает lowercase поле без ошибки.

Root воспроизвёл: один JSON с productfolders даёт длину ProductFolders=1 у всех
трёх моделей в v1 и 0 в v2. Нужен одинаковый alias decoder у всех concrete-моделей,
с приоритетом canonical productFolders и tests обоими codecs. Выбранный canonical
ключ записи менять ради этого не требуется.

Источник: md/dictionaries/_discount.md:23,38; требование spec — оба JSON codecs,
ссылки групп входят в G07–G09. Наличие lowercase в живых ответах API не проверялось;
дефект — потеря документированной формы и различие codecs.

### Дополнительная проверка root

### P2. P2 — PackRequest не позволяет адресовать существующую упаковку

`entity/product/request.go:116–120`: PackRequest не содержит ID. В PUT-примере
изменения товара источник передаёт packs[0].id и сохраняет этот ID в ответе
(_product.md:3042,3141–3152,3356–3369). SDK не может выразить этот запрос: можно
передать поля упаковки, но нельзя выбрать конкретную существующую упаковку по ID.

Это тот же случай идентификации объекта, что TrackingCode.id и StoreBalance.meta,
а не изменение идентификатора. Добавить optional ID для обновления существующих
упаковок, оставить отсутствие ID для новых; проверить сохранение ID в JSON и HTTP.
Фраза в docs/catalog-extra.md и G01 Answer о неизбежном пересоздании упаковок
не является установленной гарантией источника: источник подтверждает замену
коллекции, но отдельно показывает сохранение существующего ID. Исправить объяснение.

## Комментарий к шести решениям

1. **storebalances/storeBalances — нормально.** Принимать оба документированных
   spelling и иметь детерминированный приоритет. Это уже реализовано и проверено.
2. **productfolders/productFolders — canonical запись обоснована примерами.**
   Чтение нужно сделать одинаково терпимым во всех моделях, см. P1. Запись обоих
   ключей одновременно не нужна.
3. **specialPrice — осознанное, но не проверенное на API решение.** Таблица даёт
   priceType/value, Create-пример — meta типа цены. Текущая форма соответствует
   G09; нельзя выдавать её за подтверждённый сервером write contract. Сам факт
   формы read-объекта не доказывает тождественную форму write-объекта. При дальнейшей
   проверке нужен фактический write source или отдельно разрешённая live-проверка.
4. **expiryDate — описание противоречия неполное.** У _consignment.md три версии:
   таблица требует expiryDate; специальный Create-текст требует label + assortment;
   вводный текст (:12) допускает label или expiryDate. Optional expiryDate обоснован
   Create-текстом, но текущая безусловная проверка Label в service.go:42–43 выбирает
   одну из конфликтующих трактовок и блокирует expiryDate-only до API. Это следует
   явно включить в принятое решение/ограничение, а не утверждать однозначный контракт.
   Что реально принимает сервер, в review не проверялось.
5. **StoreBalance.meta — идентификация строки.** Передача существующей meta для
   выбора изменяемой строки не делает метаданные редактируемым бизнес-полем.
   Это не равно случаю Slot.Zone, где меняется связь с зоной. Отдельное исключение
   из readonly ради идентификации не требуется; семантику нужно закрепить в контракте.
6. **TrackingCode.id — идентификация кода.** Сам метод источника явно требует
   идентификатор обновляемого кода (_tracking_code.md:102–103). Нельзя назначать или
   менять серверный ID, но передавать его для выбора кода необходимо. Отдельное
   пользовательское исключение здесь также не требуется.

Общее уточнение правила: readonly запрещает изменять значение поля. Передача id/meta
для адресации существующего объекта допустима, когда именно это задаёт метод API.
Не нужно превращать каждый такой идентификатор в отдельный запрос разрешения.
Это уточнение относится и к P2 про PackRequest.ID; не отменяет приоритет readonly
для обычных изменяемых свойств.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., diff-check полного tracked HEAD и working tree.
- `go run -race /private/tmp/review-catalog-wave-aliases.go` — воспроизведено root
  без сети: v1 accumulation/personal/special=1; v2=0 для одного lowercase JSON.
- Источники — локальный архив 9a811cef; ранее проверенный официальный compare
  к 0e1b3d61 меняет только code-fence разметку print/labels, не эти контракты.
- HEAD оставался 2a76d13; code, planning, tracker statuses и commits reviewer не менял.
  Созданы только этот отчёт и временный repro; .DS_Store не читались/не изменялись.
- Live API не вызывался, write-запросов в реальный аккаунт не было.

Standards: 1 finding P3. Spec: 2 findings P2.

## Answer

Все замечания приняты, исправлено в `7de204f`, make check — PASS.

- S1: кейсы 0, 0.0 и дробь для Level.Amount и SpecialPrice.Value обоими кодеками.
- P1: виды скидок читают productfolders и productFolders, приоритет у productFolders; тесты обоими кодеками.
- P2: PackRequest.ID адресует существующую упаковку; описание о пересоздании исправлено (godoc, docs, Answer G01).
- Решение 4 (expiryDate): SDK требует assortment и хотя бы одно из label и expiryDate; три версии источника
  записаны в godoc и docs.
- Решения 3, 5, 6: форма specialPrice помечена как не проверенная сервером; meta storebalances и id кода
  маркировки описаны как адресация, без отдельного исключения.
- Процесс: 9f94905 закоммичен с падающим общим тестом — причина в tail, съевшем код выхода; дальше коммиты
  проверяются без конвейера.
