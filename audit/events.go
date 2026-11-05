package audit

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/url"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Event — одно событие контекста. Все поля только для чтения; enum открыты. Diff хранится как
// пришёл: его форма зависит от типа события, а значения — от версии сущности на момент события,
// поэтому SDK не раскладывает его в текущие модели сущностей.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit#3-sobytiya
type Event struct {
	Audit          *model.Reference `json:"audit,omitzero"` // контекст события
	Entity         *model.Reference `json:"entity,omitzero"`
	Name           string           `json:"name"`
	EntityType     string           `json:"entityType"`
	EventType      string           `json:"eventType"`
	ObjectType     string           `json:"objectType"`
	ObjectCount    int              `json:"objectCount"`
	Source         string           `json:"source"`
	UID            string           `json:"uid"`
	Moment         model.DateTime   `json:"moment,omitzero"`
	SupportAccess  bool             `json:"supportAccess"`
	AdditionalInfo string           `json:"additionalInfo"`
	// Diff отсутствует у событий создания.
	Diff jsontext.Value `json:"diff,omitzero"`

	// Raw — всё событие как пришло, включая поля, которых нет в типе.
	Raw jsontext.Value `json:"-"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type plain Event
	var known plain
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	*e = Event(known)
	e.Raw = jsontext.Value(append([]byte(nil), data...))
	return nil
}

// FieldChange — изменение одного атрибута в событиях обновления, удаления, корзины и архива.
// Значения — как пришли: таблица называет их Boolean, но тип зависит от атрибута.
type FieldChange struct {
	OldValue jsontext.Value `json:"oldValue,omitzero"`
	NewValue jsontext.Value `json:"newValue,omitzero"`
}

// RegistrationDiff — diff события регистрации аккаунта, для DecodeDiff.
type RegistrationDiff struct {
	Account string `json:"account"`
	Country string `json:"country"`
}

// PublicationDiff — diff событий публикации документов, для DecodeDiff.
type PublicationDiff struct {
	TemplateName    string `json:"templateName"`
	PublicationHref string `json:"publicationHref"`
}

// EmailDiff — diff событий отправки писем, для DecodeDiff.
type EmailDiff struct {
	SenderEmail  string `json:"senderEmail"`
	TargetEmail  string `json:"targetEmail"`
	SubjectEmail string `json:"subjectEmail"`
	Text         string `json:"text"`
}

// FieldChanges разбирает Diff как «имя атрибута → старое/новое значение».
// Пустой Diff — пустой результат; Diff другой формы — ошибка, Diff при этом не меняется.
// Объект атрибута без oldValue/newValue даёт пустой FieldChange: полные данные — в Diff.
func (e Event) FieldChanges() (map[string]FieldChange, error) {
	changes := map[string]FieldChange{}
	if err := e.DecodeDiff(&changes); err != nil {
		return nil, err
	}
	return changes, nil
}

// DecodeDiff декодирует Diff в dst, например *RegistrationDiff; пустой Diff dst не меняет.
func (e Event) DecodeDiff(dst any) error {
	if len(e.Diff) == 0 {
		return nil
	}
	if err := json.Unmarshal(e.Diff, dst); err != nil {
		return fmt.Errorf("moysklad: decode audit diff: %w", err)
	}
	return nil
}

var entityTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// MaxEventsLimit — максимум limit у событий; без limit сервер отдаёт 25. Источник параметры
// страницы событий не описывает; живой МойСклад (2026-11-05, контекст с 47 событиями) принимает
// limit до 100 и offset, больший limit урезает до 100.
const MaxEventsLimit = 100

// EventsOptions — страница событий: Limit 0 — по умолчанию сервера (25), не больше MaxEventsLimit.
type EventsOptions struct {
	Limit  int
	Offset int
}

func (o EventsOptions) values() (url.Values, error) {
	return query.ListOptions{Limit: o.Limit, Offset: o.Offset}.Values(MaxEventsLimit)
}

// EventsByContext читает одну страницу событий контекста; остальные страницы проходит
// pages.All по nextHref (с тем же limit).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit#3-poluchit-sobytiya-po-kontekstu
func (s *Service) EventsByContext(ctx context.Context, contextID string, options EventsOptions, opts ...transport.RequestOption) (*model.Collection[Event], error) {
	transport.ResetResponseInfo(opts)
	operation := "list audit context events"
	if !validID(contextID) {
		return nil, fmt.Errorf("moysklad: %s: context id %q must be a UUID", operation, contextID)
	}
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Event](ctx, s.client, operation, "audit/"+contextID+"/events", params, opts...)
}

// EventsByEntity читает одну страницу событий сущности entity/{entityType}/{id}/audit.
// entityType — код сущности API (product, demand …); набор открыт. Элементы пользовательских
// справочников адресуются entity/customentity/{metadataId}/{id}: этот маршрут их не покрывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit#3-poluchit-sobytiya-po-suschnosti
func (s *Service) EventsByEntity(ctx context.Context, entityType, id string, options EventsOptions, opts ...transport.RequestOption) (*model.Collection[Event], error) {
	transport.ResetResponseInfo(opts)
	operation := "list entity audit events"
	if !entityTypePattern.MatchString(entityType) || entityType == "customentity" {
		return nil, fmt.Errorf("moysklad: %s: entity type %q must match [a-z][a-z0-9]* and not be customentity", operation, entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: entity id %q must be a UUID", operation, id)
	}
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Event](ctx, s.client, operation, "entity/"+entityType+"/"+id+"/audit", params, opts...)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
