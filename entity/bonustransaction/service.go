package bonustransaction

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
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
		return nil, fmt.Errorf("moysklad: bonustransaction service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-poluchit-bonusnuyu-operaciyu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*BonusTransaction, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get bonustransaction: id must be a UUID")
	}
	return s.do(ctx, "get bonustransaction "+id, http.MethodGet, "entity/bonustransaction/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможно повторное начисление или списание.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-sozdat-bonusnuyu-operaciyu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*BonusTransaction, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create bonustransaction"
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: agent meta.href must not be empty", operation)
	}
	if input.BonusProgram.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: bonusProgram meta.href must not be empty", operation)
	}
	if input.TransactionType == "" {
		return nil, fmt.Errorf("moysklad: %s: transactionType must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/bonustransaction", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-izmenit-bonusnuyu-operaciyu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*BonusTransaction, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update bonustransaction"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/bonustransaction/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-poluchit-bonusnye-operacii
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[BonusTransaction], error) {
	return list.Entity[BonusTransaction](ctx, s.client, "bonustransaction", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*BonusTransaction, error) {
	var result BonusTransaction
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

// Delete — «Удалить Бонусную операцию»: DELETE entity/bonustransaction/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-udalit-bonusnuyu-operaciyu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete bonustransaction: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete bonustransaction "+id, "entity/bonustransaction/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/bonustransaction: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Бонусных операций».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-massovoe-sozdanie-i-obnovlenie-bonusnyh-operacij
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[BonusTransaction], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[BonusTransaction](ctx, s.client, "bonustransaction", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/bonustransaction: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Бонусных операций».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-massovoe-sozdanie-i-obnovlenie-bonusnyh-operacij
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[BonusTransaction], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[BonusTransaction](ctx, s.client, "bonustransaction", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/bonustransaction/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Бонусных операций».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-massovoe-udalenie-bonusnyh-operacij
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "bonustransaction", ids, opts)
}
