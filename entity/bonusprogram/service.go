package bonusprogram

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: bonusprogram service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-poluchit-bonusnuyu-programmu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*BonusProgram, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get bonusprogram: id must be a UUID")
	}
	return s.do(ctx, "get bonusprogram "+id, http.MethodGet, "entity/bonusprogram/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможна вторая программа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-sozdat-bonusnuyu-programmu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*BonusProgram, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create bonusprogram"
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/bonusprogram", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-izmenit-bonusnuyu-programmu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*BonusProgram, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update bonusprogram"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/bonusprogram/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-poluchit-vse-bonusnye-programmy
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[BonusProgram], error) {
	return list.Entity[BonusProgram](ctx, s.client, "bonusprogram", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*BonusProgram, error) {
	var result BonusProgram
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}

// Delete — «Удалить Бонусную программу»: DELETE entity/bonusprogram/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-udalit-bonusnuyu-programmu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete bonusprogram: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete bonusprogram "+id, "entity/bonusprogram/"+id, opts)
}
