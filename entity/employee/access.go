package employee

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// ActivateRequest — активация сотрудника. Login обязателен, если сотрудник раньше не был активным
// (это проверяет сервер); Group и Role — как в примере источника, роль — ссылка на системную
// или пользовательскую роль (role.Service).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-aktivaciya-sotrudnika
type ActivateRequest struct {
	Login model.Field[string]              `json:"login,omitzero"`
	Group model.Field[model.Reference]     `json:"group,omitzero"`
	Role  model.Field[SecurityRoleRequest] `json:"role,omitzero"`
}

// Activation — ответ активации: MailActivationRequired = true — на почту сотрудника отправлено
// письмо со ссылкой для входа; false — действует прежний пароль.
type Activation struct {
	MailActivationRequired bool `json:"mailActivationRequired"`
}

// Activate — PUT entity/employee/{id}/access/activate. Побочные эффекты (письмо, место в
// тарифе) — на сервере; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-aktivaciya-sotrudnika
func (s *Service) Activate(ctx context.Context, id string, input ActivateRequest, opts ...transport.RequestOption) (*Activation, error) {
	transport.ResetResponseInfo(opts)
	operation := "activate employee " + id
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Тело — объект даже без полей: ранее активный сотрудник активируется пустым запросом.
	var result *Activation
	resp, err := s.client.Send(ctx, http.MethodPut, "entity/employee/"+id+"/access/activate", nil, input, &result, opts...)
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

// Deactivate — PUT entity/employee/{id}/access/deactivate без тела; успех — 204 (или 200).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-deaktivaciya-sotrudnika
func (s *Service) Deactivate(ctx context.Context, id string, opts ...transport.RequestOption) error {
	return s.access(ctx, "deactivate employee "+id, id, "deactivate", opts)
}

// ResetPassword — PUT entity/employee/{id}/access/resetpassword без тела: новый пароль сервер
// отправляет на почту сотрудника. Успех — 204 (или 200).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-sbros-parolya-sotrudnika
func (s *Service) ResetPassword(ctx context.Context, id string, opts ...transport.RequestOption) error {
	return s.access(ctx, "reset employee "+id+" password", id, "resetpassword", opts)
}

func (s *Service) access(ctx context.Context, operation, id, action string, opts []transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	resp, err := s.client.Send(ctx, http.MethodPut, "entity/employee/"+id+"/access/"+action, nil, nil, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}
