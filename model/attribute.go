package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-rabota-s-dopolnitelnymi-polyami

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

type Attribute struct {
	Meta     Meta                  `json:"meta"`
	ID       string                `json:"id,omitempty"`
	Name     string                `json:"name,omitempty"`
	Type     string                `json:"type,omitempty"`
	Value    Field[jsontext.Value] `json:"value,omitzero"`
	File     Field[AttributeFile]  `json:"file,omitzero"`
	Download *Meta                 `json:"download,omitzero"`
}

type AttributeFile struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type AttributeReference struct {
	Meta Meta   `json:"meta"`
	Name string `json:"name,omitempty"`
}

// Type определяет сервер по Meta; DateTime не обрезаем: секунды обнулит сервер.
func AttrValue(meta Meta, value any) (Attribute, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return Attribute{}, fmt.Errorf("encode attribute value: %w", err)
	}
	return Attribute{Meta: meta, Value: Value(jsontext.Value(data))}, nil
}

func AttrNull(meta Meta) Attribute {
	return Attribute{Meta: meta, Value: Null[jsontext.Value]()}
}

func AttrFile(meta Meta, filename string, content []byte) Attribute {
	return Attribute{Meta: meta, File: Value(AttributeFile{Filename: filename, Content: content})}
}

func AttrFileNull(meta Meta) Attribute {
	return Attribute{Meta: meta, File: Null[AttributeFile]()}
}

func (a Attribute) AsString() (string, error) {
	return attributeValue[string](a, "string", "text", "link", "time", "file")
}

func (a Attribute) AsBool() (bool, error) {
	return attributeValue[bool](a, "boolean")
}

func (a Attribute) AsLong() (int64, error) {
	return attributeValue[int64](a, "long")
}

func (a Attribute) AsDouble() (float64, error) {
	return attributeValue[float64](a, "double")
}

func (a Attribute) AsTime() (DateTime, error) {
	return attributeValue[DateTime](a, "time")
}

func (a Attribute) AsReference() (AttributeReference, error) {
	// Набор справочников открыт; не связываем type с конкретным meta.type.
	switch a.Type {
	case "string", "text", "link", "time", "file", "boolean", "double", "long":
		return AttributeReference{}, fmt.Errorf("attribute type %q is not a reference", a.Type)
	}
	data, ok := a.Value.Get()
	if !ok || data.Kind() != '{' {
		return AttributeReference{}, fmt.Errorf("attribute type %q: expected reference object", a.Type)
	}
	value, err := attributeValue[AttributeReference](a, a.Type)
	if err != nil {
		return AttributeReference{}, err
	}
	if value.Meta.Href == "" {
		return AttributeReference{}, fmt.Errorf("attribute type %q: reference has no meta.href", a.Type)
	}
	return value, nil
}

func attributeValue[T any](a Attribute, types ...string) (T, error) {
	var zero T
	validType := false
	for _, typ := range types {
		validType = validType || a.Type == typ
	}
	if !validType {
		return zero, fmt.Errorf("attribute type %q: expected %v", a.Type, types)
	}
	data, ok := a.Value.Get()
	if !ok || data.Kind() == 'n' {
		return zero, fmt.Errorf("attribute type %q: value is absent or null", a.Type)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return zero, fmt.Errorf("attribute type %q: decode value: %w", a.Type, err)
	}
	return value, nil
}
