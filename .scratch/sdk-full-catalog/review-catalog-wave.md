# Review f2f5f09..419e52a — 2026-11-03

Независимые Standards и Spec reviewers (субагенты), только чтение. go vet — PASS; go test -p 1 —
PASS после повторного прогона (одиночный локальный TCP-сбой).

## Standards

- P2 trackingcode: тесты без матрицы ответов — добавлена (plain/gzip/API/статус/тело/отмена).
- P2 trackingcode.Request: пустые cis/type уходили при изменении по id — Field+omitzero; проверка
  вложенных кодов.
- P3: тест Replace с путём в имени; приоритет ключей storebalances; godoc storebalances; группировка
  импортов скидок; «копейки» у SpecialPriceRequest.Value — исправлено.

## Spec

- P2 trackingcode (совпадает со Standards) — исправлено; godoc о противоречии id.
- P2 aggregatepack: positions без маркера readonly и в примере создания — пишутся в Create/Update.
- P2 discount: «raw для неизвестных» — общий тип хранит Raw.
- P2 images: пустой content — ошибка до запроса.
- P3 minimumStock: вывод в Create/Update, Get без fields; meta в StoreBalanceRequest — godoc.
- P3 product: изображения в теле товара — godoc ссылается на пакет images.
- P3 productfolders/productFolders: общий список читает оба ключа; пакеты видов скидок пишут и читают
  productFolders по примерам и фикстурам — решение записано в godoc и docs/catalog-extra.md.

Исправления — `2a76d13`; make check — PASS.
