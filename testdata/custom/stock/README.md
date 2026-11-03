# Примеры ответов report/stock

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/reports/_report_stock.md (SHA-256 `da0fc332242adb8bae8deb4d4e3a9476533171f75a2ad3a73fb4176cc9ca0af2`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `all.json` | «Получить Расширенный отчет об остатках», ответ (L208–674) | `0178c564cf7d801378bf7624e90b3e743654ce36cbc8647763e26f4ac357f224` |
| `all_current.json` | «Получить Краткий отчет об остатках», первый ответ /all/current (L774–777) | `161a69788e7fd568b23643f69e97a055352feb7393cdf82377980aeaad7fab23` |
| `bystore_current.json` | там же, /bystore/current?stockType=freeStock (L834–839) | `98755edaa85a59946670ae930600f027f9bfd72610b0804ee53f6ef8683ed0e8` |
| `bystore_current_recalculate.json` | там же, /bystore/current?withRecalculate=true (L918–922) | `9e7be6d49a00419e3ba49e0a8e8e8041324c1bfda3c95f64cf32e1dde4378d18` |
| `byslot_current.json` | «Получить Краткий отчет об остатках по ячейкам», ответ (L961–980) | `f1156633e9af58a50a8fad8a3edd473db46cc5ceba30a3e5e3cfadacc7cfe4d1` |
| `bystore.json` | «Получить Остатки по складам», ответ (L1115–1372) | `9ff4fb7cb68ed4a4756520edc601308550469bc0a04a8153c08968c9695b2298` |
| `byoperation.json` | «Получить Остатки по документу», ответ (L1441–1553) | `1530734130165c752123c59fa191af128581ab0677e669d394ad90087cda9eff` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.
