package aggregatepack

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// AggregatePack — транспортная упаковка. Level — максимальный уровень вложенности, только чтение;
// Positions — метаданные коллекции позиций (ListPositions); записать позиции можно в Create/Update. ChildrenList — вложенные упаковки: приходит и
// обрабатывается сервером только с transport.WithBetaFeature(BetaChildrenList).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit
type AggregatePack struct {
	Meta         model.Meta                        `json:"meta"`
	ID           string                            `json:"id"`
	AccountID    string                            `json:"accountId"`
	Updated      model.DateTime                    `json:"updated,omitzero"`
	Moment       model.DateTime                    `json:"moment,omitzero"`
	Level        int                               `json:"level"`
	Barcodes     []model.Barcode                   `json:"barcodes"`
	ChildrenList []model.Reference                 `json:"childrenList"`
	Shared       bool                              `json:"shared"`
	Owner        *model.Reference                  `json:"owner,omitzero"`
	Group        *model.Reference                  `json:"group,omitzero"`
	Positions    *model.Collection[jsontext.Value] `json:"positions,omitzero"`
}

// BetaChildrenList — значение X-Lognex-Remap-Beta-Feature для поля childrenList.
const BetaChildrenList = "aggregatePackChildrenList"
