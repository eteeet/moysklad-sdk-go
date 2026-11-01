package transport

// Ограничения размера запроса JSON API: не более 20 МБ данных в запросе, 8 КБ в заголовках
// (URL, User-Agent, Authorization и т. д.) и вложенность JSON не глубже 10 (ошибка 2006).
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_restrictions.md
// Источник не уточняет, MB и KB десятичные или двоичные, и как считается вложенность: SDK берёт
// десятичные единицы (строже), а глубину — по одним объектам, корень — уровень 1 (мягче: ложный
// отказ хуже пропуска, который поймает сервер ошибкой 2006). Всё переопределяется в Config.
// Метаданные оцениваются как HTTP/1.1 на проводе; cookies Jar, заголовки своего RoundTripper и
// сжатие HTTP/2 не учитываются. Исходящих потоков нет: тело всегда закодировано целиком.
// Проверки выполняют NewRequest и Send; Do с собственным запросом их не делает.
// Размер ответа SDK не ограничивает.

import (
	"bytes"
	"encoding/json/jsontext"
	"fmt"
	"net/http"
	"strconv"
)

const (
	DefaultMaxBodySize   = 20_000_000
	DefaultMaxHeaderSize = 8_000
	DefaultMaxDepth      = 10
)

// effectiveLimit: 0 — default, отрицательное — без проверки (тогда 0).
func effectiveLimit(value, def int) int {
	switch {
	case value == 0:
		return def
	case value < 0:
		return 0
	}
	return value
}

func (c *Client) checkBody(data []byte) error {
	if c.maxBodySize > 0 && len(data) > c.maxBodySize {
		return fmt.Errorf("moysklad: request body is %d bytes, limit %d", len(data), c.maxBodySize)
	}
	if depth := objectDepth(data); c.maxDepth > 0 && depth > c.maxDepth {
		return fmt.Errorf("moysklad: request body object nesting depth %d exceeds %d", depth, c.maxDepth)
	}
	return nil
}

// checkHeaderSize оценивает request line, Host, Content-Length и каждое значение каждого
// заголовка; значения в ошибку не попадают.
func checkHeaderSize(req *http.Request, bodyLen, limit int) error {
	if limit <= 0 {
		return nil
	}
	size := len(req.Method) + 1 + len(req.URL.RequestURI()) + len(" HTTP/1.1\r\n") + len("Host: \r\n") + len(req.URL.Host)
	// net/http шлёт Content-Length и для записи без тела.
	if bodyLen > 0 || req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		size += len("Content-Length: \r\n") + len(strconv.Itoa(bodyLen))
	}
	for key, values := range req.Header {
		for _, value := range values {
			size += len(key) + len(": \r\n") + len(value)
		}
	}
	if size > limit {
		return fmt.Errorf("moysklad: request URL and headers are %d bytes, limit %d", size, limit)
	}
	return nil
}

// objectDepth — наибольшая вложенность объектов; data — уже проверенный json.Marshal вывод.
func objectDepth(data []byte) int {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	depth, deepest := 0, 0
	for {
		token, err := dec.ReadToken()
		if err != nil {
			return deepest
		}
		switch token.Kind() {
		case '{':
			depth++
			deepest = max(deepest, depth)
		case '}':
			depth--
		}
	}
}
