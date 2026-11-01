// Package metadata читает метаданные сущностей: общий каталог, метаданные типа, пользовательские
// справочники, описания дополнительных полей и статусы. Метаданные есть не у каждой сущности;
// создание и изменение описаний доп. полей в SDK не входит. Характеристики модификаций
// (entity/variant/metadata/characteristics) — отдельный ресурс каталога товаров.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-metadannye-suschnosti
package metadata

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Entity — метаданные одного типа сущности. Набор полей зависит от типа: известные поля
// разобраны, остальные (tags контрагента и др.) — в Raw.
type Entity struct {
	Meta         model.Meta     `json:"meta"`
	CreateShared bool           `json:"createShared"`
	Attributes   *AttributeList `json:"attributes,omitzero"`
	States       []State        `json:"states"`

	// Raw — весь ответ как пришёл.
	Raw jsontext.Value `json:"-"`
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type plain Entity
	var known plain
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	*e = Entity(known)
	e.Raw = jsontext.Value(append([]byte(nil), data...))
	return nil
}

// AttributeList — описания доп. полей в метаданных. Большинство примеров источника отдаёт
// массив, часть (product, upstream-фикстуры) — коллекцию {meta[, rows]}, где без expand есть
// только meta. Meta заполнен только у формы коллекции.
type AttributeList struct {
	Meta *model.CollectionMeta
	Rows []Attribute
}

func (l *AttributeList) UnmarshalJSON(data []byte) error {
	*l = AttributeList{}
	if jsontext.Value(data).Kind() == '[' {
		return json.Unmarshal(data, &l.Rows)
	}
	var collection struct {
		Meta model.CollectionMeta `json:"meta"`
		Rows []Attribute          `json:"rows"`
	}
	if err := json.Unmarshal(data, &collection); err != nil {
		return err
	}
	l.Meta, l.Rows = &collection.Meta, collection.Rows
	return nil
}

// Attribute — описание дополнительного поля. Type открыт (string, long, customentity …);
// CustomEntityMeta есть у полей типа «справочник».
type Attribute struct {
	Meta             model.Meta  `json:"meta"`
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Type             string      `json:"type"`
	Required         bool        `json:"required"`
	Show             bool        `json:"show"`
	Description      string      `json:"description"`
	CustomEntityMeta *model.Meta `json:"customEntityMeta,omitzero"`
}

// State — статус документа или сущности. Color — ARGB целым числом; StateType открыт.
type State struct {
	Meta       model.Meta `json:"meta"`
	ID         string     `json:"id"`
	AccountID  string     `json:"accountId"`
	Name       string     `json:"name"`
	Color      int64      `json:"color"`
	StateType  string     `json:"stateType"`
	EntityType string     `json:"entityType"`
}

// CompanySettings — метаданные настроек компании; CustomEntities — пользовательские справочники.
type CompanySettings struct {
	Meta           model.Meta     `json:"meta"`
	CustomEntities []CustomEntity `json:"customEntities"`

	// Raw — весь ответ как пришёл.
	Raw jsontext.Value `json:"-"`
}

func (c *CompanySettings) UnmarshalJSON(data []byte) error {
	type plain CompanySettings
	var known plain
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	*c = CompanySettings(known)
	c.Raw = jsontext.Value(append([]byte(nil), data...))
	return nil
}

// CustomEntity — пользовательский справочник; ID — id его метаданных для GetCustomEntity.
type CustomEntity struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	EntityMeta   model.Meta     `json:"entityMeta"`
	CreateShared bool           `json:"createShared"`
	Attributes   *AttributeList `json:"attributes,omitzero"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: metadata service requires a client")
	}
	return &Service{client: client}, nil
}

var entityTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Catalog читает GET entity/metadata: ссылку на метаданные каждого типа. types ограничивает
// выборку фильтром type=…; без types сервер отдаёт все типы. Ответ — объект «тип → {meta}»,
// не коллекция; поля кроме meta источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-metadannye-suschnosti
func (s *Service) Catalog(ctx context.Context, types []string, opts ...transport.RequestOption) (map[string]model.Reference, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get metadata catalog"
	var params url.Values
	if len(types) > 0 {
		conditions := make([]query.Condition, len(types))
		for i, t := range types {
			if !entityTypePattern.MatchString(t) {
				return nil, fmt.Errorf("moysklad: %s: entity type %q must match [a-z][a-z0-9]*", operation, t)
			}
			conditions[i] = query.Eq("type", t)
		}
		filter, err := query.Filter(conditions...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		params = url.Values{query.ParamFilter: {filter}}
	}
	var result map[string]model.Reference
	if err := s.do(ctx, operation, http.MethodGet, "entity/metadata", params, nil, &result, opts); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("moysklad: %s: response is not an object", operation)
	}
	return result, nil
}

// Get читает GET entity/{entityType}/metadata. Существование ресурса у типа проверяет сервер.
func (s *Service) Get(ctx context.Context, entityType string, opts ...transport.RequestOption) (*Entity, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get metadata: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	return s.entity(ctx, "get "+entityType+" metadata", "entity/"+entityType+"/metadata", opts)
}

// GetCustomEntity читает метаданные пользовательского справочника по id его метаданных
// (CompanySettings.CustomEntities[i].ID).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity
func (s *Service) GetCustomEntity(ctx context.Context, metadataID string, opts ...transport.RequestOption) (*Entity, error) {
	transport.ResetResponseInfo(opts)
	if !validID(metadataID) {
		return nil, fmt.Errorf("moysklad: get customentity metadata: id %q must be a UUID", metadataID)
	}
	return s.entity(ctx, "get customentity metadata "+metadataID, "entity/customentity/"+metadataID+"/metadata", opts)
}

// GetCompanySettings читает GET context/companysettings/metadata со списком пользовательских справочников.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysettings
func (s *Service) GetCompanySettings(ctx context.Context, opts ...transport.RequestOption) (*CompanySettings, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get companysettings metadata"
	var result CompanySettings
	if err := s.do(ctx, operation, http.MethodGet, "context/companysettings/metadata", nil, nil, &result, opts); err != nil {
		return nil, err
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	return &result, nil
}

// AttributesPage — страница описаний доп. полей; нули не передаются (default сервера 1000 и 0).
type AttributesPage struct {
	Limit  int
	Offset int
}

// ListAttributes читает одну страницу описаний доп. полей типа (limit до 1000).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchit-vse-dopolnitelnye-polya-dlya-ukazannogo-tipa
func (s *Service) ListAttributes(ctx context.Context, entityType string, page AttributesPage, opts ...transport.RequestOption) (*model.Collection[Attribute], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get attributes: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	operation := "get " + entityType + " attributes"
	params, err := query.ListOptions{Limit: page.Limit, Offset: page.Offset}.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Attribute](ctx, s.client, operation, "entity/"+entityType+"/metadata/attributes", params, opts...)
}

// GetAttribute читает одно описание доп. поля. Таблица помечает id обязательным, но пример
// customerorder отдаёт поле без id: критерий успеха — meta.href.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchit-dopolnitelnoe-pole
func (s *Service) GetAttribute(ctx context.Context, entityType, id string, opts ...transport.RequestOption) (*Attribute, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get attribute: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get attribute: id %q must be a UUID", id)
	}
	operation := "get " + entityType + " attribute " + id
	var result Attribute
	if err := s.do(ctx, operation, http.MethodGet, "entity/"+entityType+"/metadata/attributes/"+id, nil, nil, &result, opts); err != nil {
		return nil, err
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	return &result, nil
}

// entity — критерий успеха метаданных типа: meta.href; id у них нет.
func (s *Service) entity(ctx context.Context, operation, endpoint string, opts []transport.RequestOption) (*Entity, error) {
	var result Entity
	if err := s.do(ctx, operation, http.MethodGet, endpoint, nil, nil, &result, opts); err != nil {
		return nil, err
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	return &result, nil
}

// do выполняет один запрос и проверяет HTTP 200; содержательный критерий проверяет вызывающий метод.
func (s *Service) do(ctx context.Context, operation, method, endpoint string, params url.Values, body, dst any, opts []transport.RequestOption) error {
	resp, err := s.client.Send(ctx, method, endpoint, params, body, dst, opts...)
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
