# Примеры ответов print

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/dictionaries/_template.md (SHA-256 `be29a3ec143807ff721f35a321c39507ed76f1fbbe7141ca2c5e4495f44f7f67`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `embedded-list.json` | «Список стандартных шаблонов», ответ (L36–135) | `a3f09f8ee8953147826c97e1e8ca301c65f517d382af3e7a2d5b50df5d47799a` |
| `embedded-get.json` | «Отдельный стандартный шаблон», ответ (L160–170) | `ed68363e08a4c01a02c54600ebff1dd28e97f99eb676ff9704a46f547ec4618b` |
| `assortment-embedded-list.json` | «Список стандартных ценников и этикеток», ответ (L198–241) | `91d3b4341bafd261607fd7fe9cb1b493bb91d4b12d7791fff9051a4376bcb052` |
| `assortment-embedded-get.json` | «Отдельный стандартный ценник или этикетка», ответ (L267–277) | `261885857f0afc3740645b6851a42963bdde42a863494df6616cb76f7a9dc47c` |
| `custom-list.json` | «Список пользовательских шаблонов», ответ (L312–345) | `8a2a91e27de4bead13ca98e4349729a4d8f6445d296680cf3d371348fc6dc97d` |
| `custom-get.json` | «Отдельный пользовательский шаблон», ответ (L370–380) | `c827dfb2945969f1d7d4b8c4a3d26492e6d4be673c2266044ae647bdbcc44955` |
| `assortment-custom-get.json` | «Отдельный пользовательский ценник или этикетка», ответ (L477–487) | `6810006f296acf402f1b514f07bbae1d4f9e9c5b206ff823f13c3ef9aabdf4a9` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.
Ответ «Список пользовательских ценников и этикеток» (L408–451) побайтно совпадает с
`assortment-embedded-list.json` (meta.type `embeddedtemplate`) — ошибка копирования в источнике, отдельной
фикстурой не взят.

`*.golden.json` — независимые тела запросов печати по примерам md/documents/_print.md и
md/dictionaries/_labels.md (не копии: пример передаётся в shell-строке).
