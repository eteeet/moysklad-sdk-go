package variant

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// CharacteristicDefinition — описание характеристики модификаций (размер, цвет …) в метаданных
// variant; значение у конкретной модификации — Characteristic. Type открыт (источник называет
// string). Противоречие: _characteristic.md — required «всегда имеет значение false», _variant.md —
// «является ли характеристика обязательной»; SDK отдаёт значение сервера как есть.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/characteristic
type CharacteristicDefinition struct {
	Meta     model.Meta `json:"meta"`
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Required bool       `json:"required"`
}

// Characteristic — значение характеристики в Variant.characteristics: id и meta указывают
// на CharacteristicDefinition. При записи сервер ищет описание по id, без id — по name,
// без обоих отвечает ошибкой; новые описания при этом не создаются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#5-harakteristiki-modifikacii
type Characteristic struct {
	Meta  model.Meta `json:"meta"`
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Value string     `json:"value"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: variant service requires a client")
	}
	return &Service{client: client}, nil
}

const characteristicsPath = "entity/variant/metadata/characteristics"

// ListCharacteristics читает описания из GET entity/variant/metadata (массив characteristics) —
// форма, которую показывает источник. Это не Collection: страниц нет, остальные поля метаданных
// (attributes, priceTypes) читает metadata.Service. Для entity/variant/metadata/characteristics
// источник пример ответа не приводит, и SDK его не угадывает. Описания не создаются скрыто.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/characteristic#3-poluchit-metadannye
func (s *Service) ListCharacteristics(ctx context.Context, opts ...transport.RequestOption) ([]CharacteristicDefinition, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get variant metadata characteristics"
	var result struct {
		Meta            model.Meta                 `json:"meta"`
		Characteristics []CharacteristicDefinition `json:"characteristics"`
	}
	if err := s.send(ctx, operation, http.MethodGet, "entity/variant/metadata", nil, &result, opts); err != nil {
		return nil, err
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	if result.Characteristics == nil {
		return nil, fmt.Errorf("moysklad: %s: response characteristics must be a non-null array", operation)
	}
	return result.Characteristics, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/characteristic#3-poluchit-harakteristiku
func (s *Service) GetCharacteristic(ctx context.Context, id string, opts ...transport.RequestOption) (*CharacteristicDefinition, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get variant characteristic: id %q must be a UUID", id)
	}
	return s.definition(ctx, "get variant characteristic "+id, http.MethodGet, characteristicsPath+"/"+id, nil, opts)
}

// CreateCharacteristic создаёт одно описание с непустым именем; изменения описаний API не описывает,
// массовое создание (POST массива) в SDK v1 не входит.
// Права на редактирование товаров проверяет сервер. Не повторять после неоднозначного сбоя:
// возможен дубль.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/characteristic#3-sozdat-harakteristiku
func (s *Service) CreateCharacteristic(ctx context.Context, name string, opts ...transport.RequestOption) (*CharacteristicDefinition, error) {
	transport.ResetResponseInfo(opts)
	if name == "" {
		return nil, fmt.Errorf("moysklad: create variant characteristic: name must not be empty")
	}
	body := struct {
		Name string `json:"name"`
	}{name}
	return s.definition(ctx, "create variant characteristic", http.MethodPost, characteristicsPath, body, opts)
}

func (s *Service) definition(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*CharacteristicDefinition, error) {
	var result CharacteristicDefinition
	if err := s.send(ctx, operation, method, endpoint, body, &result, opts); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

func (s *Service) send(ctx context.Context, operation, method, endpoint string, body, dst any, opts []transport.RequestOption) error {
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, dst, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
