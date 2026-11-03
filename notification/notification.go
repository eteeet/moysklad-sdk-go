// Package notification читает ленту уведомлений текущего пользователя, отмечает уведомления
// прочитанными, удаляет их и читает и изменяет настройки уведомлений. Маршруты — notification/...,
// не entity/. Действия выполняются только над уведомлениями пользователя, от имени которого
// авторизован клиент; права и принадлежность проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notifications
package notification

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Notification — уведомление ленты: общие атрибуты всех типов. Тип уведомления — Meta.Type
// (NotificationTaskChanged, NotificationImportCompleted …); значение открытое, неизвестный тип
// не отвергается. Поля конкретного типа (task, performedBy, diff, good …) SDK не разбирает:
// они в Raw. Все поля только для чтения; Read меняют MarkAsRead и MarkAllAsRead.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notifications#3-obshie-atributy-uvedomlenij
type Notification struct {
	Meta        model.Meta     `json:"meta"`
	ID          string         `json:"id"`
	AccountID   string         `json:"accountId"`
	Created     model.DateTime `json:"created,omitzero"`
	Read        bool           `json:"read"`
	Title       string         `json:"title"`
	Description string         `json:"description"`

	// Raw — всё уведомление как пришло, включая поля конкретного типа из
	// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification-types-details;
	// разбирается вызывающим в свою структуру. При кодировании не пишется.
	Raw jsontext.Value `json:"-"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	type plain Notification
	var known plain
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	*n = Notification(known)
	n.Raw = jsontext.Value(append([]byte(nil), data...))
	return nil
}

// ListOptions — параметры одной страницы ленты; других параметров источник не описывает.
type ListOptions struct {
	Limit  int // 0 — default сервера 1000; не больше query.MaxLimit
	Offset int
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: notification service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу ленты GET notification без автообхода; следующие страницы —
// pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notifications#3-poluchit-lentu-uvedomlenij
func (s *Service) List(ctx context.Context, options ListOptions, opts ...transport.RequestOption) (*model.Collection[Notification], error) {
	transport.ResetResponseInfo(opts)
	const operation = "list notifications"
	params, err := query.ListOptions{Limit: options.Limit, Offset: options.Offset}.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Notification](ctx, s.client, operation, "notification", params, opts...)
}

// Get читает уведомление GET notification/{id}; успех — HTTP 200 и непустой id.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification#3-poluchit-uvedomlenie
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Notification, error) {
	transport.ResetResponseInfo(opts)
	operation := "get notification " + id
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	var result Notification
	resp, err := s.client.Send(ctx, http.MethodGet, "notification/"+id, nil, nil, &result, opts...)
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

// Delete — один DELETE notification/{id}; удалить можно только уведомление текущего
// пользователя. Тело ответа не читается.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification#3-udalit-uvedomlenie
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete notification " + id
	if !validID(id) {
		return fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	return write.Delete(ctx, s.client, operation, "notification/"+id, opts)
}

// MarkAsRead — один PUT notification/{id}/markasread без тела; отметить можно только
// уведомление текущего пользователя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification#3-otmetit-uvedomlenie-kak-prochitannoe
func (s *Service) MarkAsRead(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "mark notification " + id + " as read"
	if !validID(id) {
		return fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	return s.put(ctx, operation, "notification/"+id+"/markasread", nil, opts)
}

// MarkAllAsRead — один PUT notification/markasreadall без тела: все непрочитанные
// уведомления текущего пользователя становятся прочитанными.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/notification/notification#3-otmetit-vse-uvedomleniya-kak-prochitannye
func (s *Service) MarkAllAsRead(ctx context.Context, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	return s.put(ctx, "mark all notifications as read", "notification/markasreadall", nil, opts)
}

// put — один PUT без чтения ответа: источник тела ответа не описывает. Успех — 200 или 204,
// как у write.Delete. PUT не повторяется.
func (s *Service) put(ctx context.Context, operation, endpoint string, body any, opts []transport.RequestOption) error {
	resp, err := s.client.Send(ctx, http.MethodPut, endpoint, nil, body, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
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
