Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-access/issues/05-entity-metadata.md

# C05: Характеристики Variant: definitions и значения

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Разделить definition в variant metadata и value в Variant.characteristics. Чтение definitions;
Create по POST metadata/characteristics и Get definition по ID; отдельного Update в источнике
нет, не добавлять его. Это необходимое дополнение к Variant, не новый общий metadata CRUD
framework.

## Специальные проверки

Точные endpoints и list/envelope shapes; name/id/value неоднозначности из примеров
зафиксировать, unknown values, существующие definitions не создаются скрытым запросом.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/iterator/CI tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] README/guide/status отражают реальный scope, ограничения и неизвестные факты.
  Planning без команды не коммитится; один code commit по ticket.

## Источники

[characteristic](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_characteristic.md), [variant](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `bdfa7fa`, code commit `1d90849`.

| Метод | Endpoint | Примечание |
| --- | --- | --- |
| `ListCharacteristics` | GET entity/variant/metadata | массив characteristics, не Collection; meta.href и массив обязательны |
| `GetCharacteristic` | GET entity/variant/metadata/characteristics/{id} | |
| `CreateCharacteristic` | POST entity/variant/metadata/characteristics | только `{"name"}`; Update нет в источнике |

| Тип | Поля | Примечание |
| --- | --- | --- |
| `CharacteristicDefinition` | meta, id, name, type, required | readonly; type открыт |
| `Characteristic` (значение в Variant) | meta, id, name, value | запись значений — в C07 (id, иначе поиск по name) |

- Противоречие required: _characteristic.md — «всегда false», _variant.md — флаг обязательности; SDK отдаёт
  значение сервера, противоречие в godoc.
- GET entity/variant/metadata/characteristics как список не используется: источник не даёт формы ответа.
- Массовое создание (POST массива) не входит в v1 (batch исключён решением по X02).
- Неизвестные значения type и лишние поля читаются без ошибки; Raw у описаний нет — полный ответ
  метаданных доступен через metadata.Service.
- Review Standards: ListCharacteristics переведён на documenttest.HTTP (gzip, ошибки, null), фикстура
  метаданных — обоими codecs; проверка имени — как в остальных пакетах (`== ""`, без TrimSpace).
  Возврат среза вместо Collection и второй читатель metadata — осознанно, в godoc. validID в internal
  не переносится в этом ticket: дублирование общее для всех пакетов.
- Review Spec: добавлен тип значения Characteristic, правило id/name, противоречие required, отказ на
  ответ без characteristics; решения по второму endpoint и bulk — здесь.
- Проверки: go test -race ./entity/variant, go vet, gofmt — PASS (полный прогон — при коммите C06).
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
