Status: resolved
Причина: работа закрыта: `financial-documents/spec.md` resolved. Закрыто 2026-11-14 при сверке здоровья документов.
Stage: DRAFT
Type: research

# Финансовая четвёрка: проверенные источники и открытые решения

Проверка 2026-11-01. [master HEAD](https://api.github.com/repos/moysklad/api-remap-1.2-doc/commits/master)
на момент проверки — 9a811cef3dc4a3f3a46f070e60452361a06bf7cd (commit от 2026-10-26),
совпадает с revision предыдущего исследования. Корень GitHub по умолчанию показывал
pre_release; контракт взят с master и закреплён полным SHA. Никакие запросы к живому API
аккаунта не выполнялись. Доказательства — official api-remap-1.2-doc, не OpenAPI/fixtures.

## Проверенные факты

| Факт | Первичный источник |
| --- | --- |
| PaymentIn: organization+agent достаточны для минимального Create; sum, name, moment не required Create; incomingDate/incomingNumber и account-ссылки есть в таблице | [PaymentIn, таблица и Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L5-L43), [минимум Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L326-L363) |
| PaymentOut: required organization+agent+expenseItem; agent допускает Counterparty/Employee/Organization, agentAccount — счет Counterparty/Organization; accrualDate и noClosingDocs есть | [PaymentOut, таблица](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L5-L43), [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L426-L477) |
| CashIn: required organization+agent; account-ссылки и expenseItem отсутствуют в таблице; factureOut и incoming operations присутствуют | [CashIn, таблица/связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md#L5-L54), [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md#L391-L425) |
| CashOut: required organization+agent+expenseItem; agent описан только Counterparty; accrualDate, noClosingDocs, factureIn и outgoing operations присутствуют | [CashOut, таблица/связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L5-L57), [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L450-L493) |
| ExpenseItem реально передан в минимальных outgoing примерах; автоподстановка/default для него не описаны | [PaymentOut minimal](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L426-L477), [CashOut minimal](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L450-L493) |
| Operations — массив элементов с непосредственными meta и необязательным при записи linkedSum; вложенного operation/reference wrapper нет | [Привязка финансовых документов](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L507), [PaymentIn Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L1491-L1499) |
| DocumentRate содержит currency, value Float; currency обязательна в ответе, value появляется при курсе, отличном от 1 | [Валюта в документах](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L689-L699) |
| rate:null — смена на валюту по умолчанию с серверным пересчётом; contract:null — документированный сброс object | [Null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187) |
| Банковский organizationAccount в иной валюте требует согласованного rate, ошибка 22004; для cashin/cashout проверка по счету не применяется | [PaymentIn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L62-L64), [PaymentOut](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L63-L65) |
| Incoming/outgoing поддерживают разные наборы operations; noClosingDocs=true с непустым operations запрещён только в outgoing таблицах | [Incoming operations](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L45-L58), [Outgoing operations/flag](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L24-L59) |
| Общий API описывает syncId как идентификатор повторного Create; после заполнения его нельзя изменять | [syncId и обновление](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1144-L1155), таблицы четырёх документов |

## Противоречия и пределы доказательств

1. CashIn.Update и CashOut.Update требуют account-поля одновременно со сменой сторон,
   хотя эти поля не описаны в их таблицах, а банковская валютная проверка для ордеров
   прямо не применяется. Нельзя признать банковские поля cash-документов поддержанными
   по одной повторённой фразе. Proposal: не включать их; server validation передавать как есть.
   Источники: [CashIn.Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md#L1523-L1531),
   [CashOut.Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L1694-L1702).
2. PaymentOut.agent явно polymorphic, CashOut.agent назван Counterparty. По сходству
   исходящих документов нельзя объявить остальные типы CashOut подтверждёнными.
   Proposal: Reference остаётся открытым и не валидируется локально; docs описывают только
   подтверждённые варианты и неопределённость CashOut.
3. RetailShift в списке типов написан с заглавной S, meta.type примеров — retailshift.
   Proposal: не нормализовать строку и не вводить закрытый enum.
4. Связь operations описана как добавление/привязка; Update говорит передать все meta,
   но не закрепляет достаточно чётко замену всей коллекции и поведение []/null.
   Proposal: SDK передаёт точный массив, не делает merge/GET и не обещает очистку по [].
   В [общем примере](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L605)
   нельзя доказать отсутствие ранее существовавших связей.
5. Float в таблицах не определяет decimal-арифметику и не обещает сохранение исходного
   текстового вида числа. Proposal: float64 как существующий Price, без новых расчетов,
   масштабирования и округления; нужны проверки дроби и 0 в обоих JSON-кодеках.
   Поля sum/vatSum/linkedSum не переводить в int только из-за целых чисел в примерах.
6. Необязательность поля при Create не доказывает конкретный default. Для name/moment/
   externalCode минимальные примеры показывают заполнение сервером; конкретные значения
   и defaults applicable/sum/vatSum SDK не задаёт. ExpenseItem нельзя опускать по аналогии
   с другими необязательными ссылками.
7. Required-response метки не гарантируют присутствие во всех старых примерах. Proposal:
   разреженный ответ читается с нулевыми значениями, обязательность кроме ID не валидируется.
8. Общее правило null условно для optional object и UI-пустоты. Явные contract:null и
   rate:null подтверждены; null прочих scalar/collections/required ссылок не объявляется
   допустимым сбросом. Field способен передать null, но не превращает это в гарантию API.
9. SyncId в API можно заполнить до первого непустого значения, а не строго только при Create.
   Уточнённый proposal: Field[string] в Create/Update, без предварительного чтения;
   проверка запрета изменения уже заполненного syncId остаётся серверу.

## Точки triage

- Утвердить поимённый write scope и намеренное исключение files.
- Утвердить float64 и отсутствие новой финансовой арифметики.
- Утвердить FinanceOperation/FinanceOperationRequest и request linkedSum через Field.
- Утвердить минимальную validation обязательных Create Ref.Meta.Href и отсутствие локальных
  бизнес-валидаторов; согласовать null-policy и SyncID в Create/Update.
- Утвердить предложения для противоречий cash account-полей и CashOut.agent; не изобретать
  гарантии замещения/очистки operations. Общий DocumentRate согласовать в document-common.

Все эти пункты отражены в DRAFT, а не скрыты за утверждённым статусом.
