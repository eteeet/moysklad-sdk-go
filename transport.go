package moysklad

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAttempts = 5
	defaultPause    = time.Second
)

type transport struct {
	baseURL  string
	token    string
	hc       *http.Client
	attempts int
	pause    time.Duration
	limiter  Limiter
	log      *slog.Logger
}

// do выполняет запрос к API с повторами и разбирает JSON-ответ в out.
func (t *transport) do(ctx context.Context, method, path string, query url.Values, in, out any) error {
	u := t.baseURL + "/" + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body []byte

	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}

		body = b
	}

	for attempt := 1; ; attempt++ {
		err := t.send(ctx, method, u, body, out)
		if err == nil {
			return nil
		}

		wait, retry := t.retryDecision(method, attempt, err)
		if !retry || attempt >= t.attempts {
			return fmt.Errorf("transport.send: %w", err)
		}

		t.log.InfoContext(ctx, "moysklad: повтор запроса",
			slog.String("method", method), slog.String("path", path),
			slog.Int("attempt", attempt), slog.Duration("wait", wait), slog.Any("error", err))

		err = sleep(ctx, wait)
		if err != nil {
			return fmt.Errorf("sleep: %w", err)
		}
	}
}

// retryDecision решает по ошибке попытки, повторять ли запрос, и возвращает паузу перед повтором.
// 429 повторяем всегда и сообщаем ограничителю. 5xx — только для GET, PUT и DELETE:
// POST мог выполниться, и повтор создал бы дубль.
func (t *transport) retryDecision(method string, attempt int, err error) (wait time.Duration, retry bool) {
	var e *Error
	if !errors.As(err, &e) {
		return 0, false
	}

	switch {
	case e.StatusCode == http.StatusTooManyRequests:
		wait = e.retryAfter
		if wait <= 0 {
			wait = t.pause
		}

		if t.limiter != nil {
			t.limiter.Backoff(wait)

			return 0, true
		}

		return wait, true
	case e.StatusCode >= http.StatusInternalServerError && idempotent(method):
		return t.pause * time.Duration(attempt), true
	default:
		return 0, false
	}
}

func idempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

// send выполняет одну попытку запроса. body — JSON тела запроса или nil.
func (t *transport) send(ctx context.Context, method, u string, body []byte, out any) error {
	if t.limiter != nil {
		release, err := t.limiter.Acquire(ctx)
		if err != nil {
			return fmt.Errorf("Limiter.Acquire: %w", err)
		}
		defer release()
	}

	var reqBody io.Reader = http.NoBody
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bearer "+t.token)
	// МойСклад без gzip отвечает 415. Заголовок ставим сами: HTTP-клиент потребителя может не сжимать.
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := t.hc.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := uncompressed(resp)
	if err != nil {
		return fmt.Errorf("uncompressed: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return apiError(resp, respBody)
	}

	err = decode(respBody, out)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}

// retryAfter читает паузу из X-Lognex-Retry-After (мс); 0 — заголовка нет.
func retryAfter(h http.Header) time.Duration {
	ms, err := strconv.Atoi(h.Get("X-Lognex-Retry-After"))
	if err != nil || ms <= 0 {
		return 0
	}

	return time.Duration(ms) * time.Millisecond
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("ctx.Err: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

// uncompressed возвращает тело ответа, распакованное по Content-Encoding.
func uncompressed(resp *http.Response) (io.Reader, error) {
	if !strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		return resp.Body, nil
	}

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader: %w", err)
	}

	return zr, nil
}

// maxErrorBody ограничивает чтение тела ошибки: HTML-страница прокси может быть большой.
const maxErrorBody = 64 << 10

// apiError разбирает тело ошибки МойСклад {"errors":[...]}. Если тело не JSON, текст ошибки — само тело.
func apiError(resp *http.Response, body io.Reader) error {
	e := &Error{StatusCode: resp.StatusCode, retryAfter: retryAfter(resp.Header)}

	raw, err := io.ReadAll(io.LimitReader(body, maxErrorBody))
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	var parsed Errors
	if json.Unmarshal(raw, &parsed) != nil || len(parsed.Errors) == 0 {
		e.Message = strings.TrimSpace(string(raw))

		return e
	}

	e.Errors = parsed.Errors

	e.Message = parsed.Errors[0].Error

	if c := parsed.Errors[0].Code; c != nil {
		e.Code = *c
	}

	return e
}

func decode(r io.Reader, out any) error {
	if out == nil {
		return nil
	}

	err := json.NewDecoder(r).Decode(out)
	if err != nil {
		return fmt.Errorf("json.Decoder.Decode: %w", err)
	}

	return nil
}
