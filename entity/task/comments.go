package task

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// MaxCommentsLimit — максимум limit списка комментариев задачи.
const MaxCommentsLimit = 100

// CommentsOptions — параметры одной страницы комментариев. Limit 0 — default сервера (25).
// Filter собирается query.Filter; источник называет фильтры updated и updatedBy (uid)
// («updatedFrom/updatedTo/updatedBy»).
type CommentsOptions struct {
	Limit  int
	Offset int
	Filter string
}

func (s *Service) comments(taskID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "task", ParentID: taskID, Path: "notes", Item: "comment", MetaType: "tasknote"}
}

// ListComments читает одну страницу комментариев задачи без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-poluchit-kommentarii-zadachi
func (s *Service) ListComments(ctx context.Context, taskID string, options CommentsOptions, opts ...transport.RequestOption) (*model.Collection[Comment], error) {
	transport.ResetResponseInfo(opts)
	operation := "list task " + taskID + " notes"
	if !validID(taskID) {
		return nil, fmt.Errorf("moysklad: %s: task id must be a UUID", operation)
	}
	params, err := query.ListOptions{Limit: options.Limit, Offset: options.Offset, Filter: options.Filter}.Values(MaxCommentsLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Comment](ctx, s.client, operation, "entity/task/"+taskID+"/notes", params, opts...)
}

// GetComment — успех требует id в ответе (примеры Get/Update источника его не содержат, пример
// создания и фикстура задачи — содержат).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-poluchit-kommentarij-k-zadache
func (s *Service) GetComment(ctx context.Context, taskID, commentID string, opts ...transport.RequestOption) (*Comment, error) {
	return nested.Get[Comment](ctx, s.comments(taskID), commentID, opts)
}

// CreateComment — пустой text — ошибка до запроса. Ответ источника — массив из одного
// комментария; принимается и объект. Не повторять POST автоматически: возможен дубль.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-kommentarij-zadachi
func (s *Service) CreateComment(ctx context.Context, taskID string, input CommentRequest, opts ...transport.RequestOption) (*Comment, error) {
	if input.Text == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create task %s comment: text must not be empty", taskID)
	}
	return nested.Create[Comment](ctx, s.comments(taskID), input, opts)
}

// CreateManyComments — один POST notes с массивом комментариев (пример «создания сразу
// нескольких»). Массового изменения комментариев источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-kommentarij-zadachi
func (s *Service) CreateManyComments(ctx context.Context, taskID string, items []CommentRequest, opts ...transport.RequestOption) ([]batch.Result[Comment], error) {
	return nested.CreateMany[Comment](ctx, s.comments(taskID), items, opts)
}

// UpdateComment заменяет текст; пустой text — ошибка до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-izmenit-kommentarij-k-zadache
func (s *Service) UpdateComment(ctx context.Context, taskID, commentID string, input CommentRequest, opts ...transport.RequestOption) (*Comment, error) {
	return nested.Update[Comment](ctx, s.comments(taskID), commentID, input, opts)
}

// DeleteComment — DELETE комментария. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-udalit-kommentarij
func (s *Service) DeleteComment(ctx context.Context, taskID, commentID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.comments(taskID), commentID, opts)
}

// DeleteManyComments — один POST notes/delete; источник отвечает 200 с пустым телом, результата
// по элементам нет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-udalenie-komentariev-k-zadache
func (s *Service) DeleteManyComments(ctx context.Context, taskID string, ids []string, opts ...transport.RequestOption) error {
	operation := "delete many task " + taskID + " notes"
	if !validID(taskID) {
		transport.ResetResponseInfo(opts)
		return fmt.Errorf("moysklad: %s: task id must be a UUID", operation)
	}
	return nested.DeleteRows(ctx, s.comments(taskID), operation, ids, opts)
}
