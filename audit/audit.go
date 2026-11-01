// Package audit читает журнал изменений аккаунта: контексты (одно действие пользователя или
// решения), их события и каталог фильтров. Аудит — не changelog с курсором: SDK не обещает
// неограниченного хранения, exactly-once и строгого порядка; при равных moment и параллельных
// изменениях возможны пропуски и дубли между страницами. Доступность событий проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit
package audit

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// MaxContextsLimit — максимум limit у контекстов; по умолчанию сервер отдаёт 25.
const MaxContextsLimit = 100

// Context — контекст аудита: одно изменение, действие пользователя или решения над одним или
// несколькими объектами. Не путать с context.Context и model.ResponseContext коллекций.
// Все поля только для чтения. EntityType, EventType и ObjectType приходят, только если они
// одинаковы у всех событий контекста. Значения enum открыты.
type Context struct {
	Meta          model.Meta     `json:"meta"`
	ID            string         `json:"id"`
	UID           string         `json:"uid"`
	Moment        model.DateTime `json:"moment,omitzero"`
	Source        string         `json:"source"`
	EntityType    string         `json:"entityType"`
	EventType     string         `json:"eventType"`
	ObjectType    string         `json:"objectType"`
	ObjectCount   int            `json:"objectCount"`
	SupportAccess bool           `json:"supportAccess"`
	// Info — краткое описание. Таблица называет его Object, пример списка — строка
	// («Вход в МойСклад (…)»), фикстура — объект; хранится как пришло.
	Info jsontext.Value `json:"info,omitzero"`
	// Events — метаданные списка событий. Таблица источника называет поле Array(Object), но
	// примеры — объект с meta коллекции; сами события читаются отдельным запросом.
	Events *EventsRef `json:"events,omitzero"`
}

// EventsRef — ссылка на список событий контекста и его размер.
type EventsRef struct {
	Meta model.CollectionMeta `json:"meta"`
}

// FilterCatalog — значения, по которым пользователь может фильтровать контексты. Таблица источника
// называет ключи entityType и eventType, пример и upstream-фикстура — entitytype и eventtype;
// SDK читает обе формы (json/v2 различает регистр) и пишет форму примера.
type FilterCatalog struct {
	Source     []string `json:"source"`
	EntityType []string `json:"entitytype"`
	EventType  []string `json:"eventtype"`
}

func (f *FilterCatalog) UnmarshalJSON(data []byte) error {
	var raw struct {
		Source          []string `json:"source"`
		EntityType      []string `json:"entitytype"`
		EventType       []string `json:"eventtype"`
		EntityTypeTable []string `json:"entityType"`
		EventTypeTable  []string `json:"eventType"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*f = FilterCatalog{Source: raw.Source, EntityType: raw.EntityType, EventType: raw.EventType}
	if f.EntityType == nil {
		f.EntityType = raw.EntityTypeTable
	}
	if f.EventType == nil {
		f.EventType = raw.EventTypeTable
	}
	return nil
}

// ContextsOptions — параметры одной страницы контекстов; других параметров маршрут не описывает.
// Filter собирается query.Filter по полям moment (>=, <=, московское время,
// model.DateTime.QueryValue), employee (href), uid, source, eventType, entityType (код сущности
// или href пользовательского справочника) и application (href установленного решения
// entity/application/{id}, не appUid).
type ContextsOptions struct {
	Limit  int // 0 — default сервера 25; не больше MaxContextsLimit
	Offset int
	Filter string
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: audit service requires a client")
	}
	return &Service{client: client}, nil
}

// ListContexts читает одну страницу контекстов без чтения событий и автообхода;
// следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit#3-poluchit-konteksty
func (s *Service) ListContexts(ctx context.Context, options ContextsOptions, opts ...transport.RequestOption) (*model.Collection[Context], error) {
	transport.ResetResponseInfo(opts)
	params, err := query.ListOptions{Limit: options.Limit, Offset: options.Offset, Filter: options.Filter}.Values(MaxContextsLimit)
	if err != nil {
		return nil, fmt.Errorf("list audit contexts: %w", err)
	}
	return list.Page[Context](ctx, s.client, "list audit contexts", "audit", params, opts...)
}

// GetFilters читает каталог фильтров, доступных пользователю.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/audit/audit#3-poluchit-filtry
func (s *Service) GetFilters(ctx context.Context, opts ...transport.RequestOption) (*FilterCatalog, error) {
	transport.ResetResponseInfo(opts)
	var result FilterCatalog
	resp, err := s.client.Send(ctx, http.MethodGet, "audit/metadata/filters", nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("get audit filters: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: get audit filters: unexpected HTTP status %d", resp.StatusCode)
	}
	// Критерий успеха каталога — ключ source; пустой список допустим, null и {} — нет.
	if result.Source == nil {
		return nil, fmt.Errorf("moysklad: get audit filters: response has no source")
	}
	return &result, nil
}
