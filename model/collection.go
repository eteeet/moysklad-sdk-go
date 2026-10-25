package model

type Collection[T any] struct {
	Context *ResponseContext `json:"context,omitzero"`
	Meta    CollectionMeta   `json:"meta"`
	Rows    []T              `json:"rows"`
}
