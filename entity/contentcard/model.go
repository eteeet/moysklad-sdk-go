package contentcard

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// ContentCard — карточка контента товара или услуги для площадки продаж. Images — метаданные
// коллекции изображений; сами изображения — пакет images (тип contentcard).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard
type ContentCard struct {
	Meta            model.Meta                        `json:"meta"`
	ID              string                            `json:"id"`
	AccountID       string                            `json:"accountId"`
	Name            string                            `json:"name"`
	Description     string                            `json:"description"`
	CardContentName string                            `json:"cardContentName"`
	Assortment      *model.Reference                  `json:"assortment,omitzero"`
	SalePlatform    *model.Reference                  `json:"salePlatform,omitzero"`
	SalesChannels   []model.Reference                 `json:"salesChannels"`
	Images          *model.Collection[jsontext.Value] `json:"images,omitzero"`
}
