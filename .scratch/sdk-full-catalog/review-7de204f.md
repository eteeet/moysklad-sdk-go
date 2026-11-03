# Повторное review исправлений catalog — 2026-11-03

BASE: `2a76d132d0f742b6ae764c9532e32dfe6d7e167d`.
HEAD: `7de204f94129feacde71b23b0ae42fa8a7a528b4`.
Scope: один fix commit, 17 files; закрытие review-2a76d13 и пояснения противоречий.
Метод: независимые Standards/Spec reviewers, проверка root, полный make check.

## Standards

Прежний S1 закрыт: money_test.go проверяет независимые wire JSON 0/0.0/дробь
обоими codecs для Level.Amount и SpecialPrice.Value, а также запись нуля и дроби.

### S1. P3 — не хватает HTTP regression cases исправленного поведения

- `entity/product/fields_test.go:50` проверяет PackRequest.ID только через Marshal.
  Предыдущий review требовал сохранение ID в JSON и HTTP; TestServiceUpdate не
  содержит запроса с существующей упаковкой. Добавить вызов публичного Update
  с точным PUT body, содержащим packs[].id.
- `entity/consignment/consignment_test.go:71` проверяет expiryDate-only только
  через Marshal. HTTP Create использует full с непустым Label, поэтому новая
  успешная ветка hasDate не проверена. Возврат прежней проверки обязательного
  Label не сломает эти тесты. Нужен Service.Create с ExpiryDate без Label,
  проверкой body и ровно одной отправки.

Это пробел regression coverage, а не подтверждённая ошибка production-кода.
Дублирование небольших alias decoders само по себе не требует новой абстракции.
Иных существенных Standards findings в изменениях нет.

## Spec

Замечаний нет. Предыдущие P1/P2 закрыты:

- Все три concrete-модели скидок читают оба spelling productFolders/productfolders,
  canonical ключ имеет приоритет; tests проверяют оба codecs и оба ключа.
- Optional PackRequest.ID выражает обновление существующей упаковки; неверное
  утверждение о неизбежном пересоздании удалено из godoc и guide.
- Consignment допускает expiryDate-only без label в JSON. Unset, null и нулевая
  дата не удовлетворяют hasDate; непустого Label достаточно. Update не требует
  one-of, поскольку SDK не знает сохранённого состояния сущности.
- Документация раскрывает три противоречащих требования Consignment, отмечает
  неподтверждённость сервером формы specialPrice и различает идентификацию
  объекта по id/meta и изменение readonly-свойства.

Новых подтверждённых ошибок реализации в diff не найдено. Server acceptance
спорных форм живыми запросами не проверялся.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., diff-check полного tracked HEAD и working tree.
- HEAD оставался 7de204f. Code, planning, tracker statuses и commits reviewer
  не менял; создан только этот отчёт. .DS_Store не читались/не изменялись.
- Live API не вызывался.

Standards: 1 finding P3. Spec: 0 findings.
Предыдущие P2 закрыты; оставшийся P3 не блокирует начало независимой следующей захватки.

## Answer

P3 закрыт в `d21d84d`: HTTP-тесты Update товара с packs[].id и Create партии только с expiryDate
(точное тело, одна отправка). make check — PASS.
