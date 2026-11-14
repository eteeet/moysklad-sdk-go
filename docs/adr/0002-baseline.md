# ADR 0002: действующие решения (baseline)

Статус: принято 2026-11-14 (решение владельца). Заменяет ADR 0001 целиком.

Здесь собраны действующие решения ADR 0001 в той форме, в какой их исполняет код, сверенные с `~/ETEEET/AGENTS.md`,
`GIT_FLOW.md`, `INFRA.md` и `~/ETEEET/moysklad/AGENTS.md`. ADR 0001 остаётся в репозитории как история: почему опции
вызова, а не `Result[T]`, и как переходили 17 services.

## Опции вызова и ResponseInfo

- Последний параметр каждого публичного метода — `opts ...transport.RequestOption`. Результаты существующих методов
  `(*T, error)` не меняются.
- Query-параметры передаются отдельным аргументом конкретного метода, а не через RequestOption: допустимые filter,
  order и fields различаются по endpoint.
- `transport.WithHeader(key, value)` добавляет значение; повтор ключа даёт несколько значений заголовка, а не CSV.
- Приоритет заголовков: значения transport < `transport.Config.Header` < опции вызова; ключ из опций заменяет ключ
  клиента целиком.
- Служебные заголовки `Authorization`, `Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`, `Host`,
  `Idempotency-Key`, `X-Idempotency-Key` задаёт только transport; попытка задать их в Config или опциях — ошибка до
  отправки.
- `transport.WithAuth(auth)` заменяет `Config.Auth` на один вызов; заголовок `Authorization` по-прежнему ставит
  transport.
- `transport.WithResponseInfo(&info)` заполняет `StatusCode`, копию заголовков и число попыток (`Attempts`), в том
  числе при `APIError` и ошибке decode; `StatusCode` и заголовки относятся к последней попытке. Без ответа `info`
  обнуляется. Указатель принадлежит вызову, общего `LastResponse` нет.
- Публичная операция обнуляет `info` на входе через `transport.ResetResponseInfo(opts)`, до локальной валидации:
  ошибка до запроса не оставляет метаданных прошлого вызова.

## Transport

- `transport.Client.Send` и его обёртка `SendNoRedirect` (один вызов без следования 303/202, тело успеха не
  декодируется) — единственный seam JSON-вызова. Limiter и retry живут в transport, а не в services.
- `Download` копирует тело ответа в `io.Writer`; тело ответа закрывает сам на всех путях, writer не закрывает. Слоты
  limiter держатся до конца копирования; retry повторяет попытку, только пока в writer ничего не записано. При
  неразбираемом `Location` ResponseInfo у `Download` остаётся без статуса.
- Policy — настройки, limiter, retry — принадлежит `*transport.Client`; services и фасад, собранные на одном Client,
  делят её.
- `NewRequest`/`Do` — низкоуровневый API без опций; ResponseInfo они не заполняют.

## Отвергнуто

- `Result[T]{Value, Info}`: метаданные были бы доступны всегда, но это меняет сигнатуры всех методов и заставляет
  каждого вызывающего разворачивать `.Value` ради редко нужных заголовков.
