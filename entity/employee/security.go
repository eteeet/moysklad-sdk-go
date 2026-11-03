package employee

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/entity/role"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Security — доступ сотрудника к сервису и его роль. Читать и менять может системный
// администратор; права проверяет сервер. Активация, деактивация и сброс пароля в SDK не входят.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-rabota-s-pravami-sotrudnika
type Security struct {
	IsActive            bool           `json:"isActive"`
	Login               string         `json:"login"`
	Email               string         `json:"email"`
	Group               *SecurityGroup `json:"group,omitzero"`
	AuthorizedHosts     []string       `json:"authorizedHosts"`
	AuthorizedIPNetwork string         `json:"authorizedIpNetwork"`
	AuthorizedIPNetmask string         `json:"authorizedIpNetmask"`
	Role                *role.Role     `json:"role,omitzero"`
}

// SecurityGroup — отдел сотрудника: meta, id и имя.
type SecurityGroup struct {
	Meta model.Meta `json:"meta"`
	ID   string     `json:"id"`
	Name string     `json:"name"`
}

// SecurityRoleRequest задаёт роль по meta (role.Service.Preset или пользовательская роль).
// Для индивидуальной роли: переданные Permissions выставляются, остальные права сервер выключает
// (кроме view=ALL у currency, country, uom); без Permissions сервер оставляет права, которые были
// у сотрудника до смены роли, или значения по умолчанию.
type SecurityRoleRequest struct {
	Meta        model.Meta                     `json:"meta"`
	Permissions model.Field[model.Permissions] `json:"permissions,omitzero"`
}

// UpdateSecurityRequest: role, group и authorizedHosts — из примера изменения прав,
// authorizedIpNetwork и authorizedIpNetmask — из таблицы атрибутов. login и isActive меняются
// активацией и деактивацией (Activate, Deactivate), email — через Service.Update сотрудника.
// null для полей источник не описывает: Field позволяет его отправить, решает сервер.
type UpdateSecurityRequest struct {
	Role                model.Field[SecurityRoleRequest] `json:"role,omitzero"`
	Group               model.Field[model.Reference]     `json:"group,omitzero"`
	AuthorizedHosts     model.Field[[]string]            `json:"authorizedHosts,omitzero"`
	AuthorizedIPNetwork model.Field[string]              `json:"authorizedIpNetwork,omitzero"`
	AuthorizedIPNetmask model.Field[string]              `json:"authorizedIpNetmask,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-poluchit-informaciyu-o-pravah-sotrudnika
func (s *Service) GetSecurity(ctx context.Context, id string, opts ...transport.RequestOption) (*Security, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get employee security: id must be a UUID")
	}
	return s.security(ctx, "get employee "+id+" security", http.MethodGet, id, nil, opts)
}

// UpdateSecurity не активирует сотрудника, не сбрасывает пароль и не назначает роль скрытыми вызовами.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-izmenit-informaciyu-o-pravah-sotrudnika
func (s *Service) UpdateSecurity(ctx context.Context, id string, input UpdateSecurityRequest, opts ...transport.RequestOption) (*Security, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update employee security: id must be a UUID")
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update employee security: no fields specified")
	}
	return s.security(ctx, "update employee "+id+" security", http.MethodPut, id, input, opts)
}

func (s *Service) security(ctx context.Context, operation, method, id string, body any, opts []transport.RequestOption) (*Security, error) {
	var result Security
	resp, err := s.client.Send(ctx, method, "entity/employee/"+id+"/security", nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	// group обязателен в ответе: без него ответ не описывает права сотрудника.
	if result.Group == nil {
		return nil, fmt.Errorf("moysklad: %s: response has no group", operation)
	}
	return &result, nil
}
