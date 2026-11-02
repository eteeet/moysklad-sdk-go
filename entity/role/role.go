// Package role читает и изменяет пользовательские роли и читает предустановленные роли.
// Получать и изменять роли может системный администратор или решение; права проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role
package role

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Role — пользовательская роль (meta.type customrole) или предустановленная роль, у которой
// сервер возвращает только meta (systemrole, individualrole …).
type Role struct {
	Meta        model.Meta        `json:"meta"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Permissions model.Permissions `json:"permissions"`
}

// CreateRequest: права, не переданные в Permissions, сервер выставляет по своим правилам.
type CreateRequest struct {
	Name        string                         `json:"name"`
	Permissions model.Field[model.Permissions] `json:"permissions,omitzero"`
}

type UpdateRequest struct {
	Name        model.Field[string]            `json:"name,omitzero"`
	Permissions model.Field[model.Permissions] `json:"permissions,omitzero"`
}

// Preset — предустановленная роль; читается по имени, а не по UUID. Набор открыт: имя
// проверяется только как безопасный сегмент пути.
type Preset string

const (
	PresetOwner      Preset = "owner"
	PresetAdmin      Preset = "admin"
	PresetIndividual Preset = "individual"
	PresetCashier    Preset = "cashier"
	// Роль сотрудника производства доступна только с опцией «Управление производством».
	PresetWorker Preset = "worker"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: role service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role#3-poluchit-polzovatelskuyu-rol
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Role, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: role id must be a UUID")
	}
	return s.do(ctx, "get role "+id, http.MethodGet, "entity/role/"+id, nil, hasID, opts)
}

// Preset читает предустановленную роль: сервер возвращает только её meta, без прав.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-zapros-na-poluchenie-roli-admina
func (s *Service) Preset(ctx context.Context, preset Preset, opts ...transport.RequestOption) (*Role, error) {
	transport.ResetResponseInfo(opts)
	if !presetPattern.MatchString(string(preset)) {
		return nil, fmt.Errorf("moysklad: role preset %q must match [a-z][a-z0-9]*", preset)
	}
	return s.do(ctx, "get role "+string(preset), http.MethodGet, "entity/role/"+string(preset), nil, hasMeta, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role#3-sozdat-polzovatelskuyu-rol
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль роли.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Role, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create role: name must not be empty")
	}
	return s.do(ctx, "create role", http.MethodPost, "entity/role", input, hasID, opts)
}

// Update заменяет переданные поля; состав прав после записи определяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role#3-obnovit-polzovatelskuyu-rol
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Role, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update role: id must be a UUID")
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update role: no fields specified")
	}
	return s.do(ctx, "update role "+id, http.MethodPut, "entity/role/"+id, input, hasID, opts)
}

// List читает одну страницу пользовательских ролей без автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role#3-poluchit-vse-polzovatelskie-roli
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Role], error) {
	return list.Entity[Role](ctx, s.client, "role", options, nil, opts...)
}

var presetPattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Критерии успеха: у пользовательской роли — id, у предустановленной (только meta) — meta.href.
func hasID(r *Role) bool   { return r.ID != "" }
func hasMeta(r *Role) bool { return r.Meta.Href != "" }

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, valid func(*Role) bool, opts []transport.RequestOption) (*Role, error) {
	var result Role
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if !valid(&result) {
		return nil, fmt.Errorf("moysklad: %s: response has no id or meta.href", operation)
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

// Delete — «Удалить пользовательскую роль»: DELETE entity/role/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role#3-udalit-polzovatelskuyu-rol
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete role: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete role "+id, "entity/role/"+id, opts)
}
