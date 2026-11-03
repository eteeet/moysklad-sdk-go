package notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Settings — настройки уведомлений текущего пользователя: ключ — группа (customerOrder,
// dataExchange, invoice, retail, scripts, stock, task, mentions, onlineStores, followedEvents
// у источника). Ключи открытые: неизвестная группа читается, как пришла.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification-settings#3-struktura-dannyh
type Settings map[string]GroupSettings

// GroupSettings — подписка на группу и каналы доставки. Ключи ChannelsEnabled открытые
// (email, push, interface у источника); отсутствующий ключ и false различаются через ok.
type GroupSettings struct {
	Enabled         bool            `json:"enabled"`
	ChannelsEnabled map[string]bool `json:"channelsEnabled"`
}

// SettingsRequest — изменение частичное: не переданные группы, флаг Enabled и каналы
// сервер не меняет. Ключи открытые, как у Settings.
type SettingsRequest map[string]GroupSettingsRequest

// GroupSettingsRequest: Enabled — model.Field, model.Value(false) отправляется явно.
// ChannelsEnabled передаёт только перечисленные каналы, false в карте — явное отключение.
// Источник: отключение групп scripts и onlineStores недопустимо, «параметр enabled
// игнорируется» — для этих групп или вообще, источник не уточняет; SDK отправляет как задано.
type GroupSettingsRequest struct {
	Enabled         model.Field[bool] `json:"enabled,omitzero"`
	ChannelsEnabled map[string]bool   `json:"channelsEnabled,omitzero"`
}

// Settings — один GET notification/settings. Ответ без meta: ошибка, если тело не JSON-объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification-settings#3-poluchit-nastrojki-uvedomlenij
func (s *Service) Settings(ctx context.Context, opts ...transport.RequestOption) (Settings, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get notification settings"
	var result Settings
	resp, err := s.client.Send(ctx, http.MethodGet, "notification/settings", nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result == nil {
		return nil, fmt.Errorf("moysklad: %s: response is not an object", operation)
	}
	return result, nil
}

// UpdateSettings — один PUT notification/settings. Источник тела ответа не описывает:
// SDK его не читает, актуальные настройки — Settings. Пустой запрос и группа без полей
// отклоняются до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification-settings#3-izmenit-nastrojki-uvedomlenij
func (s *Service) UpdateSettings(ctx context.Context, input SettingsRequest, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	const operation = "update notification settings"
	if len(input) == 0 {
		return fmt.Errorf("moysklad: %s: no groups specified", operation)
	}
	for group, settings := range input {
		if settings.Enabled.IsZero() && len(settings.ChannelsEnabled) == 0 {
			return fmt.Errorf("moysklad: %s: group %q has no fields", operation, group)
		}
	}
	return s.put(ctx, operation, "notification/settings", input, opts)
}
