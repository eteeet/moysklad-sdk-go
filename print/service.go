// Package print читает шаблоны печатных форм и запрашивает печать документов, этикеток и
// ценников. Печать — POST entity/{type}/{id}/export без следования redirect: ответ 303 или 202
// с Location возвращается как Result, файл скачивает отдельный вызов Download.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
package print

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: print service requires a client")
	}
	return &Service{client: client}, nil
}

var entityTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Коды набора шаблонов в пути entity/{type}/metadata/{kind}.
const (
	embedded = "embeddedtemplate"
	custom   = "customtemplate"
)

// ListEmbeddedTemplates читает стандартные шаблоны типа: GET entity/{type}/metadata/embeddedtemplate.
// Для шаблонов этикеток и ценников entityType = "assortment". Скрытые формы сервер не отдаёт;
// параметров списка источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/template
func (s *Service) ListEmbeddedTemplates(ctx context.Context, entityType string, opts ...transport.RequestOption) (*model.Collection[Template], error) {
	transport.ResetResponseInfo(opts)
	return s.listTemplates(ctx, entityType, embedded, "embedded", opts)
}

// GetEmbeddedTemplate читает стандартный шаблон: GET entity/{type}/metadata/embeddedtemplate/{id}.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/template
func (s *Service) GetEmbeddedTemplate(ctx context.Context, entityType, id string, opts ...transport.RequestOption) (*Template, error) {
	transport.ResetResponseInfo(opts)
	return s.getTemplate(ctx, entityType, embedded, "embedded", id, opts)
}

// ListCustomTemplates читает пользовательские шаблоны типа: GET entity/{type}/metadata/customtemplate.
// Для шаблонов этикеток и ценников entityType = "assortment".
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/template
func (s *Service) ListCustomTemplates(ctx context.Context, entityType string, opts ...transport.RequestOption) (*model.Collection[Template], error) {
	transport.ResetResponseInfo(opts)
	return s.listTemplates(ctx, entityType, custom, "custom", opts)
}

// GetCustomTemplate читает пользовательский шаблон: GET entity/{type}/metadata/customtemplate/{id}.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/template
func (s *Service) GetCustomTemplate(ctx context.Context, entityType, id string, opts ...transport.RequestOption) (*Template, error) {
	transport.ResetResponseInfo(opts)
	return s.getTemplate(ctx, entityType, custom, "custom", id, opts)
}

func (s *Service) listTemplates(ctx context.Context, entityType, kind, name string, opts []transport.RequestOption) (*model.Collection[Template], error) {
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get %s templates: entity type %q must match [a-z][a-z0-9]*", name, entityType)
	}
	operation := "get " + entityType + " " + name + " templates"
	return list.Page[Template](ctx, s.client, operation, "entity/"+entityType+"/metadata/"+kind, nil, opts...)
}

func (s *Service) getTemplate(ctx context.Context, entityType, kind, name, id string, opts []transport.RequestOption) (*Template, error) {
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get %s template: entity type %q must match [a-z][a-z0-9]*", name, entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get %s template: id %q must be a UUID", name, id)
	}
	operation := "get " + entityType + " " + name + " template " + id
	var result Template
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/"+entityType+"/metadata/"+kind+"/"+id, nil, nil, &result, opts...)
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

// Export печатает один документ по шаблону: POST entity/{type}/{id}/export. Redirect не
// выполняется — Result несёт 303 или 202 и Location. Запись не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
func (s *Service) Export(ctx context.Context, entityType, id string, input ExportRequest, opts ...transport.RequestOption) (*Result, error) {
	transport.ResetResponseInfo(opts)
	if err := checkTarget(entityType, id); err != nil {
		return nil, fmt.Errorf("moysklad: export: %w", err)
	}
	operation := "export " + entityType + " " + id
	if input.Template.Meta.Href == "" || input.Extension == "" {
		return nil, fmt.Errorf("moysklad: %s: template.meta.href and extension are required", operation)
	}
	return s.export(ctx, operation, entityType, id, input, opts)
}

// ExportSet печатает комплект документов в pdf: POST entity/{type}/{id}/export с templates.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
func (s *Service) ExportSet(ctx context.Context, entityType, id string, input ExportSetRequest, opts ...transport.RequestOption) (*Result, error) {
	transport.ResetResponseInfo(opts)
	if err := checkTarget(entityType, id); err != nil {
		return nil, fmt.Errorf("moysklad: export set: %w", err)
	}
	operation := "export set " + entityType + " " + id
	if len(input.Templates) == 0 {
		return nil, fmt.Errorf("moysklad: %s: templates must not be empty", operation)
	}
	for i, item := range input.Templates {
		if item.Template.Meta.Href == "" || item.Count < 1 {
			return nil, fmt.Errorf("moysklad: %s: templates[%d]: template.meta.href and positive count are required", operation, i)
		}
	}
	return s.export(ctx, operation, entityType, id, input, opts)
}

// ExportLabels печатает этикетки и ценники: POST entity/{type}/{id}/export. Источник называет
// товары, услуги, комплекты и модификации (product, service, bundle, variant); тип проверяется
// только по форме, остальное — сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/labels
func (s *Service) ExportLabels(ctx context.Context, entityType, id string, input LabelRequest, opts ...transport.RequestOption) (*Result, error) {
	transport.ResetResponseInfo(opts)
	if err := checkTarget(entityType, id); err != nil {
		return nil, fmt.Errorf("moysklad: export labels: %w", err)
	}
	operation := "export labels " + entityType + " " + id
	if input.Organization.Meta.Href == "" || input.Count < 1 || input.SalePrice.PriceType.Meta.Href == "" || input.Template.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization, positive count, salePrice.priceType and template are required", operation)
	}
	return s.export(ctx, operation, entityType, id, input, opts)
}

// Download скачивает готовую форму (Result со статусом 303) в w отдельным GET через
// transport.Client.Download: Location чужого origin отклоняется до отправки, такой файл
// вызывающий код скачивает сам. Ссылка живёт 5 минут.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
func (s *Service) Download(ctx context.Context, result Result, w io.Writer, opts ...transport.RequestOption) (int64, error) {
	transport.ResetResponseInfo(opts)
	if result.StatusCode != http.StatusSeeOther || result.Location == "" {
		return 0, fmt.Errorf("moysklad: download print: result has status %d, want 303 with Location", result.StatusCode)
	}
	n, err := s.client.Download(ctx, result.Location, w, opts...)
	if err != nil {
		return n, fmt.Errorf("download print: %w", err)
	}
	return n, nil
}

func (s *Service) export(ctx context.Context, operation, entityType, id string, body any, opts []transport.RequestOption) (*Result, error) {
	resp, err := s.client.SendNoRedirect(ctx, http.MethodPost, "entity/"+entityType+"/"+id+"/export", body, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	location, err := resp.Location()
	if err != nil {
		// Текст ошибки разбора содержит ссылку целиком, а подпись в ней — секрет.
		return nil, fmt.Errorf("moysklad: %s: response %d has no valid Location", operation, resp.StatusCode)
	}
	return &Result{StatusCode: resp.StatusCode, Location: location.String()}, nil
}

func checkTarget(entityType, id string) error {
	if !entityTypePattern.MatchString(entityType) {
		return fmt.Errorf("entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return fmt.Errorf("id %q must be a UUID", id)
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
