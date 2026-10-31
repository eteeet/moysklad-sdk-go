# Organization: request golden

`create.golden.json` и `update.golden.json` — синтетические ожидаемые тела запросов,
заданные независимо от сериализатора SDK по согласованному контракту
[Organization](../../../docs/organization.md).
Они не являются ответами живого API или копиями официальных примеров.

Golden покрывают весь согласованный набор записи, обе формы адресов одновременно,
все десять вложенных адресных полей, false/true, строковые коды, даты,
references, `attributes` с целым числом больше 2^53 и явными `value:null`/`file:null`.
Create содержит `syncId`, Update его не содержит.

Модель чтения проверяется на неизменённом `testdata/upstream/organization.json`;
происхождение и revision upstream описаны в [docs/upstream.md](../../../docs/upstream.md).
Случаи дополнительных полей, unknown CompanyType, null и обоих значений
`isEgaisEnable` заданы отдельными синтетическими JSON в `model_test.go`.
