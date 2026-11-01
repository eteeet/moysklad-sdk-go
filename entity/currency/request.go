package currency

import (
	"errors"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// CreateRequest описывает обе формы создания:
//   - пользовательская валюта — Name, Code и ISOCode обязательны;
//   - системная — System=true и Code или ISOCode, остальное сервер берёт из справочника;
//     дополнительно можно указать RateUpdateType (auto или manual) и Margin.
//
// Rate не может быть 0. system и rateUpdateType таблица помечает только для чтения, но раздел
// «Создать системную валюту» прямо требует передать их при создании: это исключение из правила
// «пометка сильнее примера», поэтому оба поля есть только в CreateRequest, не в UpdateRequest.
// Остальные поля без пометки (fullName, archived, multiplicity, indirect, majorUnit, minorUnit)
// пишутся в обе формы; совместимость с системной валютой проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency#3-sozdat-novuyu-valyutu
type CreateRequest struct {
	System         model.Field[bool]    `json:"system,omitzero"`
	RateUpdateType model.Field[string]  `json:"rateUpdateType,omitzero"`
	Name           model.Field[string]  `json:"name,omitzero"`
	FullName       model.Field[string]  `json:"fullName,omitzero"`
	Code           model.Field[string]  `json:"code,omitzero"`
	ISOCode        model.Field[string]  `json:"isoCode,omitzero"`
	Archived       model.Field[bool]    `json:"archived,omitzero"`
	Rate           model.Field[float64] `json:"rate,omitzero"`
	Multiplicity   model.Field[int]     `json:"multiplicity,omitzero"`
	Indirect       model.Field[bool]    `json:"indirect,omitzero"`
	Margin         model.Field[float64] `json:"margin,omitzero"`
	MajorUnit      model.Field[Unit]    `json:"majorUnit,omitzero"`
	MinorUnit      model.Field[Unit]    `json:"minorUnit,omitzero"`
}

// UpdateRequest: Rate не может быть 0, Name, Code и ISOCode — пустыми. У системных валют сервер
// не даёт менять name, fullName, code, isoCode, majorUnit, minorUnit, у валюты учёта — курс.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency#3-izmenit-valyutu
type UpdateRequest struct {
	Name         model.Field[string]  `json:"name,omitzero"`
	FullName     model.Field[string]  `json:"fullName,omitzero"`
	Code         model.Field[string]  `json:"code,omitzero"`
	ISOCode      model.Field[string]  `json:"isoCode,omitzero"`
	Archived     model.Field[bool]    `json:"archived,omitzero"`
	Rate         model.Field[float64] `json:"rate,omitzero"`
	Multiplicity model.Field[int]     `json:"multiplicity,omitzero"`
	Indirect     model.Field[bool]    `json:"indirect,omitzero"`
	Margin       model.Field[float64] `json:"margin,omitzero"`
	MajorUnit    model.Field[Unit]    `json:"majorUnit,omitzero"`
	MinorUnit    model.Field[Unit]    `json:"minorUnit,omitzero"`
}

func (r CreateRequest) validate() error {
	if system, _ := r.System.Get(); system {
		if !nonEmpty(r.Code) && !nonEmpty(r.ISOCode) {
			return errors.New("system currency requires code or isoCode")
		}
	} else if !nonEmpty(r.Name) || !nonEmpty(r.Code) || !nonEmpty(r.ISOCode) {
		return errors.New("name, code and isoCode are required")
	}
	return checkValues(r.Rate, r.Name, r.Code, r.ISOCode)
}

// checkValues отклоняет переданные rate=0 и пустые name, code, isoCode — их API не принимает.
func checkValues(rate model.Field[float64], texts ...model.Field[string]) error {
	if !rate.IsZero() {
		if value, ok := rate.Get(); !ok || value == 0 {
			return errors.New("rate must not be zero or null")
		}
	}
	for _, text := range texts {
		if !text.IsZero() && !nonEmpty(text) {
			return errors.New("name, code and isoCode must not be empty or null")
		}
	}
	return nil
}

func nonEmpty(f model.Field[string]) bool {
	value, ok := f.Get()
	return ok && value != ""
}
