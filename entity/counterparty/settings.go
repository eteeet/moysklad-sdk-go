package counterparty

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/single"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Settings — настройки справочника контрагентов (meta.type counterpartysettings).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-nastrojki-spravochnika-kontragentov
type Settings struct {
	Meta            model.Meta            `json:"meta"`
	UniqueCodeRules model.UniqueCodeRules `json:"uniqueCodeRules"`
	CreateShared    bool                  `json:"createShared"`
}

// SettingsRequest — частичное изменение: незаданные поля остаются прежними.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-nastrojki-spravochnika-kontragentov
type SettingsRequest struct {
	UniqueCodeRules model.Field[model.UniqueCodeRulesRequest] `json:"uniqueCodeRules,omitzero"`
	CreateShared    model.Field[bool]                         `json:"createShared,omitzero"`
}

// GetSettings — GET entity/counterparty/settings.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-nastrojki-spravochnika-kontragentov
func (s *Service) GetSettings(ctx context.Context, opts ...transport.RequestOption) (*Settings, error) {
	return single.Get[Settings](ctx, s.client, "get counterparty settings", "entity/counterparty/settings", opts)
}

// UpdateSettings — PUT entity/counterparty/settings; ответ — настройки целиком.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-nastrojki-spravochnika-kontragentov
func (s *Service) UpdateSettings(ctx context.Context, input SettingsRequest, opts ...transport.RequestOption) (*Settings, error) {
	return single.Put[Settings](ctx, s.client, "update counterparty settings", "entity/counterparty/settings", input, opts)
}
