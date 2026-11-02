// Package files работает с файлами, прикреплёнными к операции, номенклатуре, задаче или
// контрагенту: entity/{type}/{id}/files. Отдельного ресурса файла нет; изображения товаров
// в SDK не входят.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/files
package files

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// MaxUpload — файлов в одном запросе добавления; всего у объекта не больше 100 (проверяет сервер).
const MaxUpload = 10

var entityType = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// File — метаданные файла. Meta.DownloadHref — ссылка на содержимое; Miniature и Tiny есть только
// у изображений. Тип содержимого SDK не проверяет: файл скачивается как есть.
type File struct {
	Meta      model.Meta       `json:"meta"`
	Title     string           `json:"title"`
	Filename  string           `json:"filename"`
	Size      int64            `json:"size"` // байты
	Created   model.DateTime   `json:"created,omitzero"`
	CreatedBy *model.Reference `json:"createdBy,omitzero"`
	Miniature *model.Meta      `json:"miniature,omitzero"`
	Tiny      *model.Meta      `json:"tiny,omitzero"`
}

// Upload — новый файл: имя с расширением и содержимое; JSON кодирует Content в Base64, nil — как
// пустой файл. Длину имени (String(255)) проверяет сервер.
type Upload struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: files service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу файлов объекта без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/files#3-poluchit-spisok-fajlov-operacii-nomenklatury-zadachi-ili-kontragenta
func (s *Service) List(ctx context.Context, typ, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[File], error) {
	transport.ResetResponseInfo(opts)
	operation := "get " + typ + " " + id + " files"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[File](ctx, s.client, operation, "entity/"+typ+"/"+id+"/files", params, opts...)
}

// Upload добавляет от 1 до 10 файлов и возвращает весь список файлов объекта из ответа сервера.
// Существующие файлы не удаляются; одинаковые filename и content дают один id. Лимит 100 файлов
// на объект проверяет сервер, SDK не считает их заранее. Размер закодированного тела ограничен
// transport.Config.MaxBodySize (Base64 больше исходного файла примерно на треть).
// Не повторять POST автоматически после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/files#3-dobavit-fajly-k-operacii-nomenklature-ili-kontragentu
func (s *Service) Upload(ctx context.Context, typ, id string, uploads []Upload, opts ...transport.RequestOption) ([]File, error) {
	transport.ResetResponseInfo(opts)
	operation := "create " + typ + " " + id + " files"
	if err := owner(typ, id); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if len(uploads) == 0 || len(uploads) > MaxUpload {
		return nil, fmt.Errorf("moysklad: %s: uploads must contain 1 to %d files", operation, MaxUpload)
	}
	for i, u := range uploads {
		// Имя уходит на сервер и потом приходит в filename; разделители пути и NUL не пропускаем.
		if u.Filename == "" || strings.ContainsAny(u.Filename, "/\\\x00") || u.Filename == "." || u.Filename == ".." {
			return nil, fmt.Errorf("moysklad: %s: uploads[%d]: filename must be a plain non-empty name", operation, i)
		}
	}
	// nil []byte json v1 пишет как null, а источник ждёт строку Base64.
	body := make([]Upload, len(uploads))
	for i, u := range uploads {
		body[i] = u
		if u.Content == nil {
			body[i].Content = []byte{}
		}
	}
	var result []File
	resp, err := s.client.Send(ctx, http.MethodPost, "entity/"+typ+"/"+id+"/files", nil, body, &result, opts...)
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

// Delete удаляет файл fileID объекта typ/id. По источнику сервер удаляет первый найденный файл
// с этим id у операции, номенклатуры, задачи или контрагента. Успех — 200 или 204; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/files#3-udalit-fajl
func (s *Service) Delete(ctx context.Context, typ, id, fileID string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete " + typ + " " + id + " file " + fileID
	if err := owner(typ, id); err != nil {
		return fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if !validID(fileID) {
		return fmt.Errorf("moysklad: %s: file id must be a UUID", operation)
	}
	return write.Delete(ctx, s.client, operation, "entity/"+typ+"/"+id+"/files/"+fileID, opts)
}

// Download копирует содержимое файла в w по file.Meta.DownloadHref (transport.Client.Download):
// один GET, redirect на подписанный адрес без Authorization, повтор только до первой записи.
// Время жизни ссылки задаёт transport.WithDownloadExpiration. Скачивание миниатюр источник не
// описывает (в примере есть только поля miniature) — SDK его не предлагает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-ssylki-na-fajly
func (s *Service) Download(ctx context.Context, file File, w io.Writer, opts ...transport.RequestOption) (int64, error) {
	transport.ResetResponseInfo(opts)
	const operation = "download file"
	if file.Meta.DownloadHref == "" {
		return 0, fmt.Errorf("moysklad: %s: meta.downloadHref must not be empty", operation)
	}
	n, err := s.client.Download(ctx, file.Meta.DownloadHref, w, opts...)
	if err != nil {
		return n, fmt.Errorf("%s: %w", operation, err)
	}
	return n, nil
}

func owner(typ, id string) error {
	if !entityType.MatchString(typ) {
		return fmt.Errorf("entity type must match [a-z][a-z0-9]*")
	}
	if !validID(id) {
		return fmt.Errorf("id must be a UUID")
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
