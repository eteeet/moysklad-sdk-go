# Примеры report/counterparty

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/reports/_report_counterparty.md (SHA-256 `fb2cbdaf5d49258cda49830dc207da33c801b658fc9db05556276ecbaf04bc70`)

Каталог назван `reportcounterparty`: `counterparty` занят фикстурами entity/counterparty.

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `list.json` | «Получить показатели контрагентов», ответ (L84–269) | `66bdcad4ea05676c439cd8546c2edb236f43d0f1cdd1a74d9704d15ef23a626a` |
| `selected.json` | «Выборочные показатели контрагентов», ответ (L322–413) | `4ce969e2f99b854558a227442c539802887677da1ff4b0b069865aa4b14faf72` |
| `get.json` | «Показатели контрагента», ответ (L438–470) | `d1b91d094a95c1c7c3eb67d02d12a5f944478fbdb37f1f8015fd1807b99e6274` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.

`selected.golden.json` — тело POST выборочных показателей, набрано вручную по примеру запроса
(L292–315) и сравнивается с телом, которое собирает SDK; SHA-256 не фиксируется.
