package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-podderzhka-null

import (
	"bytes"
	"encoding/json/v2"
	"errors"
)

const (
	fieldUnset uint8 = iota
	fieldNull
	fieldValue
)

// Field используется в моделях запросов с тегом json:",omitzero"; нулевое значение означает «не передано».
type Field[T any] struct {
	value T
	state uint8
}

func Value[T any](value T) Field[T] {
	return Field[T]{value: value, state: fieldValue}
}

func Null[T any]() Field[T] {
	return Field[T]{state: fieldNull}
}

func (f Field[T]) IsZero() bool {
	return f.state == fieldUnset
}

func (f Field[T]) IsNull() bool {
	return f.state == fieldNull
}

func (f Field[T]) Get() (T, bool) {
	return f.value, f.state == fieldValue
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	switch f.state {
	case fieldNull:
		return []byte("null"), nil
	case fieldValue:
		// v2 кодирует nil-слайс (кроме []byte) как [] даже при внешнем json v1.
		return json.Marshal(f.value)
	default:
		return nil, errors.New("cannot marshal an unset Field: use the omitzero tag")
	}
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*f = Null[T]()
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*f = Value(value)
	return nil
}
