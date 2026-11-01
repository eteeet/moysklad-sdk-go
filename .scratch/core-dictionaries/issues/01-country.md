Status: claimed
Blocked by: None

# 01: Country — Get, Create, Update

## What to build

Разработчик решения получает Country по ID, создаёт пользовательскую Country и изменяет
её отдельные поля через SDK. Поддержать весь согласованный набор: 11 полей чтения и 7 записи.
Это законченный срез с моделями, тремя методами, тестами и документацией.

Контракт: [общая spec](../spec.md), разделы Country, общих правил и Testing Decisions.
Образцы — Product/UOM; их архитектуру и общие типы не перерабатывать.

## Acceptance criteria

- [ ] Пакет country предоставляет модель Country, обёртку Service и NewService.
  Конструктор принимает *transport.Client, возвращает *Service и error, отклоняет nil client.
- [ ] Модель читает meta (model.Meta), id/accountId/name/description/code/externalCode (string),
  updated (model.DateTime), owner/group (*model.Reference), shared (bool).
  Ссылки и дата оформлены по образцу Product/UOM; code не преобразуется в число.
- [ ] CreateRequest содержит обязательный Name string; остальные шесть полей записи
  description, code, externalCode, owner, group, shared используют Field с omitzero.
  В UpdateRequest все семь полей используют Field. Readonly-полей в requests нет.
- [ ] Get выполняет GET entity/country/{id} без тела/query; Create — POST entity/country;
  Update — PUT entity/country/{id}. Все возвращают (*Country, error) из ответа сервера.
- [ ] Get/Update отклоняют неверную форму UUID без ограничения версии; Create отклоняет
  только Name=="" без trim; полностью незаданный Update отклоняется до HTTP.
  Явные false, "" и null считаются заданными полями.
- [ ] Успех требует HTTP 200 и непустой ID; некорректный ответ возвращает error и nil.
  Предварительного Get, локального merge, переходов по ссылкам и retry нет.
- [ ] Неизменённый upstream Country проверен по значениям всех согласованных полей обоими
  JSON-кодеками. Добавлен небольшой официальный предустановленный пример с источником/revision;
  отсутствие accountId/owner/group/shared допускается вопреки общей таблице обязательности.
- [ ] Проверены неизвестные/отсутствующие поля ответа, строковый code с ведущими нулями,
  shared=false и даты; Shared не используется для определения происхождения Country.
- [ ] Документированный запрет Update предустановленных Country не проверяется локальной
  классификацией. Синтетическая APIError-тестовая ошибка передаётся вызывающему;
  тест не выдаётся за проверку реального статуса или текста живого API.
- [ ] Requests проверены через encoding/json и encoding/json/v2: минимальный Create,
  полный Create/Update по независимым golden, частичный Update, пропуски, false/true, "" и null.
  Возможность кодировать null не обещает его принятие сервером.
- [ ] Локальные HTTP-тесты с настоящим transport проверяют метод/путь, заголовки,
  отсутствие лишних query/запросов, plain/gzip и данные ответа, отличающиеся от input.
- [ ] Для трёх методов проверены nil/отменённый context, локальная валидация,
  ошибки авторизации/транспорта и APIError 4xx/5xx/429, errors.As/errors.Is,
  контекст без двойного moysklad:, неожиданный статус, malformed/null JSON и пустой/отсутствующий ID.
- [ ] Документация отражает только фактически готовый Country и его TODO.
  Приватные do/validID остаются по образцу; если обнаружено общее исправление,
  проверены все применимые копии согласно правилу coverage, без нового общего исполнителя.
- [ ] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и встроенное ревью Standards/Spec относительно зафиксированного BASE.

## Blocked by

Нет. Ссылки не требуют реализации связанных сущностей. Можно начинать сразу.

## Выполнение и границы

- Один основной агент, одна новая сессия, полный ticket через implement/TDD без учебных остановок.
- До правок прочитать инструкции, всю spec, согласованный scope и текущее состояние Git;
  записать исходный HEAD как BASE ревью. Сохранить уже существующие изменения документов.
- Не включать в commit существовавшие до начала изменения планирования или чужие изменения.
  Свои результаты документации учитывать отдельно от исходного dirty diff.
- Не добавлять List/Delete/expand/batch, другие сущности, фасад, транспортные доработки или live API.
- Reviewer не меняет код; основной агент исправляет замечания в scope. Новые контракты согласовать.
- После полного ticket — локальный commit, Status: resolved и ## Answer с SHA, проверками
  и результатами ревью; сообщить результат пользователю и остановиться. Push запрещён без команды.

## Sources

[Официальная документация Country, зафиксированный revision](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_country.md).
Upstream уже имеется и не редактируется; provenance закреплён механизмом SDK.
