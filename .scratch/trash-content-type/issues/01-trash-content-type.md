Status: ready-for-agent
Type: bug
Blocked by: —

# 01: `Trash` уходит без `Content-Type`

Живая проверка ms_documan 2026-11-05 (`ms_documan/.scratch/replica/live-checks.md`, «Найдено попутно»): `POST entity/{entity}/{id}/trash` без тела уходит без `Content-Type`, и МойСклад отвечает 400 [1037] «Неверно указан Content-Type запроса». С телом `{}` и `Content-Type: application/json` — 200 с пустым ответом. Задать заголовок снаружи нельзя: `transport` его не принимает. Ошибка касается `Trash` всех документов: общий вызов — `internal/write/write.go`, функция `Trash`.

## What to build

`Trash` отправляет запрос, который МойСклад принимает: тело `{}` с `Content-Type: application/json` или заголовок без тела — что подтвердит проверка на живом аккаунте. POST по-прежнему не повторяется.

## Acceptance criteria

- [ ] Тест на `httptest`: запрос `Trash` несёт `Content-Type: application/json`, тело — по выбранному варианту.
- [ ] Живая проверка на тестовом аккаунте: документ уходит в корзину, ответ 200 или 204; результат с датой записан в тикет.
- [ ] Выпуск SDK с исправлением; ms_documan убирает обход в `.scratch/replica/livecheck/main.go`, когда поднимет версию.
