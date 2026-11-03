package processingprocess

import "github.com/eteeet/moysklad-sdk-go/model"

// ProcessingProcess — техпроцесс. Positions — метаданные коллекции позиций (этапов); строки читает
// GetPositions, rows в ответе бывают только при expand (вне SDK).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#2-tehprocess
type ProcessingProcess struct {
	Meta         model.Meta                  `json:"meta"`
	ID           string                      `json:"id"`
	AccountID    string                      `json:"accountId"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	ExternalCode string                      `json:"externalCode"`
	Archived     bool                        `json:"archived"`
	Shared       bool                        `json:"shared"`
	Updated      model.DateTime              `json:"updated,omitzero"`
	Owner        *model.Reference            `json:"owner,omitzero"`
	Group        *model.Reference            `json:"group,omitzero"`
	Positions    *model.Collection[Position] `json:"positions,omitzero"`
}

// Position — позиция техпроцесса (meta.type processingprocessposition): один этап. Собственная
// форма, не positions.Position документов. NextPositions — ссылки на следующие позиции этого же
// техпроцесса (тип MetaArray; в примерах источника поля нет, форма — массив {meta}, как у
// MetaArray performers этапа производства).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#2-tehprocess
type Position struct {
	Meta            model.Meta        `json:"meta"`
	ID              string            `json:"id"`
	AccountID       string            `json:"accountId"`
	ProcessingStage *model.Reference  `json:"processingstage,omitzero"`
	NextPositions   []model.Reference `json:"nextPositions"`
}
