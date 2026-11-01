# Промпты реализации: пять отдельных сессий

Планирование завершено. Начинай новую сессию в том же рабочем каталоге SDK
и копируй один соответствующий промпт целиком. После завершения ticket проверь результат;
следующий запускай в новой сессии. Рекомендуемый порядок: 01 → 02 → 03 → 04 → 05.

Единственные блокирующие зависимости: 04 и 05 требуют готового 03.
Organization не зависит от Counterparty. Реализацию нескольких tickets одновременно
в одном рабочем дереве здесь не запускаем.

Spec, tickets и текущие изменения планирования пока локальные, без отдельного commit.
Новая сессия должна сохранить исходный dirty diff, а не включать его целиком в commit реализации.
Точные scope, критерии приёмки и полномочия находятся в spec/ticket и temp_plan.

## Сессия 1 — Country

```text
Работай в /Users/av/ETEEET/moysklad/moysklad-sdk-go.
Прочитай AGENTS.md, CLAUDE.md, все документы docs/agents/ и docs/temp_plan.md.
Прочитай .scratch/core-dictionaries/spec.md и ticket .scratch/core-dictionaries/issues/01-country.md.
Выполни только этот согласованный ticket через implement: TDD на согласованных точках, проверки, встроенное code-review Standards/Spec, локальный commit.
Один основной агент; субагенты только для ревью. Ticket выполняй целиком без учебных остановок.
До правок запиши исходный HEAD как BASE и проверь рабочее дерево. Сохрани начальный dirty diff; в commit включи только свои изменения реализации этого ticket, без предшествующих изменений планирования и чужих правок.
Архитектуру Product/UOM сохраняй. Новые контракты сначала согласуй. .env, живой API, push, PR и VPS требуют отдельной команды.
По завершении обнови ticket фактическим результатом, сообщи SHA, проверки и результаты ревью; остановись. Следующий ticket самостоятельно не начинай.
```

## Сессия 2 — ExpenseItem

```text
Работай в /Users/av/ETEEET/moysklad/moysklad-sdk-go.
Прочитай AGENTS.md, CLAUDE.md, все документы docs/agents/ и docs/temp_plan.md.
Прочитай .scratch/core-dictionaries/spec.md и ticket .scratch/core-dictionaries/issues/02-expenseitem.md.
Выполни только этот согласованный ticket через implement: TDD на согласованных точках, проверки, встроенное code-review Standards/Spec, локальный commit.
Один основной агент; субагенты только для ревью. Ticket выполняй целиком без учебных остановок.
До правок запиши исходный HEAD как BASE и проверь рабочее дерево. Сохрани начальный dirty diff; в commit включи только свои изменения реализации этого ticket, без предшествующих изменений планирования и чужих правок.
Архитектуру Product/UOM сохраняй. Новые контракты сначала согласуй. .env, живой API, push, PR и VPS требуют отдельной команды.
По завершении обнови ticket фактическим результатом, сообщи SHA, проверки и результаты ревью; остановись. Следующий ticket самостоятельно не начинай.
```

## Сессия 3 — Address, AddressRequest, CompanyType

```text
Работай в /Users/av/ETEEET/moysklad/moysklad-sdk-go.
Прочитай AGENTS.md, CLAUDE.md, все документы docs/agents/ и docs/temp_plan.md.
Прочитай .scratch/core-dictionaries/spec.md и ticket .scratch/core-dictionaries/issues/03-shared-requisites.md.
Выполни только этот согласованный ticket через implement: TDD публичного JSON-контракта общих типов, проверки, отдельное встроенное code-review Standards/Spec, локальный commit.
Один основной агент; субагенты только для ревью. Ticket выполняй целиком без учебных остановок.
До правок запиши исходный HEAD как BASE и проверь рабочее дерево. Сохрани начальный dirty diff; в commit включи только свои изменения реализации этого ticket, без предшествующих изменений планирования и чужих правок.
Реализуй самостоятельные общие типы; HTTP-методы Counterparty/Organization остаются для отдельных tickets. Новые контракты сначала согласуй. .env, живой API, push, PR и VPS требуют отдельной команды.
По завершении обнови ticket фактическим результатом, сообщи SHA, проверки и результаты ревью; остановись. Следующий ticket самостоятельно не начинай.
```

## Сессия 4 — Counterparty

```text
Работай в /Users/av/ETEEET/moysklad/moysklad-sdk-go.
Прочитай AGENTS.md, CLAUDE.md, все документы docs/agents/ и docs/temp_plan.md.
Прочитай .scratch/core-dictionaries/spec.md и ticket .scratch/core-dictionaries/issues/04-counterparty.md.
Проверь зависимость .scratch/core-dictionaries/issues/03-shared-requisites.md: Status: resolved и общие типы действительно реализованы. Если зависимость не выполнена, сообщи об этом и остановись, не реализуя её попутно.
Выполни только ticket 04 через implement: TDD на согласованных точках, проверки, встроенное code-review Standards/Spec, локальный commit.
Один основной агент; субагенты только для ревью. Ticket выполняй целиком без учебных остановок.
До правок запиши исходный HEAD как BASE и проверь рабочее дерево. Сохрани начальный dirty diff; в commit включи только свои изменения реализации этого ticket, без предшествующих изменений планирования и чужих правок.
Архитектуру Product/UOM сохраняй. Новые контракты сначала согласуй. .env, живой API, push, PR и VPS требуют отдельной команды.
По завершении обнови ticket фактическим результатом, сообщи SHA, проверки и результаты ревью; остановись. Следующий ticket самостоятельно не начинай.
```

## Сессия 5 — Organization

```text
Работай в /Users/av/ETEEET/moysklad/moysklad-sdk-go.
Прочитай AGENTS.md, CLAUDE.md, все документы docs/agents/ и docs/temp_plan.md.
Прочитай .scratch/core-dictionaries/spec.md и ticket .scratch/core-dictionaries/issues/05-organization.md.
Проверь зависимость .scratch/core-dictionaries/issues/03-shared-requisites.md: Status: resolved и общие типы действительно реализованы. Если зависимость не выполнена, сообщи об этом и остановись, не реализуя её попутно. Ticket 04 не является зависимостью.
Выполни только ticket 05 через implement: TDD на согласованных точках, проверки, встроенное code-review Standards/Spec, локальный commit.
Один основной агент; субагенты только для ревью. Ticket выполняй целиком без учебных остановок.
До правок запиши исходный HEAD как BASE и проверь рабочее дерево. Сохрани начальный dirty diff; в commit включи только свои изменения реализации этого ticket, без предшествующих изменений планирования и чужих правок.
Архитектуру Product/UOM сохраняй. Новые контракты сначала согласуй. .env, живой API, push, PR и VPS требуют отдельной команды.
По завершении обнови ticket фактическим результатом, сообщи SHA, проверки и результаты ревью; остановись. Сущность Service самостоятельно не начинай.
```
