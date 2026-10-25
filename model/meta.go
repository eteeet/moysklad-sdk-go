package model

// Набор значений открыт: сервер добавляет новые ресурсы, значения не валидируем.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-zamechaniya-po-razrabotke-klientskih-prilozhenij
type MetaType string

// Документация Meta — https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-metadannye
type Meta struct {
	Href         string   `json:"href,omitempty"`
	MetadataHref string   `json:"metadataHref,omitempty"`
	Type         MetaType `json:"type,omitempty"`
	MediaType    string   `json:"mediaType,omitempty"`
	UUIDHref     string   `json:"uuidHref,omitempty"`
	// null у несозданной миниатюры декодируется в "": отличать его от отсутствия поля не нужно.
	DownloadHref string `json:"downloadHref,omitempty"`

	// Permitted приходит только при expand объекта без права просмотра;
	// nil — поле не пришло.
	Permitted *bool `json:"permitted,omitzero"`
}

// Документация Meta (коллекция для пагинации) — https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-meta#3-metadannye-kollekcii
type CollectionMeta struct {
	Meta

	Size   int `json:"size"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	NextHref     string `json:"nextHref,omitempty"`
	PreviousHref string `json:"previousHref,omitempty"`
}
