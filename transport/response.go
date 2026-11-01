package transport

import (
	"compress/gzip"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
)

// Ограничиваем только тело ошибки; успешные ответы могут содержать большие отчёты.
const maxErrorBodySize = 1 << 20

// Body ответа закрыт до возврата; resp нужен для статуса и заголовков.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-szhatie-soderzhimogo-otvetov
func (c *Client) Do(req *http.Request, dst any) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("moysklad: send request: %w", err)
	}
	defer resp.Body.Close()

	hasBody := req.Method != http.MethodHead &&
		resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusResetContent

	var reader io.Reader = resp.Body
	if hasBody {
		decoded, err := decodedBody(resp)
		if err != nil {
			return resp, err
		}
		reader = decoded
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return resp, readAPIError(resp.StatusCode, reader)
	}

	if dst == nil || !hasBody {
		if _, err := io.Copy(io.Discard, reader); err != nil {
			return resp, fmt.Errorf("moysklad: read response: %w", err)
		}
		return resp, nil
	}

	if err := json.UnmarshalRead(reader, dst); err != nil {
		return resp, fmt.Errorf("moysklad: decode response: %w", err)
	}
	return resp, nil
}

// decodedBody распаковывает gzip-ответ; ошибка уже несёт статус ответа (APIError для не 2xx).
func decodedBody(resp *http.Response) (io.Reader, error) {
	switch encoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding"))); encoding {
	case "", "identity":
		return resp.Body, nil
	case "gzip":
		compressed, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, responseError(resp.StatusCode, fmt.Errorf("open gzip response: %w", err))
		}
		return compressed, nil
	default:
		return nil, responseError(resp.StatusCode, fmt.Errorf("unsupported content encoding %q", encoding))
	}
}

// readAPIError читает тело ошибки до maxErrorBodySize; пустой или невалидный JSON не скрывает статус.
func readAPIError(statusCode int, reader io.Reader) error {
	apiErr := &sdkerrors.APIError{StatusCode: statusCode}
	data, err := io.ReadAll(io.LimitReader(reader, maxErrorBodySize+1))
	if err != nil {
		return errors.Join(apiErr, fmt.Errorf("read error response: %w", err))
	}
	if len(data) > maxErrorBodySize {
		return errors.Join(apiErr, fmt.Errorf("error response exceeds %d bytes", maxErrorBodySize))
	}
	var payload sdkerrors.APIError
	if err := json.Unmarshal(data, &payload); err == nil {
		apiErr.Errors = payload.Errors
	}
	return apiErr
}

// Префикс moysklad: даёт APIError; без него err получает префикс сам.
func responseError(statusCode int, err error) error {
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return errors.Join(&sdkerrors.APIError{StatusCode: statusCode}, err)
	}
	return fmt.Errorf("moysklad: %w", err)
}
