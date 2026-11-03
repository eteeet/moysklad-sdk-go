package metadata

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// AttributeCreateRequest — новое доп. поле: name и type обязательны. Type — тип значения
// (string, long, double, boolean, time, link, file, text) или тип сущностей справочника
// (product, counterparty, …, customentity); для customentity нужен CustomEntityMeta. Имена
// уникальны в пределах сущности, show=false вместе с required=true запрещены — проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud
type AttributeCreateRequest struct {
	Name             string                  `json:"name"`
	Type             string                  `json:"type"`
	Required         model.Field[bool]       `json:"required,omitzero"`
	Show             model.Field[bool]       `json:"show,omitzero"`
	Description      model.Field[string]     `json:"description,omitzero"`
	CustomEntityMeta model.Field[model.Meta] `json:"customEntityMeta,omitzero"`
}

// AttributeUpdateRequest: type после создания не меняется; пример массового изменения источника
// передаёт его вместе с meta, поэтому поле есть — допустимо только текущее значение.
type AttributeUpdateRequest struct {
	Name        model.Field[string] `json:"name,omitzero"`
	Type        model.Field[string] `json:"type,omitzero"`
	Required    model.Field[bool]   `json:"required,omitzero"`
	Show        model.Field[bool]   `json:"show,omitzero"`
	Description model.Field[string] `json:"description,omitzero"`
}

// CreateAttribute — POST entity/{type}/metadata/attributes. Доп. поля услуг и комплектов
// хранятся в метаданных товаров: entityType = "product". Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud#3-sozdanie-novogo-dopolnitelnogo-polya-cherez-json-api
func (s *Service) CreateAttribute(ctx context.Context, entityType string, input AttributeCreateRequest, opts ...transport.RequestOption) (*Attribute, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: create attribute: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	operation := "create " + entityType + " attribute"
	if input.Name == "" || input.Type == "" {
		return nil, fmt.Errorf("moysklad: %s: name and type are required", operation)
	}
	return s.attribute(ctx, operation, http.MethodPost, "entity/"+entityType+"/metadata/attributes", input, opts)
}

// UpdateAttribute — PUT entity/{type}/metadata/attributes/{id}; незаданные свойства не меняются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud#3-redaktirovanie-dopolnitelnyh-polej-cherez-json-api
func (s *Service) UpdateAttribute(ctx context.Context, entityType, id string, input AttributeUpdateRequest, opts ...transport.RequestOption) (*Attribute, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: update attribute: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update attribute: id %q must be a UUID", id)
	}
	operation := "update " + entityType + " attribute " + id
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.attribute(ctx, operation, http.MethodPut, "entity/"+entityType+"/metadata/attributes/"+id, input, opts)
}

// DeleteAttribute удаляет доп. поле, даже обязательное и заполненное: DELETE, успех — 200 или
// 204; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud#3-udalenie-dopolnitelnogo-polya-cherez-json-api
func (s *Service) DeleteAttribute(ctx context.Context, entityType, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return fmt.Errorf("moysklad: delete attribute: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return fmt.Errorf("moysklad: delete attribute: id %q must be a UUID", id)
	}
	return write.Delete(ctx, s.client, "delete "+entityType+" attribute "+id, "entity/"+entityType+"/metadata/attributes/"+id, opts)
}

// CreateManyAttributes и UpdateManyAttributes — один POST entity/{type}/metadata/attributes с
// массивом: элемент без meta создаёт поле, с meta (тип attributemetadata) — изменяет. Правила
// ответа — пакет batch; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud
func (s *Service) CreateManyAttributes(ctx context.Context, entityType string, items []AttributeCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Attribute], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: create many attributes: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	for i, item := range items {
		if item.Name == "" || item.Type == "" {
			return nil, fmt.Errorf("moysklad: create many %s attributes: items[%d]: name and type are required", entityType, i)
		}
	}
	return write.CreateManyNested[Attribute](ctx, s.client, "create many "+entityType+" attributes", "entity/"+entityType+"/metadata/attributes", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud
func (s *Service) UpdateManyAttributes(ctx context.Context, entityType string, items []batch.Update[AttributeUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Attribute], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: update many attributes: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	return write.UpdateManyNested[Attribute](ctx, s.client, "update many "+entityType+" attributes", "entity/"+entityType+"/metadata/attributes", "attributemetadata", items, opts)
}

// DeleteManyAttributes — POST entity/{type}/metadata/attributes/delete с meta полей; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-attributes-crud#3-massovoe-udalenie-dopolnitelnyh-polej-cherez-json-api
func (s *Service) DeleteManyAttributes(ctx context.Context, entityType string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: delete many attributes: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	return write.DeleteManyNested(ctx, s.client, "delete many "+entityType+" attributes", "entity/"+entityType+"/metadata/attributes", "attributemetadata", ids, opts)
}

// attribute — критерий успеха записи: meta.href, как у GetAttribute.
func (s *Service) attribute(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Attribute, error) {
	var result Attribute
	if err := s.do(ctx, operation, method, endpoint, nil, body, &result, opts); err != nil {
		return nil, err
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	return &result, nil
}
