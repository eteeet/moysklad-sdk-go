package bonustransaction

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
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
