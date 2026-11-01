package transport

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Инвентарь заголовков JSON API:
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md
// Заголовки без типизированной опции задаются через WithHeader и читаются из ResponseInfo.Header.
const (
	HeaderAcceptRegion           = "X-Lognex-Accept-Region"
	HeaderContentRegion          = "X-Lognex-Content-Region"
	HeaderAcceptLanguage         = "Accept-Language"
	HeaderAcceptTimezone         = "X-Lognex-Accept-Timezone"
	HeaderContentTimezone        = "X-Lognex-Content-Timezone"
	HeaderBetaFeature            = "X-Lognex-Remap-Beta-Feature"
	HeaderWebHookDisableByPrefix = "X-Lognex-WebHook-DisableByPrefix"
	HeaderDownloadExpiration     = "X-Lognex-Download-Expiration-Seconds"
	HeaderPrettyPrintJSON        = "Lognex-Pretty-Print-JSON"
	HeaderFormatMillisecond      = "X-Lognex-Format-Millisecond"
	HeaderRateLimitLimit         = "X-RateLimit-Limit"
	HeaderRateLimitRemaining     = "X-RateLimit-Remaining"
	HeaderRetryTimeInterval      = "X-Lognex-Retry-TimeInterval"
	HeaderReset                  = "X-Lognex-Reset"
	HeaderLognexRetryAfter       = "X-Lognex-Retry-After"
	HeaderRetryAfter             = "Retry-After"
	HeaderAuth                   = "X-Lognex-Auth"
	HeaderAuthMessage            = "X-Lognex-Auth-Message"
	HeaderAPIVersionDeprecated   = "X-Lognex-API-Version-Deprecated"
	HeaderLocation               = "Location"
	HeaderContentLocation        = "Content-Location"
)

// Формат даты X-Lognex-Accept-Timezone и X-Lognex-Content-Timezone: «Wed, 16 Aug 2017 23:07:01 +0700».
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_money.md#L41-L46
const timezoneLayout = time.RFC1123Z

// Регионы X-Lognex-Accept-Region; значение не ограничено этим списком.
const (
	RegionRU = "ru"
	RegionUZ = "uz"
	RegionKZ = "kz"
)

// WithAcceptRegion нужен для региональных полей аккаунтов не из России.
func WithAcceptRegion(region string) RequestOption {
	return withSingleHeader(HeaderAcceptRegion, region)
}

// WithAcceptLanguage меняет только язык сообщений об ошибках; сейчас поддержаны ru и en.
func WithAcceptLanguage(language string) RequestOption {
	return withSingleHeader(HeaderAcceptLanguage, language)
}

// WithAcceptTimezone передаёт текущую дату клиента с числовой зоной для отчётов;
// сервер смещает операции в эту зону. DateTime моделей от заголовка не зависит.
func WithAcceptTimezone(now time.Time) RequestOption {
	return withSingleHeader(HeaderAcceptTimezone, now.Format(timezoneLayout))
}

// API принимает один beta-тег: повторный вызов заменяет значение, CSV не поддерживается.
// Пустой и неизвестный тег сервер молча игнорирует.
func WithBetaFeature(tag string) RequestOption {
	return withSingleHeader(HeaderBetaFeature, tag)
}

// WithWebHookDisableByPrefix отключает webhooks с этими URL-префиксами (scheme и полный домен)
// в рамках запроса; каждый префикс — отдельное значение заголовка, не CSV. Формат префикса проверяет API.
// Устаревший X-Lognex-WebHook-Disable не поддерживается: API отключил его 2026-10-01.
func WithWebHookDisableByPrefix(prefixes ...string) RequestOption {
	return func(o *RequestOptions) {
		for _, prefix := range prefixes {
			WithHeader(HeaderWebHookDisableByPrefix, prefix)(o)
		}
	}
}

// WithDownloadExpiration задаёт срок жизни временной ссылки на файл; API принимает целые
// секунды (по умолчанию 60, максимум 86400), дробная часть отбрасывается. Проверку значения выполняет API.
func WithDownloadExpiration(ttl time.Duration) RequestOption {
	return withSingleHeader(HeaderDownloadExpiration, strconv.FormatInt(int64(ttl/time.Second), 10))
}

func withSingleHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Header == nil {
			o.Header = make(http.Header)
		}
		o.Header.Set(key, value)
	}
}

// Разобранные значения ниже возвращают ok=false для отсутствующего, нечислового,
// отрицательного или переполняющего значения; сырое значение остаётся в Header.

// X-RateLimit-* не обязательно выражены во взвешенных единицах лимита запросов.
func (i ResponseInfo) RateLimitLimit() (int64, bool) {
	return headerInt(i.Header, HeaderRateLimitLimit)
}

func (i ResponseInfo) RateLimitRemaining() (int64, bool) {
	return headerInt(i.Header, HeaderRateLimitRemaining)
}

// Значения X-Lognex-* ниже — длительности в миллисекундах, а не Unix time.
func (i ResponseInfo) RetryTimeInterval() (time.Duration, bool) {
	return headerMillis(i.Header, HeaderRetryTimeInterval)
}

// ResetAfter — время до сброса ограничения; 0, если ограничение не установлено.
func (i ResponseInfo) ResetAfter() (time.Duration, bool) {
	return headerMillis(i.Header, HeaderReset)
}

func (i ResponseInfo) LognexRetryAfter() (time.Duration, bool) {
	return headerMillis(i.Header, HeaderLognexRetryAfter)
}

// RetryAfter разбирает стандартный Retry-After: секунды или HTTP-date относительно now.
// Дата в прошлом даёт 0.
func (i ResponseInfo) RetryAfter(now time.Time) (time.Duration, bool) {
	value := strings.TrimSpace(i.Header.Get(HeaderRetryAfter))
	if value == "" {
		return 0, false
	}
	if seconds, ok := parseNonNegative(value); ok {
		if seconds > math.MaxInt64/int64(time.Second) {
			return 0, false
		}
		return time.Duration(seconds) * time.Second, true
	}
	date, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	return max(date.Sub(now), 0), true
}

// ContentTimezone — дата клиента в понимании сервера с числовой зоной из ответа отчёта.
func (i ResponseInfo) ContentTimezone() (time.Time, bool) {
	date, err := time.Parse(timezoneLayout, strings.TrimSpace(i.Header.Get(HeaderContentTimezone)))
	return date, err == nil
}

func headerInt(header http.Header, key string) (int64, bool) {
	return parseNonNegative(strings.TrimSpace(header.Get(key)))
}

func headerMillis(header http.Header, key string) (time.Duration, bool) {
	ms, ok := headerInt(header, key)
	if !ok || ms > math.MaxInt64/int64(time.Millisecond) {
		return 0, false
	}
	return time.Duration(ms) * time.Millisecond, true
}

func parseNonNegative(value string) (int64, bool) {
	n, err := strconv.ParseInt(value, 10, 64)
	return n, err == nil && n >= 0
}
