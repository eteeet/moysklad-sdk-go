// Package batch — типы массовых операций сервисов: CreateMany, UpdateMany и DeleteMany.
//
// Один вызов — один POST с массивом от 1 до MaxItems элементов: SDK не делит массив на части,
// не повторяет запрос и не обещает «всё или ничего». Ответ — по элементу на каждый элемент
// запроса; ошибка одного элемента приходит в его Errors, остальные могут быть записаны.
// Порядок элементов ответа источник не описывает: в примерах создания и удаления товаров он
// совпадает с запросом, в примере удаления ассортимента — нет; сверяйте записанные объекты по
// id, удаления — по Info.
//
// Ошибка вызова — только ошибка всего запроса: сеть, HTTP-статус, лимит элементов, разбор
// ответа. Массив разбирается только у HTTP 200: статус ответа с ошибками отдельных элементов
// источник не указывает, и при другом статусе исход элементов неизвестен. Обработка массива
// на сервере ограничена 5 минутами (ошибка 1088); после таймаута, обрыва или отмены ctx часть
// элементов может быть записана — результат неизвестен, повторять запрос вслепую нельзя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-sozdanie-i-obnovlenie-neskolkih-obuektov
package batch

import (
	"encoding/json/v2"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
)

// MaxItems — лимит элементов в одном массовом запросе.
const MaxItems = 1000

// Update — элемент UpdateMany: UUID объекта и изменяемые поля. SDK передаёт объект с meta
// по ID и полями Request.
type Update[R any] struct {
	ID      string
	Request R
}

// Result — элемент ответа CreateMany или UpdateMany: записанный объект в Value либо ошибки
// API этого элемента в Errors. Пустой массив errors даёт Errors без деталей и nil Value.
type Result[T any] struct {
	Value  *T
	Errors []sdkerrors.Detail
}

// UnmarshalJSON: объект с полем errors — ошибка элемента, иначе — объект T.
func (r *Result[T]) UnmarshalJSON(data []byte) error {
	var head struct {
		Errors []sdkerrors.Detail `json:"errors"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return err
	}
	if head.Errors != nil {
		*r = Result[T]{Errors: head.Errors}
		return nil
	}
	value := new(T)
	if err := json.Unmarshal(data, value); err != nil {
		return err
	}
	*r = Result[T]{Value: value}
	return nil
}

// Deleted — элемент ответа DeleteMany: сообщение сервера об удалении в Info либо ошибки API
// этого элемента в Errors.
type Deleted struct {
	Info   string             `json:"info"`
	Errors []sdkerrors.Detail `json:"errors"`
}
