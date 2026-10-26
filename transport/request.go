package transport

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Явные query-параметры заменяют одноимённые параметры из endpoint.
func (c *Client) NewRequest(ctx context.Context, method, endpoint string, query url.Values, body any) (*http.Request, error) {
	resource, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("moysklad: invalid endpoint: %w", err)
	}
	if resource.Path == "" || resource.Opaque != "" || resource.User != nil || resource.Fragment != "" {
		return nil, fmt.Errorf("moysklad: endpoint must contain a path and no credentials or fragment")
	}
	for segment := range strings.SplitSeq(resource.Path, "/") {
		if segment == ".." {
			return nil, fmt.Errorf("moysklad: endpoint must not contain parent path segments")
		}
	}

	requestURL := resource
	if resource.IsAbs() || resource.Host != "" {
		if !sameOrigin(resource, c.baseURL) {
			return nil, fmt.Errorf("moysklad: endpoint must have the same origin as base URL")
		}
	} else {
		requestURL = c.baseURL.JoinPath(resource.EscapedPath())
	}

	rawQuery, err := mergeQuery(resource.RawQuery, query)
	if err != nil {
		return nil, fmt.Errorf("moysklad: invalid query parameters: %w", err)
	}
	requestURL.RawQuery = rawQuery
	requestURL.ForceQuery = false

	var payload io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("moysklad: encode request body: %w", err)
		}
		payload = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("moysklad: create request: %w", err)
	}
	// ошибка 1062 требует application/json;charset=utf-8.
	// https://github.com/moysklad/api-remap-1.2-doc/blob/master/md/_errors.md#L62
	req.Header.Set("Accept", "application/json;charset=utf-8")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "moysklad-sdk-go")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.auth != nil {
		if err := c.auth(req); err != nil {
			return nil, fmt.Errorf("moysklad: authorize request: %w", err)
		}
	}

	return req, nil
}

// Сырой query сохраняется побайтно (nextHref); делим только по &: в filter МойСклад ; — часть значения,
// а url.ParseQuery на ; возвращает ошибку.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-filtraciya-vyborki-s-pomoshyu-parametra-filter
func mergeQuery(raw string, query url.Values) (string, error) {
	var parts []string
	for part := range strings.SplitSeq(raw, "&") {
		if part == "" {
			continue
		}
		key, value, _ := strings.Cut(part, "=")
		name, err := url.QueryUnescape(key)
		if err != nil {
			return "", err
		}
		if _, err := url.QueryUnescape(value); err != nil {
			return "", err
		}
		if _, ok := query[name]; !ok {
			parts = append(parts, part)
		}
	}

	if len(query) == 0 {
		return raw, nil
	}
	if encoded := query.Encode(); encoded != "" {
		parts = append(parts, encoded)
	}
	return strings.Join(parts, "&"), nil
}
