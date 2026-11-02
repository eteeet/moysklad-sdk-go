package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Download выполняет GET href и копирует тело успешного ответа в w. Возвращает число записанных
// байт — и при ошибке. href — относительный endpoint или абсолютный URL того же origin, что
// BaseURL (meta.downloadHref); чужой origin отклоняется до отправки. Redirect на подписанный URL
// другого origin разрешён, Authorization на нём удаляется; query подписанного URL в текст ошибки
// не попадает. ResponseInfo описывает последний ответ цепочки redirect, в том числе отклонённого
// CheckRedirect; при неразбираемом Location net/http ответ не возвращает, и StatusCode равен 0.
//
// Размер ответа не ограничивается, и общий Config.Timeout к чтению тела не применяется: срок
// загрузки задаёт ctx. Тело закрывается до возврата, слоты Limiter держатся до конца копирования.
// Config.Retry повторяет попытку, только пока в w ничего не записано; после частичной записи
// повтор — решение вызывающего кода. Ответ не 2xx возвращается как APIError, w не трогается.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-ssylki-na-fajly
func (c *Client) Download(ctx context.Context, href string, w io.Writer, opts ...RequestOption) (int64, error) {
	ResetResponseInfo(opts)
	if w == nil {
		return 0, fmt.Errorf("moysklad: download: writer must not be nil")
	}
	options := newRequestOptions(opts...)
	client := *c.httpClient
	client.Timeout = 0
	var policy RetryPolicy
	var retry bool
	for attempt := 1; ; attempt++ {
		req, err := c.newRequest(ctx, http.MethodGet, href, nil, nil, options)
		if err != nil {
			return 0, err
		}
		// Файл — не JSON; Accept-Encoding: gzip оставлен, ответ распаковывается здесь.
		req.Header.Set("Accept", "*/*")
		if attempt == 1 && c.retry != nil && retryableRequest(req, options) {
			retry, policy = true, c.retry.withDefaults()
		}
		n, resp, err, sent := c.downloadOnce(ctx, &client, req, w, options)
		if !sent {
			return 0, err
		}
		if info := options.ResponseInfo; info != nil {
			*info = ResponseInfo{Attempts: attempt}
			if resp != nil {
				info.StatusCode, info.Header = resp.StatusCode, resp.Header.Clone()
			}
		}
		if err == nil || n > 0 || !retry {
			return n, err
		}
		delay, ok := policy.retryDelay(attempt, resp, err)
		if !ok {
			return n, err
		}
		if waited, ctxErr := policy.sleep(ctx, delay); !waited {
			return n, errors.Join(err, ctxErr)
		}
	}
}

// downloadOnce — одна отправка; sent=false, если запрос не ушёл (отмена в ожидании Limiter).
func (c *Client) downloadOnce(ctx context.Context, client *http.Client, req *http.Request, w io.Writer, options RequestOptions) (n int64, resp *http.Response, err error, sent bool) {
	release, err := c.acquire(ctx, req, options)
	if err != nil {
		return 0, nil, err, false
	}
	defer release()
	resp, err = client.Do(req)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			// Подпись временной ссылки — секрет: в ошибке остаются схема, хост и путь.
			if u, parseErr := url.Parse(urlErr.URL); parseErr == nil {
				u.RawQuery, u.Fragment = "", ""
				urlErr.URL = u.String()
			}
			// Неразбираемый Location попадает целиком в текст вложенной ошибки.
			if strings.Contains(urlErr.Err.Error(), "?") {
				urlErr.Err = redactedError{urlErr.Err}
			}
		}
		// При отказе CheckRedirect Do возвращает и полученный ответ (body закрыт) — для ResponseInfo.
		return 0, resp, fmt.Errorf("moysklad: send request: %w", err), true
	}
	defer resp.Body.Close()
	body, err := decodedBody(resp)
	if err != nil {
		return 0, resp, err, true
	}
	if closer, ok := body.(io.Closer); ok && body != io.Reader(resp.Body) {
		defer closer.Close()
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return 0, resp, readAPIError(resp.StatusCode, body), true
	}
	n, err = io.Copy(w, body)
	if err != nil {
		return n, resp, fmt.Errorf("moysklad: download: %w", err), true
	}
	return n, resp, nil, true
}

// redactedError обрезает текст на первом "?": граница query в тексте неизвестна, в ней бывают
// пробелы и кавычки. Unwrap сохраняет errors.Is/As, Timeout/Temporary — классификацию повторов:
// url.Error проверяет их у Err без Unwrap.
type redactedError struct{ err error }

func (e redactedError) Error() string {
	text, _, _ := strings.Cut(e.err.Error(), "?")
	return text + "?…"
}

func (e redactedError) Unwrap() error { return e.err }

func (e redactedError) Timeout() bool {
	t, ok := e.err.(interface{ Timeout() bool })
	return ok && t.Timeout()
}

func (e redactedError) Temporary() bool {
	t, ok := e.err.(interface{ Temporary() bool })
	return ok && t.Temporary()
}
