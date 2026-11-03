# Review исправлений и memnet — 2026-11-03

BASE: `5de2fa76148380386a7dec4cce9f81c2ff3b92c0`.
HEAD: `c9399716f063284840341b2410957e7ce9edb460`.
Scope: 783dfb3 (wizard/docs), c939971 (миграция HTTP-тестов на memnet),
150 committed files; working diff AGENTS.md добавляет только заголовок.
Метод: независимые Standards/Spec reviewers, полный make check и lifecycle repro root.

## Standards

Замечаний нет. Предыдущие замечания закрыты:

- Wizard проверяет значения Result и Raw из HTTP-ответа, plain/gzip,
  exact query encoding, request headers и ResponseInfo.
- Docs/godoc объясняют перенос полей read model в отдельный CreateRequest;
  Prepayment/PrepaymentReturn больше не обещают PositionRequest.

Standards reviewer сравнил 117 изменённых test-files после нормализации импорта,
конструктора и TCP-комментариев: assertions не ослаблены. Wizard и новые memnet
files проверены отдельно. Снятие HTTPClient:server.Client() в documenttest согласовано
с механизмом memnet: server.Client() имеет собственный Transport, не использующий hook.
Исключение TCP в DownloadGzip сохраняет проверку собственного HTTPClient/Transport.
Существенных baseline smells и дополнительных подтверждённых нарушений нет.

## Spec

### P1. P2 — закрытый memnet-адрес снова уходит в TCP

`internal/memnet/memnet.go:79` удаляет запись listener из registry. Следующий запрос
к этому адресу попадает в fallback настоящего DialContext на строке 37. Порт
memnet синтетический, фактически не резервируется и может принадлежать другому
процессу. Поэтому закрытый тестовый сервер может неожиданно успешно ответить
через чужой TCP listener, а в Download redirect test туда может попасть и query
с тестовой signature.

Это нарушает изоляцию теста и цель устранить зависимость таких запросов от TCP.
Контракт docs/client.md — сеть для тестов не нужна; sdk-v1/spec.md требует
корректно проверять close/cancellation/redirect/auth/partial output.

Root воспроизвёл детерминированно: пока memnet-сервер ещё открыт, настоящий
TCP listener занял тот же адрес. Затем memnet.Close(), HTTP GET старого URL.
Результат: calls=1 у TCP handler, err=nil, body="other-process".

Исправление: сохранять знание об уже выданных/закрытых memnet-адресах и возвращать
ошибку закрытого listener без TCP fallback. Fallback оставить для адресов,
которые не принадлежали memnet. Добавить regression test с TCP listener на том
же синтетическом адресе.

Других подтверждённых Spec findings нет. net.Pipe сохраняет настоящий
http.Client → http.Transport → http.Server и HTTP serialization: codecs, gzip,
status errors, redirects и cancellation не заменены RoundTripper mock.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., tracked HEAD и working tree diff-check.
  В этом полном прогоне прежние TCP timeouts не воспроизвелись.
- Дополнительный lifecycle repro — **FAIL**, подтверждает P1:
  `go test -race -count=1 -overlay /private/tmp/memnet-review-overlay.json
  -run '^TestReviewClosedAddressDoesNotReachTCP$' ./internal/memnet`.
- Repro source: `/private/tmp/memnet-review-closed_test.go`; Go overlay виртуально
  подключает тест, не меняя source tree SDK. TCP listener только локальный,
  управляемый тестом; внешняя сеть и live API не использовались.
- HEAD оставался c939971. Working AGENTS.md diff и существующие untracked файлы
  reviewer не менял. Созданы только отчёт и временные overlay/repro files.
  .DS_Store не читались/не изменялись; commits/push не выполнялись.

Standards: 0 findings. Spec: 1 finding P2.

## Answer

P1 исправлен в c3bace0: Close больше не удаляет адрес из реестра, dial закрытого
memnet-сервера возвращает ECONNREFUSED без TCP fallback; fallback остался только для адресов,
которые memnet не выдавал. TestClosedServerDoesNotReachTCP (TCP listener на том же синтетическом
адресе) был красным до исправления (calls=1, err=nil) и зелёный после. make check — PASS.
