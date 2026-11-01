Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-access/issues/01-employee.md

# A02: Role и права Employee

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Role Get/List/Create/Update и Employee GetSecurity/UpdateSecurity, чтение documented role
presets. Permissions — неоднородный JSON object по примерам, несмотря на Array(Object) в таблице
Role; known typed values + lossless unknown jsontext.Value, без map[string]any/float64 для
сырого JSON.

## Специальные проверки

Приоритет admin/App permissions проверяет сервер; не активировать Employee, не сбрасывать пароль
и не присваивать admin/owner скрытыми вызовами. Bool/scope/object task permissions, неизвестные
ключи/enum, exact security routes, отсутствие потерянных values.

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

[role](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md), [employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c7b17f3`, code commit `1c6f83a` (тип прав `model.Permissions` введён в A03 `842adec`).

JSON-контракт permissions (противоречие таблицы и примеров): объект «имя права → значение», где
значение — флаг (пользовательские права роли, `importData: true`), объект «действие → область»
(`supply: {"view":"ALL"}`, у задач — ключ `script` и действие `done`) или иная форма, которая
сохраняется в `Raw` без изменений; большие целые не проходят через float64, map[string]any нет.
Области открыты: NO, OWN, OWN_SHARED, OWN_GROUP, OWN_GROUP_SHARED, ALL, AUTHOR, ASSIGNEE,
AUTHOR_OR_ASSIGNEE. Ключи без нормализации регистра.

| Ресурс | Метод | Путь | Запись |
| --- | --- | --- | --- |
| Role | Get / List / Create / Update | entity/role[/{uuid}] | name (Create required), permissions |
| Role preset | Preset | entity/role/{owner,admin,individual,cashier,worker,…} | — (только meta) |
| Employee security | GetSecurity / UpdateSecurity | entity/employee/{id}/security | role{meta, permissions}, group, authorizedHosts, authorizedIpNetwork, authorizedIpNetmask |

- Не входят: Delete (X01), активация/деактивация и сброс пароля, login/isActive (меняются активацией),
  email — через Employee.Update. null источник не описывает: Field позволяет, решает сервер.
- Правила индивидуальной роли в godoc: переданные права выставляются, остальные выключаются (кроме
  view=ALL у currency, country, uom); без permissions сохраняются прежние права сотрудника.
- Критерии успеха: Role — id; Preset — meta.href; security — group (обязателен в ответе).
  Имя preset — открытое, проверяется только как безопасный сегмент пути.
- Права администратора/решения проверяет сервер; скрытых вызовов нет.
- Тесты: фикстуры custom_role, employee_role, employee_role_individual, employee_security обоими codecs,
  round-trip прав без потерь, golden Create/Update/security, documenttest.HTTP + LocalErrors для Get/
  Create/Update/Preset/GetSecurity/UpdateSecurity, listtest для List, матрица сброса ResponseInfo.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Standards: общий HTTP-контракт для security/preset, критерии успеха, role.Role вместо дубля,
  jsontext.Value в тесте. Review Spec: открытый preset, точные комментарии (email, authorized*),
  правила индивидуальной роли, ключ script, порядок полей при записи.

## Comments

Ревью 3306518 (P2): null внутри объекта прав терялся. Исправлено в `251f638`: Actions только для объекта из строк, иначе Raw без потерь.
