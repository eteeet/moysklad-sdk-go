// Package images работает с изображениями товара, комплекта, модификации и карточки контента:
// entity/{type}/{id}/images. Содержимое скачивается по Image.Meta.DownloadHref (Download):
// источник отдаёт 302 на временную ссылку, живущую 1 минуту, — transport идёт по ней сам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images
package images

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// PermanentHref — значение fields для List: meta.downloadHref станет постоянной ссылкой
// (только на платных тарифах — проверяет сервер).
const PermanentHref = "downloadPermanentHref"

// Типы владельцев изображений по источнику; другие значения отклоняются до запроса.
var owners = map[string]bool{"product": true, "bundle": true, "variant": true, "contentcard": true}

var entityType = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Image — изображение (meta.type image). Miniature — миниатюра, Tiny — уменьшенное изображение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#2-izobrazheniya
type Image struct {
	Meta      model.Meta     `json:"meta"`
	Title     string         `json:"title"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"` // байты
	Updated   model.DateTime `json:"updated,omitzero"`
	Miniature *model.Meta    `json:"miniature,omitzero"`
	Tiny      *model.Meta    `json:"tiny,omitzero"`
}

// Upload — новое изображение: имя файла с расширением и непустое содержимое (JSON кодирует его в
// Base64); пустые filename и content источник не принимает — ошибка до запроса.
type Upload struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

// Item — элемент нового списка изображений в Replace: существующее по Meta (Keep) или новое по
// Filename и Content (New).
type Item struct {
	Meta     *model.Meta `json:"meta,omitzero"`
	Filename string      `json:"filename,omitempty"`
	Content  []byte      `json:"content,omitzero"`
}

// Keep оставляет существующее изображение в списке Replace.
func Keep(image Image) Item {
	meta := image.Meta
	return Item{Meta: &meta}
}

// New добавляет новое изображение в списке Replace.
func New(upload Upload) Item {
	return Item{Filename: upload.Filename, Content: upload.Content}
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: images service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу изображений объекта. page.Fields — только PermanentHref или "".
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#3-poluchit-spisok-izobrazhenij-tovara-komplekta-modifikacii-ili-kartochki-kontenta
func (s *Service) List(ctx context.Context, typ, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Image], error) {
	transport.ResetResponseInfo(opts)
	operation := "get " + typ + " " + id + " images"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if page.Fields != "" && page.Fields != PermanentHref {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Image](ctx, s.client, operation, "entity/"+typ+"/"+id+"/images", params, opts...)
}

// Add добавляет одно изображение и возвращает весь список изображений объекта из ответа.
// Не повторять POST после неоднозначного сбоя. Размер тела ограничен transport.Config.MaxBodySize.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#3-dobavit-izobrazhenie-k-tovaru-komplektu-modifikacii-ili-kartochke-kontenta
func (s *Service) Add(ctx context.Context, typ, id string, upload Upload, opts ...transport.RequestOption) ([]Image, error) {
	transport.ResetResponseInfo(opts)
	operation := "create " + typ + " " + id + " image"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if err := checkName(upload.Filename); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if len(upload.Content) == 0 {
		return nil, fmt.Errorf("moysklad: %s: content must not be empty", operation)
	}
	return s.post(ctx, operation, "entity/"+typ+"/"+id+"/images", upload, opts)
}

// Replace задаёт список изображений целиком: Keep оставляет существующее, New добавляет новое,
// не переданные удаляются. Пустой список удаляет все. Ответ — новый список.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#3-izmenenie-spiska-izobrazhenij-u-tovara-komplekta-modifikacii-ili-kartochki-kontenta
func (s *Service) Replace(ctx context.Context, typ, id string, items []Item, opts ...transport.RequestOption) ([]Image, error) {
	transport.ResetResponseInfo(opts)
	operation := "update " + typ + " " + id + " images"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if items == nil {
		items = []Item{}
	}
	for i, item := range items {
		if (item.Meta == nil) == (item.Filename == "") {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: either meta or filename must be set", operation, i)
		}
		if item.Meta == nil {
			if err := checkName(item.Filename); err != nil {
				return nil, fmt.Errorf("moysklad: %s: items[%d]: %w", operation, i, err)
			}
			if len(item.Content) == 0 {
				return nil, fmt.Errorf("moysklad: %s: items[%d]: content must not be empty", operation, i)
			}
		} else if item.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: meta.href must not be empty", operation, i)
		}
	}
	return s.post(ctx, operation, "entity/"+typ+"/"+id+"/images", items, opts)
}

// Delete удаляет одно изображение; успех — 200 или 204, повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#3-udalit-izobrazhenie
func (s *Service) Delete(ctx context.Context, typ, id, imageID string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete " + typ + " " + id + " image " + imageID
	if err := owner(typ, id); err != nil {
		return fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if !validID(imageID) {
		return fmt.Errorf("moysklad: %s: image id must be a UUID", operation)
	}
	return write.Delete(ctx, s.client, operation, "entity/"+typ+"/"+id+"/images/"+imageID, opts)
}

// DeleteMany — POST entity/{type}/{id}/images/delete с meta изображений (тип image).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#3-udalit-gruppu-izobrazhenij
func (s *Service) DeleteMany(ctx context.Context, typ, id string, imageIDs []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	operation := "delete many " + typ + " " + id + " images"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return write.DeleteManyNested(ctx, s.client, operation, "entity/"+typ+"/"+id+"/images", "image", imageIDs, opts)
}

// Download копирует содержимое изображения в w по image.Meta.DownloadHref (transport.Client.Download).
func (s *Service) Download(ctx context.Context, image Image, w io.Writer, opts ...transport.RequestOption) (int64, error) {
	transport.ResetResponseInfo(opts)
	const operation = "download image"
	if image.Meta.DownloadHref == "" {
		return 0, fmt.Errorf("moysklad: %s: meta.downloadHref must not be empty", operation)
	}
	n, err := s.client.Download(ctx, image.Meta.DownloadHref, w, opts...)
	if err != nil {
		return n, fmt.Errorf("%s: %w", operation, err)
	}
	return n, nil
}

func (s *Service) post(ctx context.Context, operation, endpoint string, body any, opts []transport.RequestOption) ([]Image, error) {
	var result []Image
	resp, err := s.client.Send(ctx, http.MethodPost, endpoint, nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result == nil {
		return nil, fmt.Errorf("moysklad: %s: response must be a non-null array", operation)
	}
	return result, nil
}

func owner(typ, id string) error {
	if !entityType.MatchString(typ) || !owners[typ] {
		return fmt.Errorf("entity type must be product, bundle, variant or contentcard")
	}
	if !validID(id) {
		return fmt.Errorf("id must be a UUID")
	}
	return nil
}

// Имя уходит на сервер и потом приходит в filename; разделители пути и NUL не пропускаем.
func checkName(name string) error {
	if name == "" || strings.ContainsAny(name, "/\\\x00") || name == "." || name == ".." {
		return fmt.Errorf("filename must be a plain non-empty name")
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
