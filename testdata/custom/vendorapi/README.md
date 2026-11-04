# Vendor API 1.0

- Repository: https://github.com/moysklad/api-vendor-1.0-doc
- Commit: be7188e50e29b4ab9acb7e85d0cda08455d0b9fc
- Source: md/vendor_api/_moysklad_endpoints.md (SHA-256 `f23f79fe2c2d89d8634ef5d7408e9886a7dd086bbda3f5da1dde92df501983c6`),
  раздел «Получение контекста пользователя», пример ответа: `context.json`
- Source: md/vendor_api/_vendor_endpoints.md (SHA-256 `3be3ffe933fa06ae421604e396dd5c1df6bc3bfb046a4e777d345583af8fe7cb`),
  разделы «Активация решения на аккаунте» и «Деактивация решения на аккаунте»: тела `-d '…'`
  примеров curl (`activation-*`, `deactivation-*`) и примеры ответов (`response-*`)

| File | SHA-256 |
| --- | --- |
| `context.json` | `d670e44c49a7dc8fb44654048980cab6215df4d048355d66493f7dc00cc9c225` |
| `activation-install.json` | `e17a41ce7ed0722f94d073c69425a3247d569400c14f9a9ab9c94b9cef3ef050` |
| `activation-resume.json` | `03a4d996b069674aed5d58282fe6954b91b85839b25a9dba5390eba2ee2b7d9b` |
| `activation-permissions.json` | `f9da01a948bd444bae88b6eeff4cb77b9d9c2119e08ab0b73e22ce3e5dce0314` |
| `activation-tariffchanged.json` | `07ba70c793995e8d641dc2146acf45fb2495014ef1b6f93270b03f7838f363b1` |
| `activation-autoprolongation.json` | `c27a0f8675dca7ffaa8e2ff5118d51d95a951b4d50cdb642c7320a98691d196e` |
| `activation-fiscalapi.json` | `44bb486b6459bc7e6622db7e8873bd5d6d24015ffcf852e60e44789fbd718780` |
| `deactivation-uninstall.json` | `d2fc60a81de3aa898c175365f754d0b154eeeb9564f4514661595a589d51976a` |
| `deactivation-suspend.json` | `6588078c0fa86728acb9129cb021a1ce1c11b7aa3b89463f9540f9a4b7654791` |
| `response-settingsrequired.json` | `ec3104a627502692f39d6d7d98932e55fd46e24577db68043f88791c72534959` |
| `response-activated.json` | `5e9b584a929822cd2a07650fb645b4091fd0d2e1ed14d26a80d2bc73aaa01328` |

JSON скопирован без изменений. Подтверждает форму ответов из примеров документации,
а не поведение живого API.
